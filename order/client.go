package order

import (
	"github.com/Aditya7880900936/microservices_go/order/pb"
	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.OrderServiceClient
}


