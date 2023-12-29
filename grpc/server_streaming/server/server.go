package main

import (
	lotspb "server_streaming/proto"
	"strconv"

	// system
	"math/rand"
	"net"
	"time"

	// external
	log "github.com/alexfmsu/GoLog"
	getopt "github.com/pborman/getopt/v2"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) ActiveLots(req *lotspb.LotsRequest, resp lotspb.LotsService_ActiveLotsServer) error {
	startPrice := rand.Intn(100)

	var i int = 0
	for {
		res := lotspb.LotsResponse{
			Lot: &lotspb.Lot{
				ID:    int64(i),
				Desc:  "Description",
				Price: float64(startPrice * i),
			},
		}
		resp.Send(&res)
		time.Sleep(time.Second * 3)
		if req.Limit < int64(i)+1 {
			break
		}

		i++
	}
	return nil
}

var (
	port int = 0
)

func main() {
	getopt.FlagLong(&port, "port", 'p', "port")
	getopt.Parse()

	listen, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("can't listen on port: %v", err)
	}

	s := grpc.NewServer()
	lotspb.RegisterLotsServiceServer(s, &server{})
	if err := s.Serve(listen); err != nil {
		log.Fatalf("can't register service server: %v", err)
	}

	log.Info("111")
}
