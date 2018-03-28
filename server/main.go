package main

import (
	"fmt"
	"gaoyl/linkDb"
	pb "gaoyl/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const (
	port = ":50052"
)

type netConfig struct{}

func (s *netConfig) RequestConfigLink(ctx context.Context, inpara *pb.PortParameter) (*pb.PortParameter, error) {
	err := linkDb.AddToUnLinkedPort(*inpara)
	fmt.Println(inpara)
	if err != nil {
		return &pb.PortParameter{}, err
	}
	state := linkDb.WaitPortBeConfiged(linkDb.PortKey{inpara.AppName, inpara.PortName})
	if state != linkDb.LINKCONFIGOK {
		return inpara, fmt.Errorf("配置链路失败%d:\r\b", state)
	}
	return &pb.PortParameter{}, nil
}
func (s *netConfig) QueryAllLinks(ctx context.Context, nodeInfo *pb.NodeInfo) (*pb.Links, error) {
	links := make([]*pb.LinkParameter, 0)
	_, linkInfo := linkDb.GetAllInfo()
	i := 0
	for _, v := range linkInfo {
		links = append(links, &v)
		i = i + 1
		fmt.Println(links)
	}
	return &pb.Links{Lp: links}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterNetConfigServer(s, &netConfig{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
