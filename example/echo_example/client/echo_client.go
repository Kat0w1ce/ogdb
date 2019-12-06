package main

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"ogdb/example/echo_example/demo"
)

const (
	ADDRESS string = "localhost:8081"
)

func main() {
	conn, err := grpc.Dial(ADDRESS, grpc.WithInsecure())
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
