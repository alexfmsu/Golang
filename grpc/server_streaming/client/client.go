package main

import (
	lotspb "server_streaming/proto"
	"strconv"

	// system
	"context"
	"io"

	// external
	log "github.com/alexfmsu/GoLog"
	getopt "github.com/pborman/getopt/v2"
	"google.golang.org/grpc"
)

var (
	port int = 0
)

func main() {
	getopt.FlagLong(&port, "port", 'p', "port")
	getopt.Parse()

	conn, err := grpc.Dial("localhost:"+strconv.Itoa(port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can't connect to server:")
		return
	}

	defer conn.Close()
	client := lotspb.NewLotsServiceClient(conn)

	req := lotspb.LotsRequest{
		Limit: 30,
	}

	resp, err := client.ActiveLots(context.Background(), &req)
	if err != nil {
		log.Fatalf("can't get active lots: %v", err)
		return
	}

	for {
		lots, err := resp.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("can't receive from server: %v", err)
		}
		log.Infof("active lot: %+v\n", lots.Lot)
	}
}
