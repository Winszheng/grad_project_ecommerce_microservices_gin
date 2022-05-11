package handler

import (
	"context"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/stock/global"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/stock/model"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/stock/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm/clause"
)

type StockServer struct {
	proto.UnimplementedStockServer
}

func (StockServer) GetStock(ctx context.Context, req *proto.ProductInfo) (*proto.ProductInfo, error) {
	var stock model.Stock
	if result := global.DB.Where(&model.Stock{ProductID: req.ProductID}).First(&stock); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "id为%d的商品不存在", req.ProductID)
	}
	return &proto.ProductInfo{
		ProductID: req.ProductID,
		Num:       stock.Num,
	}, nil
}

func (StockServer) SetStock(ctx context.Context, req *proto.ProductInfo) (*emptypb.Empty, error) {
	tx := global.DB.Begin()
	var stock model.Stock
	if result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.Stock{ProductID: req.ProductID}).First(&stock); result.RowsAffected == 0 {

	}
	stock.ProductID = req.ProductID
	stock.Num = req.Num
	tx.Where("product_id = ?", req.ProductID).Save(&stock)
	tx.Commit()

	return &emptypb.Empty{}, nil
}

func (StockServer) Deduction(ctx context.Context, req *proto.OrderInfo) (*emptypb.Empty, error) {
	order := model.Order{
		No: req.OrderNo, // 订单编号
	}

	tx := global.DB.Begin()
	for _, product := range req.ProductList {

		var stock model.Stock
		if result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.Stock{ProductID: product.ProductID}).First(&stock); result.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.NotFound, "id为%d的商品不存在", product.ProductID)
		}

		if stock.Num >= product.Num {
			stock.Num -= product.Num
		} else {
			tx.Rollback()
			return nil, status.Errorf(codes.NotFound, "商品库存不合法")
		}
		tx.Where("product_id = ?", product.ProductID).Save(&stock)
		order.Detail = append(order.Detail, stock)
	}
	order.Status = 1
	tx.Where("no = ?", order.No).Save(&order) // 作为扣减过的记录，备归还使用
	tx.Commit()                               // 提交事务

	return &emptypb.Empty{}, nil
}
func (StockServer) Back(ctx context.Context, req *proto.OrderInfo) (*emptypb.Empty, error) {
	var order model.Order
	if result := global.DB.Where(&model.Order{No: req.OrderNo}).First(&order); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "单号为%s的订单不存在", req.OrderNo)
	}
	if order.Status == 2 {
		return nil, status.Errorf(codes.InvalidArgument, "单号为%s的订单已归还", req.OrderNo)
	}
	tx := global.DB.Begin()
	for _, product := range req.ProductList {
		var stock model.Stock
		if result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.Stock{ProductID: product.ProductID}).First(&stock); result.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.NotFound, "id为%d的商品不存在", product.ProductID)
		}
		stock.Num += product.Num

		tx.Where("product_id = ?", product.ProductID).Save(&stock)
	}
	order.Status = 2
	tx.Where("no = ?", order.No).Save(&order)
	tx.Commit()
	return &emptypb.Empty{}, nil
}
