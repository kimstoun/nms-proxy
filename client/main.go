package main

import (
	"flag"
	"fmt"
	pb "gaoyl/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

const (
	address = "localhost:50052"
)

var (
	R              bool
	Q              bool
	appname        string
	portname       string
	rioid          int
	slotsize       int
	porttype       int
	remoteappname  string
	remoteportname string
)

func init() {
	flag.BoolVar(&R, "R", false, "配置链路")
	flag.BoolVar(&Q, "Q", false, "获取链路")
	flag.StringVar(&appname, "appname", "error", "在配置链路参数下的app名称")
	flag.StringVar(&portname, "portname", "error", "在配置链路参数下的port名称")
	flag.StringVar(&remoteappname, "remoteappname", "error", "在配置链路参数下远端的app名称")
	flag.StringVar(&remoteportname, "remoteportname", "error", "在配置链路参数下远端的app名称")
	flag.IntVar(&rioid, "rioid", 25, "在配置链路参数下的本地rioid")
	flag.IntVar(&slotsize, "slotsize", 0, "在配置链路参数下的本地slotsize")
	flag.IntVar(&porttype, "porttype", 4, "在配置链路参数下的本地porttype  0发送 1接受 2 双向")

}
func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	flag.Parse()
	defer conn.Close()
	c := pb.NewNetConfigClient(conn)

	// Contact the server and print out its response.
	if R == true {
		if appname == "error" {
			fmt.Println("参数输入错1误\r\n")
		}
		inpara := &pb.PortParameter{int32(rioid), appname, portname, int32(slotsize), pb.PortType(porttype), remoteappname, remoteportname}
		_, err := c.RequestConfigLink(context.Background(), inpara)
		if err != nil {
			log.Fatalf("could not Config: %v", err)
		}
	} else if Q == true {

		fmt.Println("参数输入错误2\r\n")
	} else {
		fmt.Println("参数输入错22222\r\n")
	}

}
