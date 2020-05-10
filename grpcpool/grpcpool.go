package grpcpool

import (
	"google.golang.org/grpc"
	"sync"
)



type Factory func() (*grpc.ClientConn,error)
type GrpcPool struct {
	conns chan *grpc.ClientConn
	factory Factory
}

func NewGrpcPool(address string,cap int)  *GrpcPool{
	return &GrpcPool{
		conns: make(chan *grpc.ClientConn,cap),
		factory: func() (*grpc.ClientConn,error) {
		res,err:=grpc.Dial(address)
		return res,err
		},
	}
}

func (p *GrpcPool) new() (*grpc.ClientConn,error){
	return p.factory()
}

func (p *GrpcPool) Get() (conn *grpc.ClientConn) {
	select {
		case conn = <- p.conns:{}
	default:
		conn,_=p.new()
	}
	return
}

func (p *GrpcPool) Put(conn *grpc.ClientConn)  {
	select {
	case p.conns<- conn:{}
	default:
		conn.Close()
	}
}