package main

import (
	"bufio"
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	rocksdb_example "ogdb/example/rocksdb_example/proto"
	"os"
	"ogdb/consistenthash"
	"strings"
)

var (
	connpool map[string]*grpc.ClientConn
	filepath string
	config string
	replicas int
	cnt0 int32=0
	cnt1 int32=0
	cnt2 int32=0
	cnt	int32 =0
	ip []string
	hasher string
	//ip =[]string {"localhost:9999","127.0.0.1:2233","localhost:6655"}
	//ip	=[]string{"192.168.1.20:9999","192.168.1.20:6655","192.168.1.20:2233","192.168.1.20:8888"}
	)

func main() {
	flag.StringVar(&filepath,"f","data","choose file")
	flag.StringVar(&config,"c","config","cluster config")
	flag.IntVar(&replicas,"r",12,"number of replications")
	flag.StringVar(&hasher,"h","mur","choose hash function")
	flag.Parse()
	conf,err:=os.Open(config)

	if err !=nil{
		log.Fatalln("failed to open config file")
		return
	}
	confScanner:=bufio.NewScanner(conf)
	for confScanner.Scan(){
		ip=append(ip,confScanner.Text())
	}
	conf.Close()
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
	//read from file
	f,err:=os.Open(filepath)
	if err !=nil{
		panic("failed to open data file")
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
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
	//time.Sleep(0.2)

}
//todo copy or ref
func getClient(key string,c *consistent.Consistent) (*rocksdb_example.RocksdbClient,error){
	 addr,err:=c.Get(key)

	 fmt.Println(key,"at",addr)
	 if err!=nil{
	 	log.Fatalln("failed to get ip address")
	 	return nil,err
	 }
	switch addr {
	case ip[0]:
		cnt0++
	case ip[1]:
		cnt1++
	default:
		cnt2++
	}
	 conn:=connpool[addr]
	 client:=rocksdb_example.NewRocksdbClient(conn)
	 return &client,nil
}
func consistHashing(ip []string) *consistent.Consistent{
	c:=consistent.New(hasher)
	c.NumberOfReplicas=replicas
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
		if	resp.OK{
			log.Println("put", key, value, resp.GetOK())
			cnt++
		}else {
			log.Println("write forbidden")
		}
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
	fmt.Println(cnt0,cnt1,cnt2,cnt)
	//fmt.Println(put0,put1,put2)
}
