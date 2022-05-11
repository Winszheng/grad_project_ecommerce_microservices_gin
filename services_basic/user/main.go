package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/user/global"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/user/handler"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/user/proto"
	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
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

	// 1.创建grpc server
	server := grpc.NewServer()

	proto.RegisterUserServer(server, &handler.UserServer{})

	//3.建立监听
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	zap.S().Infof("*IP:*Port == %s:%d", *IP, *Port)
	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	conf := api.DefaultConfig()
	consul := global.ServiceConfig.ConsulInfo
	conf.Address = fmt.Sprintf("%s:%d", consul.Host, consul.Port) // 127.0.0.1:8500
	// Create a Consul API client
	client, err := api.NewClient(conf)
	if err != nil {
		panic(err)
	}

	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	err = client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		Name:    global.ServiceConfig.Name,
		ID:      serviceID,
		Port:    *Port,
		Tags:    []string{"Yuno", "user", "basic_service"},
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
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) 
	<-quit                                               // syscall.SIGINT: ctrl-c; syscall.SIGTERM:结束程序
	if err = client.Agent().ServiceDeregister(serviceID); err != nil {
		zap.S().Infof("服务 [%s] 注销失败", serviceID)
		return
	}
	zap.S().Infof("服务 [%s] 注销成功", serviceID)
}
