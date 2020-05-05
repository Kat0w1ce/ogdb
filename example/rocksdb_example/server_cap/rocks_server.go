package main

//TODO config server by flags

import (
	"context"
	"fmt"
	"log"
	"net"
	rocksdb_example "ogdb/example/rocksdb_example/proto"
	"github.com/tecbot/gorocksdb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"flag"
)

var (
	ADDRESS string = "localhost"
	PORT    string = "8081"
	MAX		int = 10000
)

type rocksServer struct {
	Db *gorocksdb.DB
	Wo *gorocksdb.WriteOptions
	Ro *gorocksdb.ReadOptions
	cnt int
}

//how to
func (db *rocksServer) Put(ctx context.Context, request *rocksdb_example.PutRequest) (response *rocksdb_example.PutResponse, err error) {
	key, value := request.Key, request.Value
	err = db.Db.Put(db.Wo, []byte(key), []byte(value))
	fmt.Println("1")
	if err != nil {
		return &rocksdb_example.PutResponse{OK: false}, err
	} else {
		if db.cnt < MAX {
			log.Println("put ", key, value)
			db.cnt++
			return &rocksdb_example.PutResponse{OK: true}, nil
		}else {
			log.Println("no more space")
			return &rocksdb_example.PutResponse{OK: false},nil
		}
	}
}
func (db *rocksServer) Delete(ctx context.Context, request *rocksdb_example.DeleteRequest) (response *rocksdb_example.DeleteResponse, err error) {
	key := request.Key
	err = db.Db.Delete(db.Wo, []byte(key))
	if err != nil {
		log.Fatal(err)
		return nil, err
	} else {
		log.Println("delete ", key)
		return &rocksdb_example.DeleteResponse{Ok: true}, nil
	}

}
func (db *rocksServer) Get(ctx context.Context, request *rocksdb_example.GetRequest) (response *rocksdb_example.GetResponse, err error) {
	key := request.Key
	value, err := db.Db.Get(db.Ro, []byte(key))
	if err != nil {
		log.Fatal(err)
		return nil, err
	} else {

		log.Println("get ", key, string(value.Data()))
		return &rocksdb_example.GetResponse{Key: key, Value: string(value.Data())}, nil
	}
}
func main() {
	flag.StringVar(&ADDRESS,"a","localhost","listen address")
	flag.StringVar(&PORT, "p", "2233", "port")
	flag.IntVar(&MAX,"m",10000,"max capilite")
	flag.Parse()
	dbServer := new(rocksServer)
	if err := dbServer.init(); err != nil {
		panic(err)
	}
	defer dbServer.Db.Close()
	listener, err := net.Listen("tcp", ADDRESS+":"+PORT)
	if err != nil {
		log.Fatal("failed listen at " + ADDRESS + ":" + PORT)
	} else {
		log.Println("server is listening")
	}
	rocksdbServer := grpc.NewServer()
	rocksdb_example.RegisterRocksdbServer(rocksdbServer, dbServer)
	reflection.Register(rocksdbServer)
	if err = rocksdbServer.Serve(listener); err != nil {
		log.Fatal("Error")
	}
}

func (db *rocksServer) init() error {
	var err error
	bbto := gorocksdb.NewDefaultBlockBasedTableOptions()
	bbto.SetBlockCache(gorocksdb.NewLRUCache(3 << 30))
	opts := gorocksdb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(bbto)
	opts.SetCreateIfMissing(true)
	//db_server=new(rocksServer{gorocksdb.OpenDb(opts,"/server"),gorocksdb.NewDefaultWriteOptions(),gorocksdb.NewDefaultReadOptions()})
	db.Db, err = gorocksdb.OpenDb(opts,fmt.Sprintln(ADDRESS,':',PORT) )
	if err != nil {
		return err
	}
	db.Ro = gorocksdb.NewDefaultReadOptions()
	db.Wo = gorocksdb.NewDefaultWriteOptions()
	db.cnt=0
	return nil
}
