package order

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/Aditya7880900936/microservices_go/account"
	"github.com/Aditya7880900936/microservices_go/catalog"
	"github.com/Aditya7880900936/microservices_go/order/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service       Service
	accountClient *account.Client
	catalogClient *catalog.Client
}

func ListenGRPC(s Service, accountURL, catalogURL string, port int) error {
	accountClient, err := account.NewClient(accountURL)
	if err != nil {
		return err
	}

	catalogClient, err := catalog.NewClient(catalogURL)
	if err != nil {
		accountClient.Close()
		return err
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		catalogClient.Close()
		accountClient.Close()
		return err
	}
	serv := grpc.NewServer()
	pb.RegisterOrderServiceServer(serv, &grpcServer{
		s,
		accountClient,
		catalogClient,
	})
	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) PostOrder(ctx context.Context, r *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {
	_, err := s.accountClient.GetAccount(ctx, r.AccountId)
	if err != nil {
		log.Println("error getting account", err)
		return nil, errors.New("account not found")
	}

	productIDs := []string{}
	orderedProducts, err := s.catalogClient.GetProducts(ctx, 0, 0, productIDs, "")
	if err != nil {
		log.Println("error getting products", err)
		return nil, errors.New("internal error")
	}
	products := []OrderedProduct{}

	for _, p := range orderedProducts {
		product := OrderedProduct{
			ID:          p.ID,
			Quantity:    0,
			Price:       p.Price,
			Name:        p.Name,
			Description: p.Description,
		}
		for _, rp := range r.Products {
			if rp.ProductId == p.ID {
				product.Quantity = rp.Quantity
				break
			}			
		}
		if product.Quantity != 0 {
			products = append(products, product)
		}
	}
	order , err := s.service.PostOrder(ctx , r.AccountId, products)
	if err!= nil {
		log.Println("error posting order", err)
		return nil, errors.New("internal error")
	}

	orderProto := &pb.Order{
		Id:        order.ID,
		AccountId: order.AccountID,
		Products:  []*pb.Order_OrderedProduct{},
		TotalPrice:     order.TotalPrice,
	}
	orderProto.CreatedAt , _ = order.CreatedAt.MarshalBinary()
}
