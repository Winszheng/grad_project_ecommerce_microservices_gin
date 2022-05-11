package main

import (
	"context"
	"fmt"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/user/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	userClient proto.UserClient // 复用连接比较方便
	conn       *grpc.ClientConn
)

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	userClient = proto.NewUserClient(conn)
}

func TestGetUserList() {
	rsp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    1, // 页码从1开始
		PSize: 3,
	})
	if err != nil {
		panic(err)
	}
	for _, user := range rsp.Data {
		fmt.Println("user_info:", user)
		check, err := userClient.CheckPassword(context.Background(), &proto.CheckPasswordInfo{
			Password:         "2333",
			EncrytedPassword: user.Password, // 从数据库取出的密码
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(check.Success) // 橙色的f，含义为field
	}
}

func TestCreateUser() {
	password := "2333"

	for i := 0; i < 10; i++ {
		userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
			Nickname: fmt.Sprintf("user%d", i),
			Password: password,
			Mobile:   fmt.Sprintf("012345689%d", i),
		})
	}
}

func main() {
	Init()
	TestCreateUser()
	//TestGetUserList()

	conn.Close()
}
