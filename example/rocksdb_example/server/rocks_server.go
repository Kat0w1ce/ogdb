package main


//TODO config server by flags

import (
	"context"
	"fmt"
	"github.com/tecbot/gorocksdb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"ogdb/example/rocksdb_example/proto"
	//"golang.org/x/net/ipv4"
	"flag"
)

var (
	ADDRESS string = "localhost"
	PORT    string = "8081"
)

type rocksServer struct {
	Db *gorocksdb.DB
	Wo *gorocksdb.WriteOptions
	Ro *gorocksdb.ReadOptions
}

//how to
func (db *rocksServer) Put(ctx context.Context, request *rocksdb_example.PutRequest) (response *rocksdb_example.PutResponse, err error) {
	key, value := request.Key, request.Value
	err = db.Db.Put(db.Wo, []byte(key), []byte(value))
	if err != nil {
		return &rocksdb_example.PutResponse{OK: false}, err
	} else {
		fmt.Println("put ", key, value)
		return &rocksdb_example.PutResponse{OK: true}, nil
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

	flag.StringVar(&PORT,"p","2233","port")
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
	db.Db, err = gorocksdb.OpenDb(opts, "dump")
	if err != nil {
		return err
	}
	db.Ro = gorocksdb.NewDefaultReadOptions()
	db.Wo = gorocksdb.NewDefaultWriteOptions()
	return nil
}
