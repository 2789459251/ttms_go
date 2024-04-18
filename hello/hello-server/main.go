package main

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	pb "hello/hello-server/proto"
	"net"
)

type server struct {
	pb.UnsafeSayHelloServer
}

func (s *server) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("未传token")
	}
	var appId string
	var appKey string
	if v, ok := md["appid"]; ok {
		appId = v[0]
	}
	if v, ok := md["appkey"]; ok {
		appKey = v[0]
	}
	//校验
	if appId != "zy" || appKey != "666" {
		return nil, errors.New("token不正确")
	}
	fmt.Println("hello", request.RequestName)
	return &pb.HelloResponse{ResponseMsg: "hello" + request.RequestName}, nil
}
func main() {
	//TSL认证只需要把pem文件和key文件传入就可以
	//creds, _ := credentials.NewServerTLSFromFile("/usr/local/sdk/go/src/hello/key/test.pem",
	//	"/usr/local/sdk/go/src/hello/key/test.key")

	//开启端口
	listen, _ := net.Listen("tcp", ":9090")
	//创建grpc服务
	grpcServer := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	//注册服务
	pb.RegisterSayHelloServer(grpcServer, &server{})
	//启动服务
	err := grpcServer.Serve(listen)
	if err != nil {
		fmt.Println("failed to server:", err)
		return
	}
}
