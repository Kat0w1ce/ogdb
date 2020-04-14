package main

import (
	"google.golang.org/grpc"
	"log"
	"stathat.com/c/consistent"
)

//type ip struct {
//	address string `json:"IP"`
//	port string `json:"PORT"`
//}
var(
	pool map[string]*grpc.ClientConn
)
func main() {
	ip:=[]string {"localhost:9999","192.168.0.1:2233","localhost:6655"}
	//pool,err:=grpcpool.New(grpc.Dial)
	pool=make(map[string]*grpc.ClientConn)
	for _,addr:=range ip{
		conn,err:=grpc.Dial(addr,grpc.WithInsecure())
		if err!=nil{

			log.Fatalln("unexpect error happened when dial to",addr)
		}
		if conn!=nil {
			log.Println("connect with",addr)
			pool[addr] = conn
		}
	}
	defer clear()
	c:=consistent.New()
	for _,addr:=range ip{
		c.Add(addr)
	}

/*
	todo read config from json

	jsonfile,err:=os.Open(os.Args[1])
	if err!=nil {
		log.Fatal("there's something wrong with configure file")
	}
	defer jsonfile.Close()
	data,err:=ioutil.ReadAll(jsonfile)
	if()
	}

 */
	// grpcpool
	//test data
}

func clear() {
	log.Println("clear")
	if pool != nil {
		for ip, conn := range pool {
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