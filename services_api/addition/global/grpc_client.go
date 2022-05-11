package global

import (
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/addition/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var conn *grpc.ClientConn

func InitGrpcClient() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	Client = proto.NewAddressClient(conn)
}
