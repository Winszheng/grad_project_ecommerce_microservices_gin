package main

import (
	"context"
	"fmt"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/order/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var conn *grpc.ClientConn    // gRPC客户端连接
var client proto.OrderClient // gRPC测试客户端

func Init() { // 初始化gRPC连接
	var err error
	conn, err = grpc.Dial("127.0.0.1:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client = proto.NewOrderClient(conn)
}

func TestCreateCart() {
	_, err := client.CreateCart(context.Background(), &proto.CartRequest{
		UserID:    1,
		ProductID: 113,
		Num:       6,
	})
	if err != nil {
		fmt.Println("创建购物车条目失败")
	}
}

func main() {
	Init()
	TestCreateCart()
}

func TestDeleteCart() {
	client.DeleteCart(context.Background(), &proto.CartRequest{UserID: 1, ProductID: 113})
}

func TestUpdateCart() {

	client.UpdateCart(context.Background(), &proto.CartRequest{
		UserID:    1,
		ProductID: 113,
		Num:       10,
		Checked:   true,
	})
}

func TestGetCartList() {
	rsp, _ := client.GetCartList(context.Background(), &proto.UserInfo{Id: 1})
	for _, cart := range rsp.Data {
		fmt.Printf("%+v\n", cart)
	}
}

func TestCreateOrder() {
	client.CreateOrder(context.Background(), &proto.OrderRequest{
		UserID:  1,
		Address: "广东省广州市番禺大学城小谷围街道 中山大学明德园六号",
		Name:    "李小狼",
		Mobile:  "18576808293",
		Comment: "不要放丰巢谢谢",
	})
}

func TestGetOrderDetail() {
	rsp, _ := client.GetOrderDetail(context.Background(), &proto.OrderRequest{
		OrderID: 1,
	})
	fmt.Println(rsp)
}

func TestGetOrderList() {
	rsp, _ := client.GetOrderList(context.Background(), &proto.OrderFilterRequest{
		UserID:     1,
		PagePerNum: 1,
		Page:       0,
	})
	fmt.Println(rsp)
}

//func main() {
//	Init()
//	//TestCreateCart()
//	//TestDeleteCart() 成功删除购物车条目，赞！
//	//TestUpdateCart()
//	//TestGetCartList()
//	//TestCreateOrder()
//	//TestGetOrderDetail()
//	TestGetOrderList()
//}
