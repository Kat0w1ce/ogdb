package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	echo_example "ogdb/example/echo_example/demo"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	HOST string = "0.0.0.0"
	PORT string = "2233"
	MAX  int32	=10000
)

type FormatData struct{}

func (fd *FormatData) Echo(ctx context.Context, in *echo_example.Msg) (out *echo_example.Msg, err error) {
	str := in.Text
	fmt.Println("server: ", str)
	out = &echo_example.Msg{Text: str}
	return out, nil
}
func main() {
	flag.StringVar(&HOST,"h","0.0.0.0","ip")
	flag.StringVar(&PORT,"p","9999","port")
	flag.IntVar(&MAX,"m",10000,"max cap")
	flag.Parse()
	address:=fmt.Sprint(HOST,":",PORT)
	fmt.Println("litsen at ",address)
	listener, err := net.Listen("tcp", HOST+":"+PORT)
	if err != nil {
		log.Fatal("failed listen at: localhost")
	} else {
		log.Println("server is listening  ",address)
	}
	rpcServer := grpc.NewServer()
	echo_example.RegisterEchoServer(rpcServer, &FormatData{})
	reflection.Register(rpcServer)
	if err = rpcServer.Serve(listener); err != nil {
		log.Fatal("Error")
	}
}

