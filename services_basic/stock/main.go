package main

import (
	"flag"
	"fmt"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/stock/global"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/stock/handler"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/stock/proto"
	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var IP = flag.String("ip", "0.0.0.0", "ip地址")

var Port = flag.Int("port", 0, "端口号")

func main() {
	flag.Parse()
	if *Port == 0 {
		*Port, _ = global.GetFreePort()
	}

	global.InitLogger()
	global.InitConfig()
	global.InitDB()

	server := grpc.NewServer()

	proto.RegisterStockServer(server, &handler.StockServer{})

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	zap.S().Infof("*IP:*Port == %s:%d", *IP, *Port)
	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	// 注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	conf := api.DefaultConfig()
	consul := global.ServiceConfig.ConsulInfo
	conf.Address = fmt.Sprintf("%s:%d", consul.Host, consul.Port)

	client, err := api.NewClient(conf)
	if err != nil {
		panic(err)
	}

	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	err = client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		Name:    global.ServiceConfig.Name,
		ID:      serviceID,
		Port:    *Port,
		Tags:    []string{"stock_service", "basic_service"},
		Address: "127.0.0.1",
		Check: &api.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("127.0.0.1:%d", *Port), // grpc服务的地址
			Timeout:                        "5s",
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "15s",
		},
	})

	if err != nil {
		zap.S().Info("registration.Port:", *Port)
		zap.S().Info(err)
	}

	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()

	quit := make(chan os.Signal)                         // 定义终止信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // signal.Notify: 注册要接收的信号
	<-quit                                               // syscall.SIGINT: ctrl-c; syscall.SIGTERM:结束程序
	if err = client.Agent().ServiceDeregister(serviceID); err != nil {
		zap.S().Infof("服务 [%s] 注销失败", serviceID)
		return
	}
	zap.S().Infof("服务 [%s] 注销成功", serviceID)
}
