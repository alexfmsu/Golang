package main

import (
	orderspb "client_streaming/proto"

	// system
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"

	// external
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) PostOrder(stream orderspb.OrdersService_PostOrderServer) error {
	var executedOrders int
	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				resp := orderspb.OrdersResponse{
					ExecutedOrders: int64(executedOrders),
				}
				stream.SendAndClose(&resp)
				break
			}
			log.Fatalf("can't receive message from client: %v", err)
		}
		if rand.Intn(1000)%2 == 0 {
			fmt.Printf("executed order with price %.2f and quantity %d\n", req.Price, req.Quantity)
			executedOrders++
		}
	}
	return nil
}

func main() {
	listen, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatalf("can't listen on port: %v", err)
	}

	s := grpc.NewServer()
	orderspb.RegisterOrdersServiceServer(s, &server{})
	if err := s.Serve(listen); err != nil {
		log.Fatalf("can't register service server: %v", err)
	}
}
