package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "hello/hello-server/proto"
	"log"
)

type PerRPCCredentials interface {
	GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error)
	RequireTransportSecurity() bool
}
type ClientTokenAuth struct {
}

func (c ClientTokenAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"appId":  "zy",
		"appkey": "666",
	}, nil
}
func (c ClientTokenAuth) RequireTransportSecurity() bool {
	return false
}

func main() {
	//creds, err := credentials.NewClientTLSFromFile("/usr/local/sdk/go/src/hello/key/test.pem", "*.zy.com")
	//if err != nil {
	//	fmt.Println("creds_err:", err)
	//}
	//通过grpc连接到server端
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithPerRPCCredentials(new(ClientTokenAuth)))
	conn, err := grpc.Dial("127.0.0.1:9090", opts...)
	if err != nil {
		log.Fatalf("did not connect:%v", err)
	}
	defer conn.Close()
	//与客户端连接
	client := pb.NewSayHelloClient(conn)
	//执行rpc调用（在服务端实现并返回结果）
	resp, _ := client.SayHello(context.Background(), &pb.HelloRequest{RequestName: "zy"})
	fmt.Println(resp.GetResponseMsg())
}
