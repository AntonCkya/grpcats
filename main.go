package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	grpcats "github.com/AntonCkya/grpcats/grpc"
	grpc "google.golang.org/grpc"
)

type CatsServerImpl struct {
	grpcats.UnimplementedCatsServer
}

func (s *CatsServerImpl) GetCat(request *grpcats.CatRequest, stream grpc.ServerStreamingServer[grpcats.CatResponse]) error {
	says := request.Says
	url := "https://cataas.com/cat/says/"

	clientResponse, err := http.Get(url + says)
	if err != nil {
		return err
	}
	defer clientResponse.Body.Close()

	catFile, err := io.ReadAll(clientResponse.Body)
	if err != nil {
		return err
	}

	response := &grpcats.CatResponse{
		Cat: catFile,
	}
	err = stream.Send(response)
	if err != nil {
		return err
	}
	return nil
}

func newServer() grpcats.CatsServer {
	srv := &CatsServerImpl{}
	return srv
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	grpcats.RegisterCatsServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
