package global

import (
	"fmt"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/user/proto"
	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitGrpcClient() {
	consul := ServiceConfig.ConsulInfo
	userService := ServiceConfig.UserSrvInfo
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",
			consul.Host, consul.Port, userService.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitGrpcClient] 连接 【用户服务失败】")
	}
	UserSrvClient = proto.NewUserClient(userConn)
}

func InitGrpcClient2() {

	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", ServiceConfig.ConsulInfo.Host,
		ServiceConfig.ConsulInfo.Port)

	userSrvHost := ""
	userSrvPort := 0

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`,
		ServiceConfig.UserSrvInfo.Name))

	if err != nil {
		panic(err)
	}

	for _, value := range data {
		userSrvHost = value.Address
		userSrvPort = value.Port
		break
	}

	if userSrvHost == "" {
		zap.S().Fatal("[InitGrpcClient] 连接 【用户服务失败】")
		return
	}

	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接【用户服务】失败", "msg", err.Error())
	}

	UserSrvClient = proto.NewUserClient(userConn)

}
