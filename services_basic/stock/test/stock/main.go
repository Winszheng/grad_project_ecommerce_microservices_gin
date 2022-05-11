package main

import (
	"context"
	"fmt"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/stock/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"strconv"
	"time"
)

var stockClient proto.StockClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.Dial(fmt.Sprintf("%s:%d", "127.0.0.1", 50053),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	stockClient = proto.NewStockClient(conn)
}

func main() {
	Init()
	//_, err := stockClient.SetStock(context.Background(), &proto.ProductInfo{
	//	ProductID: 113,
	//	Num:       1000,
	//})
	//if err != nil {
	//	fmt.Println("调用[SetStock]有误:", err)
	//}
	//rsp, _ := stockClient.GetStock(context.Background(), &proto.ProductInfo{
	//	ProductID: 113,
	//})
	//if rsp.Num == 1000 {
	//	fmt.Println("[GetStock]正常工作")
	//}

	//批量库存扣减：因为我懒得生成订单号了，先随便整整吧

	for i := 0; i < 50; i++ {
		go func() {
			stockClient.Deduction(context.Background(), &proto.OrderInfo{
				OrderNo: strconv.Itoa(i),
				ProductList: []*proto.ProductInfo{
					{
						ProductID: 113,
						Num:       1,
					},
				},
			})
		}()
	}

	time.Sleep(5 * time.Second)
}
