package handler

import (
	"context"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/product/global"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/product/model"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/product/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ProductServer) ItemBrandList(ctx context.Context, req *proto.ItemBrandFilterRequest) (*proto.ItemBrandListResponse, error) {
	var itemBrand []model.ItemBrand
	rsp := proto.ItemBrandListResponse{}

	var total int64
	global.DB.Model(&model.ItemBrand{}).Count(&total)
	rsp.Total = int32(total)

	global.DB.Preload("Item").Preload("Brand").Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&itemBrand)

	var categoryResponses []*proto.ItemBrandResponse
	for _, categoryBrand := range itemBrand {
		categoryResponses = append(categoryResponses, &proto.ItemBrandResponse{
			Category: &proto.ItemInfoResponse{
				Id:    categoryBrand.Item.ID,
				Name:  categoryBrand.Item.Name,
				Level: categoryBrand.Item.Level,
				//IsTab:      categoryBrand.Item.IsTab,
				ParentItem: categoryBrand.Item.ParentItemID,
			},
			Brand: &proto.BrandInfoResponse{
				Id:   categoryBrand.Brand.ID,
				Name: categoryBrand.Brand.Name,
				Logo: categoryBrand.Brand.LogoUrl,
			},
			Id: categoryBrand.ID,
		})
	}

	rsp.Data = categoryResponses
	return &rsp, nil
}
func (s *ProductServer) GetBrandListByItem(ctx context.Context, req *proto.ItemInfoRequest) (*proto.BrandListResponse, error) {
	zap.S().Info("调用[GetBrandListByItem]")
	rsp := proto.BrandListResponse{}

	var category model.Item
	if result := global.DB.Find(&category, req.Id).First(&category); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var categoryBrands []model.ItemBrand
	if result := global.DB.Preload("Brand").Where(&model.ItemBrand{ItemID: req.Id}).Find(&categoryBrands); result.RowsAffected > 0 {
		rsp.Total = int32(result.RowsAffected)
	}

	var brandInfoResponses []*proto.BrandInfoResponse
	for _, categoryBrand := range categoryBrands {
		brandInfoResponses = append(brandInfoResponses, &proto.BrandInfoResponse{
			Id:   categoryBrand.Brand.ID,
			Name: categoryBrand.Brand.Name,
			Logo: categoryBrand.Brand.LogoUrl,
		})
	}

	rsp.Data = brandInfoResponses

	return &rsp, nil
}

func (s *ProductServer) CreateItemBrand(ctx context.Context, req *proto.ItemBrandRequest) (*proto.ItemBrandResponse, error) {
	var item model.Item
	if result := global.DB.First(&item, req.ItemID); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}
	var brand model.Brand
	if result := global.DB.First(&brand, req.BrandID); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}
	itemBrand := model.ItemBrand{
		//Model:   model.Model{},
		ItemID: req.ItemID,
		//Item:    model.Item{},
		BrandID: req.BrandID,
		//Brand:   model.Brand{},	数据库不会给关系另开字段; 中间表靠外键关联
	}
	global.DB.Save(&itemBrand)
	return &proto.ItemBrandResponse{Id: itemBrand.ID}, nil

}
func (s *ProductServer) DeleteItemBrand(ctx context.Context, req *proto.ItemBrandRequest) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.ItemBrand{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌分类不存在")
	}
	return &emptypb.Empty{}, nil
}
func (s *ProductServer) UpdateItemBrand(ctx context.Context, req *proto.ItemBrandRequest) (*emptypb.Empty, error) {
	return nil, nil
}
