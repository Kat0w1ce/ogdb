package main

import (
	"bufio"
	"fmt"
	"log"
	rocksdb_example "ogdb/example/rocksdb_example/proto"
	"os"
	"stathat.com/c/consistent"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	ADDRESS string = "localhost"
	PORT    string = "8081"
	connpool map[string]*grpc.ClientConn
)

func main() {
	ip:=[]string {"localhost:9999","localhost:6655"}
	connpool=make(map[string]*grpc.ClientConn)
	for _,addr:=range ip{
		if conn,err:=grpc.Dial(addr,grpc.WithInsecure());err==nil{
			if conn!=nil {
				log.Println("connecting to", addr)
				//todo copy or ref?
				//todo 保存连接和客户的效率
				connpool[addr]=conn;
			}
		}
	}
	defer clear()
	hashring:=consistHashing(ip)
	//client := rocksdb_example.NewRocksdbClient(conn)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		cmd := strings.Split(line, " ")
		//todo add new servers
		if len(cmd)<2 {
			log.Println("invaild cmd")
			continue
		}
		client,err:=getClient(cmd[1],hashring)
		if err!=nil||client==nil{
			log.Println("failed to get client")
			continue
		}
		switch cmd[0] {
		case "put":
			put(client, cmd[1], cmd[2])
		case "get":
			get(client, cmd[1])
		case "delete":
			delete(client, cmd[1])
		default:
			fmt.Println("invalid cmd")
		}
		fmt.Print(">")
	}

}
//todo copy or ref
func getClient(key string,c *consistent.Consistent) (*rocksdb_example.RocksdbClient,error){
	 addr,err:=c.Get(key)
	 fmt.Println(key,"at",addr)
	 if err!=nil{
	 	log.Fatalln("failed to get ip address")
	 	return nil,err
	 }
	 conn:=connpool[addr]
	 client:=rocksdb_example.NewRocksdbClient(conn)
	 return &client,nil
}
func consistHashing(ip []string) *consistent.Consistent{
	c:=consistent.New()
	for _,addr:=range ip{
		c.Add(addr)
	}
	return  c
}
func put(client *rocksdb_example.RocksdbClient, key string, value string) {
	resp, err := (*client).Put(context.Background(), &rocksdb_example.PutRequest{
		Key:   key,
		Value: value,
	})
	if err != nil {
		log.Fatal("put error")
	} else {
		log.Println("put", key, value, resp.GetOK())
	}
}

func get(client *rocksdb_example.RocksdbClient, key string) {
	resp, err := (*client).Get(context.Background(), &rocksdb_example.GetRequest{
		Key: key,
	})
	if err != nil {
		log.Fatal("get error")
	} else {
		log.Println("get", resp.GetKey(), resp.GetValue())
	}
}
func delete(client *rocksdb_example.RocksdbClient, key string) {
	resp, err := (*client).Delete(context.Background(), &rocksdb_example.DeleteRequest{
		Key: key,
	})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("delete", resp.GetOk())
	}
}


func clear() {
	log.Println("clear")
	if connpool != nil {
		for ip, conn := range connpool {
			if conn != nil {
				err := conn.Close()
				if err != nil {
					log.Panic("error when closing connection to", ip)
				}
				log.Println(ip,"closed")
			}
		}
	}
}