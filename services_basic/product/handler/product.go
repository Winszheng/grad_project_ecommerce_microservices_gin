package handler

import (
	"context"
	"fmt"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/product/global"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/product/model"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/product/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProductServer struct {
	proto.UnsafeProductServer
}

func Model2Response(product model.Product) proto.ProductInfoResponse {
	return proto.ProductInfoResponse{
		Id:     product.ID,
		ItemID: product.ItemID,
		//Item:
		Name:       product.Name,
		ArticleNum: product.ArticleNum,
		//ClickNum:    product.ClickNum,
		SoldNum:     product.SoldNum,
		FavoriteNum: product.FavoriteNum,
		NormalPrice: product.NormalPrice,
		ProPrice:    product.ProPrice,
		Brief:       product.Brief,
		IsShipFree:  product.IsShipFree,
		FrontImage:  product.FrontImage,
		IsNew:       product.IsNew,
		//IsHot:       product.IsHot,
		OnSale:     product.IsOnSale,
		DescImages: product.DescImages,
		Images:     product.Images,

		Brand: &proto.BrandInfoResponse{
			Id:   product.BrandID,
			Name: product.Brand.Name,
			Logo: product.Brand.LogoUrl,
		},
		Item: &proto.ItemInfoResponse{
			Id:   product.ItemID,
			Name: product.Item.Name,
		},
	}
}

func (s *ProductServer) GetProductByFilter(ctx context.Context, req *proto.ProductByFilterRequest) (*proto.ProductByFilterResponse, error) {

	zap.S().Info("调用[GetProductByFilter]")

	var product []model.Product
	localDB := global.DB.Model(model.Product{})

	if req.Keyword != "" {
		localDB = localDB.Where("name LIKE ?", "%"+req.Keyword+"%")
	}

	if req.MinPrice > 0 {
		localDB = localDB.Where("pro_price >= ?", req.MinPrice)
	}

	if req.MaxPrice > 0 {
		localDB = localDB.Where("pro_price <= ?", req.MaxPrice)
	}

	if req.Brand > 0 {
		localDB = localDB.Where("brand_id=?", req.Brand)
	}

	var sql string
	if req.TopItem > 0 {
		var item model.Item
		if result := global.DB.First(&item, req.TopItem); result.RowsAffected == 0 {
			return nil, status.Error(codes.NotFound, "筛选分类不存在")
		}
		if item.Level == 1 {
			sql = fmt.Sprintf("select id from item where parent_item_id in (select id from item WHERE parent_item_id=%d)", req.TopItem)
		} else if item.Level == 2 {
			sql = fmt.Sprintf("select id from item WHERE parent_item_id=%d", req.TopItem)
		} else if item.Level == 3 {
			sql = fmt.Sprintf("select id from item WHERE id=%d", req.TopItem)
		}
		localDB = localDB.Where(fmt.Sprintf("item_id in (%s)", sql))
	}

	var count int64
	localDB.Count(&count)
	var rsp proto.ProductByFilterResponse
	rsp.Total = int32(count)

	result := localDB.Preload("Brand").Preload("Item").Scopes(Paginate(int(req.Page), int(req.PagePerNum))).Find(&product)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, pro := range product { // 构造返回数据
		zap.S().Infof("pro: %+v\n", pro.Item)
		response := Model2Response(pro)
		rsp.Data = append(rsp.Data, &response)
	}

	return &rsp, nil
}

func (s *ProductServer) GetProductQuantity(ctx context.Context, req *proto.ProductQuantityRequest) (*proto.ProductQuantityResponse, error) {
	var product []model.Product
	result := global.DB.Find(&product, req.Id) // 批量获取model
	var rsp proto.ProductQuantityResponse
	for _, pro := range product {
		p := Model2Response(pro)
		rsp.Data = append(rsp.Data, &p)
	}
	rsp.Total = int32(result.RowsAffected)
	return &rsp, nil
}

func (s *ProductServer) CreateProduct(ctx context.Context, req *proto.CreateProductInfo) (*proto.ProductInfoResponse, error) {
	var item model.Item
	if result := global.DB.First(&item, req.ItemID); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brand
	if result := global.DB.First(&brand, req.BrandID); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	product := model.Product{
		Brand:       brand,
		BrandID:     brand.ID,
		Item:        item,
		ItemID:      item.ID,
		Name:        req.Name,
		ArticleNum:  req.ArticleNum,
		NormalPrice: req.NormalPrice,
		ProPrice:    req.ProPrice,
		Brief:       req.Brief,
		IsShipFree:  req.IsShipFree,
		Images:      req.Images,
		DescImages:  req.DescImages,
		FrontImage:  req.FrontImage,
		IsNew:       req.IsNew,
		//IsHot:       req.IsHot,	这个字段用不上，被我在数据库删掉了，所以gorm就插不进数据库
		IsOnSale: req.IsOnSale,
	}

	//srv之间互相调用了
	tx := global.DB.Begin()
	result := tx.Save(&product)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()

	zap.S().Info("逻辑服务传来的库存为:", req.Stocks)
	_, err := global.StockClient.SetStock(context.Background(), &proto.ProductInfo{
		ProductID: product.ID,
		Num:       req.Stocks,
	})
	if err != nil {
		zap.S().Infof("设置主键为%d的商品库存失败", product.ID)
		return nil, err
	}

	return &proto.ProductInfoResponse{
		Id: product.ID,
	}, nil
}

func (s *ProductServer) DeleteProductByID(ctx context.Context, req *proto.ProductID) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Product{}, req.Id); result.Error != nil {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	return &emptypb.Empty{}, nil
}

func (s *ProductServer) UpdateProduct(ctx context.Context, req *proto.CreateProductInfo) (*emptypb.Empty, error) {
	var product model.Product

	if result := global.DB.First(&product, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}

	var item model.Item
	if result := global.DB.First(&item, req.ItemID); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id为%d的商品分类不存在", req.ItemID)
	}

	var brand model.Brand
	if result := global.DB.First(&brand, req.BrandID); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	product.Brand = brand
	product.BrandID = brand.ID
	product.Item = item
	product.ItemID = item.ID
	product.Name = req.Name
	product.ArticleNum = req.ArticleNum
	product.NormalPrice = req.NormalPrice
	product.ProPrice = req.ProPrice
	product.Brief = req.Brief
	product.IsShipFree = req.IsShipFree
	product.Images = req.Images
	product.DescImages = req.DescImages
	product.FrontImage = req.FrontImage
	product.IsNew = req.IsNew

	product.IsOnSale = req.IsOnSale

	tx := global.DB.Begin() // 为啥要这样
	result := tx.Save(&product)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}
func (s *ProductServer) GetProductDetailByID(ctx context.Context, req *proto.ProductID) (*proto.ProductInfoResponse, error) {
	zap.S().Info("调用[product.GetProductDetailByID]")
	var product model.Product

	if result := global.DB.Preload("Item").First(&product, "id = ?", req.Id); result.RowsAffected == 0 {
		zap.S().Infof("id为%d的商品不存在, 查询到的是%+v", req.Id, product)
		return nil, status.Error(codes.NotFound, fmt.Sprintf("id为%d的商品不存在", req.Id))
	}
	zap.S().Infof("查询到的是%+v", product)
	var rsp proto.ProductInfoResponse
	rsp = Model2Response(product)
	return &rsp, nil
}
