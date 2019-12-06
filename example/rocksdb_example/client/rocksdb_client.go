package main

import (
	rocksdb_example "awesomeProject/example/rocksdb_example/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)
const (
	ADDRESS	string = "localhost:8081"
)

func main()  {
	conn,err:=grpc.Dial(ADDRESS,grpc.WithInsecure())
	if err!=nil{
		log.Fatal("Can't connect: "+ADDRESS)
	}
	defer conn.Close()
	client:=rocksdb_example.NewRocksdbClient(conn)
	delete(&client,"hello")
	put(&client,"og","db")
	put(&client,"psg","lxo")
	put(&client,"kato","wizz")
	get(&client,"og")
	put(&client,"psg","sb")
	get(&client,"psg")
	delete(&client,"psg")
	get(&client,"psg")
	delete(&client,"psg")
}
func put(client *rocksdb_example.RocksdbClient,key string ,value string)   {
	resp,err:=(*client).Put(context.Background(),&rocksdb_example.PutRequest{
		Key: key,
		Value: value,
	})
	if err!=nil{
		log.Fatal("put error")
	}else {
		log.Println("put ",key,value,resp.GetOK())
	}
}

func get(client *rocksdb_example.RocksdbClient,key string )   {
	resp,err:=(*client).Get(context.Background(),&rocksdb_example.GetRequest{
		Key: key,
	})
	if err!=nil{
		log.Fatal("get error")
	}else {
		log.Println("get ",resp.GetKey(),resp.GetValue())
	}
}
func delete(client *rocksdb_example.RocksdbClient,key string)  {
	resp,err:=(*client).Delete(context.Background(),&rocksdb_example.DeleteRequest{
		Key: key,
	})
	if err!=nil{
		log.Fatal("delete error")
	}else {
		log.Println("delete ",resp.GetOk())
	}
}