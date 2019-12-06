package main

import (
	rocksdb_example "awesomeProject/example/rocksdb_example/proto"
	"fmt"
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
	client:=rocksdb_example.NewRocksdbClient(conn)
	//resp,err:=client.Put(context.Background(),&rocksdb_example.PutRequest{Key:"hello",Value:"world2"})
	//if err!=nil{
	//	panic("put error")
	//}else {
	//	fmt.Println(resp.OK)
	//}
	resp2,err:=client.Get(context.Background(),&rocksdb_example.GetRequest{Key:"hello"})
	if err!=nil{
		panic("get error")
	}else {
		fmt.Println(resp2.Key,resp2.Value)
	}

}