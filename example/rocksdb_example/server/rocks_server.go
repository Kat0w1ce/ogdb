package server

import (
	"awesomeProject/example/rocksdb_example/proto"
	"context"
	"fmt"
	"github.com/tecbot/gorocksdb"
	"golang.org/x/crypto/argon2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const (
	ADDRESS string = "localhost"
	PORT    string = "8081"
)

type rocksServer struct {
	Db *gorocksdb.DB
	Wo *gorocksdb.WriteOptions
	Ro *gorocksdb.ReadOptions
}

//how to
func (database *rocksServer) Put(ctx context.Context, request *rocksdb_example.PutRequest) (response *rocksdb_example.PutResponse, err error) {
	key, value := request.Key, request.Value
	err = database.Db.Put(database.Wo, []byte(key), []byte(value))
	if err != nil {
		return &rocksdb_example.PutResponse{OK: false}, err
	} else {
		fmt.Println("put ", key, value)
		return &rocksdb_example.PutResponse{OK: true}, nil
	}
}

func (database *rocksServer) Get(ctx context.Context, request *rocksdb_example.GetRequest) (response *rocksdb_example.GetResponse, err error) {
	key := request.Key
	value, err := database.Db.Get(database.Ro, []byte(key))
	if err != nil {
		return nil, err
	} else {
		return &rocksdb_example.GetResponse{Key: key, Value: string(value.Data())}, nil
	}
}
func main() {
	var dbServer *rocksServer
	var err error
	bbto := gorocksdb.NewDefaultBlockBasedTableOptions()
	bbto.SetBlockCache(gorocksdb.NewLRUCache(3 << 30))
	opts := gorocksdb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(bbto)
	opts.SetCreateIfMissing(true)
	//db_server=new(rocksServer{gorocksdb.OpenDb(opts,"/server"),gorocksdb.NewDefaultWriteOptions(),gorocksdb.NewDefaultReadOptions()})
	dbServer = new(rocksServer)
	dbServer.Db, err = gorocksdb.OpenDb(opts, "dump")
	if err != nil {
		log.Fatal("failed open database")
	}
	dbServer.Ro = gorocksdb.NewNativeReadOptions()
	dbServer.Wo = gorocksdb.NewNativeWriteOptions()
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
