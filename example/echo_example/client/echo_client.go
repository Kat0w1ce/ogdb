package main

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"ogdb/example/echo_example/demo"
	"flag"
)

var (
	ADDRESS string = "localhost"
	PORT string ="8081"
)

func main() {
	flag.StringVar(&ADDRESS,"a","localhost","address")
	flag.StringVar(&PORT,"p","2233","address")
	flag.Parse()
	conn, err := grpc.Dial(ADDRESS+":"+PORT, grpc.WithInsecure())
	if err != nil {
		log.Fatal("Can't connect: " + ADDRESS)
	}
	defer conn.Close()
	client := echo_example.NewEchoClient(conn)
	resp, err := client.Echo(context.Background(), &echo_example.Msg{Text: "hello,world!"})
	if err != nil {
		log.Fatal("Echo error" + err.Error())
	}
	log.Println(resp.Text)
}
