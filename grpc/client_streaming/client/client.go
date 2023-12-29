package main

import (
	orderspb "client_streaming/proto"

	// system
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	// external
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:5000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can't connect to server: %v", err)
	}
	defer conn.Close()
	client := orderspb.NewOrdersServiceClient(conn)

	stream, err := client.PostOrder(context.Background())
	if err != nil {
		log.Fatalf("can't post order: %v", err)
	}

	for i := 0; i < 10; i++ {
		order := orderspb.OrderRequest{
			Price:    float64(rand.Intn(1000)),
			Quantity: int64(rand.Intn(10)),
		}
		stream.Send(&order)
		time.Sleep(time.Second)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("can't get response from server: %v", err)
	}
	fmt.Printf("executed orders: %d", resp.ExecutedOrders)
}
