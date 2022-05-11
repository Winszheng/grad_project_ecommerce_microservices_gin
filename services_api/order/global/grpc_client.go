package global

import (
	"fmt"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/order/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitGrpcClient() {
	consul := ServiceConfig.ConsulInfo

	orderService := ServiceConfig.OrderSrvConfig
	orderConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",
			consul.Host, consul.Port, orderService.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitGrpcClient] 连接 [订单&购物车服务] 失败")
	}
	OrderClient = proto.NewOrderClient(orderConn)

	productService := ServiceConfig.ProductInfo
	productConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",
			consul.Host, consul.Port, productService.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitGrpcClient] 连接 [商品服务失败]:", err)
	} else {
		zap.S().Info("成功连接商品服务")
	}
	ProductClient = proto.NewProductClient(productConn)

	stockService := ServiceConfig.StockInfo
	stockConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",
			consul.Host, consul.Port, stockService.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitGrpcClient] 连接 【库存服务失败】")
	} else {
		zap.S().Info("成功连接库存服务")
	}
	StockClient = proto.NewStockClient(stockConn)
}
