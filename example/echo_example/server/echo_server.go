package main

import (
	"awesomeProject/example/echo_example/demo"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const (
	HOST string="localhost"
	PORT string="8081"
)

type FormatData struct {}

func (fd *FormatData)  Echo(ctx context.Context , in *echo_example.Msg)(out *echo_example.Msg,err error){
	str:= in.Text
	fmt.Println("server: ",str)
	out=&echo_example.Msg{Text: str}
	return  out,nil
}
func main()  {
	listener,err:=net.Listen("tcp",HOST+":"+PORT)
	if err!=nil{
		log.Fatal("failed listen at: localhost")
	}else{
		log.Println("server is listening  ")
	}
	rpcServer:=grpc.NewServer()
	echo_example.RegisterEchoServer(rpcServer,&FormatData{})
	reflection.Register(rpcServer)
	if err=rpcServer.Serve(listener);err!=nil{
		log.Fatal("Error")
	}
}