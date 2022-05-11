package handler

import (
	"context"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/order/global"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/order/model"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/order/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type OrderServer struct {
	proto.UnimplementedOrderServer
}

func (OrderServer) CreateCart(ctx context.Context, req *proto.CartRequest) (*proto.CartResponse, error) {
	var cart model.Cart

	if result := global.DB.Where(&model.Cart{ProductID: req.ProductID, UserID: req.UserID}).First(&cart); result.RowsAffected != 0 {
		cart.Num += req.Num // 合并数量
	} else {
		cart.UserID = req.UserID
		cart.ProductID = req.ProductID
		cart.Num = req.Num
		cart.Checked = false
	}
	global.DB.Save(&cart)
	return &proto.CartResponse{
		CartID: cart.ID,
	}, nil

}

func (OrderServer) DeleteCart(ctx context.Context, req *proto.CartRequest) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Cart{}, model.Cart{ProductID: req.ProductID, UserID: req.UserID}); result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "待删除的商品条目不存在")
	}
	return &emptypb.Empty{}, nil
}

func (OrderServer) UpdateCart(ctx context.Context, req *proto.CartRequest) (*emptypb.Empty, error) {
	var cart model.Cart
	if result := global.DB.First(&cart, &model.Cart{UserID: req.UserID, ProductID: req.ProductID}); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "待更新待条目不存在")
	}
	if req.Num > 0 {
		cart.Num = req.Num
	}

	cart.Checked = req.Checked
	global.DB.Save(&cart)
	return &emptypb.Empty{}, nil
}

func (OrderServer) GetCartList(ctx context.Context, req *proto.UserInfo) (*proto.CartListResponse, error) {
	var cartList []model.Cart
	var rsp proto.CartListResponse
	if result := global.DB.Find(&cartList, model.Cart{UserID: req.Id}); result.Error != nil {
		return nil, status.Errorf(codes.NotFound, "查找购物车list出错")
	} else {
		rsp.Total = int32(result.RowsAffected)
	}
	for _, cart := range cartList {

		product, err :=
			global.ProductClient.GetProductDetailByID(context.Background(), &proto.ProductID{Id: cart.ProductID})
		if err != nil {
			return nil, status.Errorf(codes.NotFound, "查找商品出错")
		}
		rsp.Data = append(rsp.Data, &proto.CartResponse{
			CartID:    cart.ID,
			UserID:    cart.UserID,
			ProductID: cart.ProductID,
			Num:       cart.Num,
			Checked:   cart.Checked,

			ProductName: product.Name,
			Price:       product.ProPrice,
			Image:       product.FrontImage,
		})
	}
	return &rsp, nil
}

func (OrderServer) CreateOrder(ctx context.Context, req *proto.OrderRequest) (*proto.OrderResponse, error) {

	var cartList []model.Cart
	if result := global.DB.Find(&cartList, model.Cart{Checked: true}); result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "购物车未勾选下单商品")
	}

	orderNo := GenerateNo(req.UserID)
	var productList []*proto.ProductInfo

	for _, cart := range cartList {
		productList = append(productList, &proto.ProductInfo{
			ProductID: cart.ProductID,
			Num:       cart.Num,
		})
	}
	_, err := global.StockClient.Deduction(context.Background(), &proto.OrderInfo{ // 果然应该先写出来再查
		OrderNo:     orderNo,
		ProductList: productList,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "批量扣除库存失败:", err)
	}

	var amount float32 = 0
	var productDetailList []*proto.ProductInfoResponse
	for _, cart := range cartList {
		product, err :=
			global.ProductClient.GetProductDetailByID(context.Background(), &proto.ProductID{Id: cart.ProductID})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "获取商品金额失败:", err)
		}
		product.SoldNum = cart.Num
		productDetailList = append(productDetailList, product) // 得到商品详情list
		amount += product.ProPrice * float32(cart.Num)
	}

	createTime := time.Now()
	rsp := proto.OrderResponse{
		UserID:     req.UserID,
		No:         orderNo,
		Status:     "WAIT",
		Comment:    req.Comment,
		Amount:     amount,
		Address:    req.Address,
		Name:       req.Name,
		Mobile:     req.Mobile,
		CreateTime: createTime.String(),
	}
	order := model.Order{
		UserID: req.UserID,
		No:     orderNo,
		Status: "WAIT",

		Amount: amount,

		Address: req.Address,
		Name:    req.Name,
		Mobile:  req.Mobile,
		Comment: req.Comment,
	}
	global.DB.Save(&order)

	for _, pro := range productDetailList {
		item := model.OrderItem{
			OrderID:     order.ID,
			UserID:      req.UserID,
			ProductName: pro.Name,
			Image:       pro.FrontImage,
			Price:       pro.ProPrice,
			Num:         pro.SoldNum,
		}
		global.DB.Save(&item) // 把商品条目写进数据库
	}

	for _, cart := range cartList {
		global.DB.Delete(&cart)
	}
	return &rsp, nil
}

func (OrderServer) GetOrderDetail(cxt context.Context, req *proto.OrderRequest) (*proto.OrderDetailResponse, error) {
	var order model.Order

	if result := global.DB.Where(&model.Order{Model: model.Model{ID: req.OrderID}, UserID: req.UserID}).First(&order); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "找不到订单号为%d的订单", req.OrderID)
	}

	rsp := proto.OrderDetailResponse{
		OrderBasicInfo: &proto.OrderResponse{
			Id: req.OrderID,

			UserID:     order.UserID,
			No:         order.No,
			Status:     order.Status,
			Comment:    order.Comment,
			Amount:     order.Amount,
			Address:    order.Address,
			Name:       order.Name,
			Mobile:     order.Mobile,
			CreateTime: order.CreatedAt.String(), // 获取订单创建时间
		},
	}

	var orderGoods []model.OrderItem
	if result := global.DB.Where(&model.OrderItem{OrderID: order.ID}).Find(&orderGoods); result.Error != nil {
		return nil, result.Error
	}

	for _, orderGood := range orderGoods {
		rsp.ProductList = append(rsp.ProductList, &proto.OrderItemResponse{
			GoodsId:    orderGood.ID,
			GoodsName:  orderGood.ProductName,
			GoodsPrice: orderGood.Price,
			GoodsImage: orderGood.Image,
			Nums:       orderGood.Num,
		})
	}

	return &rsp, nil
}

func (OrderServer) GetOrderList(ctx context.Context, req *proto.OrderFilterRequest) (*proto.OrderListResponse, error) {
	var rsp proto.OrderListResponse
	var total int64
	global.DB.Where(&model.Order{UserID: req.UserID}).Count(&total)
	rsp.Total = int32(total)

	var orderList []model.Order

	for _, order := range orderList {
		rsp.Data = append(rsp.Data, &proto.OrderResponse{
			Id:         order.ID,
			UserID:     order.UserID,
			No:         order.No,
			Status:     order.Status,
			Comment:    order.Comment,
			Amount:     order.Amount,
			Address:    order.Address,
			Name:       order.Name,
			Mobile:     order.Mobile,
			CreateTime: order.CreatedAt.String(),
		})
	}
	return &rsp, nil
}

func (OrderServer) UpdateOrderStatus(ctx context.Context, req *proto.OrderStatusRequest) (*emptypb.Empty, error) {
	zap.S().Info("调用[order.UpdateOrderStatus]")

	var order model.Order
	if result := global.DB.Where("no = ?", req.OrderNo).First(&order); result.RowsAffected == 0 {
		zap.S().Infof("编号为%s的订单不存在", req.OrderNo)
		return nil, status.Errorf(codes.Internal, "编号为%s的订单不存在", req.OrderNo)
	}
	order.Status = req.Status
	global.DB.Save(&order)

	return &emptypb.Empty{}, nil
}
