package global

import (
	"fmt"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/product/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitGrpcClient() {
	consul := ServiceConfig.ConsulInfo
	userService := ServiceConfig.ProductServiceInfo
	productConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",
			consul.Host, consul.Port, userService.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitGrpcClient] 连接 【商品服务失败】")
	}
	ProductClient = proto.NewProductClient(productConn)
}
