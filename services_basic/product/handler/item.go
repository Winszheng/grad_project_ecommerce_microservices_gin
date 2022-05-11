package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/product/global"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/product/model"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/product/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ProductServer) GetAllItemList(ctx context.Context, e *emptypb.Empty) (*proto.ItemListResponse, error) {
	zap.S().Info("调用[GetAllItemList]")
	var items []model.Item

	result := global.DB.Where("level=?", 1).Preload("Sub.Sub").Find(&items)
	b, _ := json.Marshal(items)

	return &proto.ItemListResponse{
		Total:    int32(result.RowsAffected),
		JsonData: string(b),
	}, nil
}

func (s *ProductServer) GetSubItem(ctx context.Context, req *proto.ItemListRequest) (*proto.SubItemListResponse, error) {
	zap.S().Info("调用[GetSubItem]")

	rsp := proto.SubItemListResponse{}

	var item model.Item
	if result := global.DB.First(&item, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}

	rsp.Info = &proto.ItemInfoResponse{
		Id:         item.ID,
		Name:       item.Name,
		Level:      item.Level,
		ParentItem: item.ParentItemID,
	}

	var subCategorys []model.Item
	var subCategoryResponse []*proto.ItemInfoResponse

	global.DB.Where(&model.Item{ParentItemID: req.Id}).Find(&subCategorys)
	for _, subItem := range subCategorys {
		subCategoryResponse = append(subCategoryResponse, &proto.ItemInfoResponse{
			Id:         subItem.ID,
			Name:       subItem.Name,
			Level:      subItem.Level,
			ParentItem: subItem.ParentItemID,
		})
	}

	rsp.SubItem = subCategoryResponse
	return &rsp, nil
}

func (s *ProductServer) CreateItem(ctx context.Context, req *proto.ItemInfoRequest) (*proto.ItemInfoResponse, error) {
	zap.S().Info("调用【CreateItem】")
	if result := global.DB.Where("name=?", req.Name).Find(&model.Item{}); result.RowsAffected != 0 {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("分类【%s】已存在", req.Name))
	}

	item := model.Item{
		Name:  req.Name,
		Level: req.Level,
	}

	if req.Level != 1 {
		item.ParentItemID = req.ParentItem
	}

	global.DB.Save(&item)

	zap.S().Infof("创建新分类: %+v\n", item)

	return &proto.ItemInfoResponse{
		Id: item.ID,
	}, nil
}

func (s *ProductServer) DeleteItem(ctx context.Context, req *proto.DeleteItemRequest) (*emptypb.Empty, error) {
	zap.S().Info("调用[DeleteItem]")

	if result := global.DB.Delete(&model.Item{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}

	return &emptypb.Empty{}, nil
}

// UpdateItem 更新要先把东西查出来
func (s *ProductServer) UpdateItem(ctx context.Context, req *proto.ItemInfoRequest) (*emptypb.Empty, error) {
	zap.S().Info("调用【UpdateItem】")

	var item model.Item

	if result := global.DB.First(&item, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "item分类不存在")
	}

	if req.Name != "" {
		item.Name = req.Name
	}

	if req.ParentItem != 0 {
		item.ParentItemID = req.ParentItem
	}

	if req.Level != 0 {
		item.Level = req.Level
	}

	global.DB.Save(item)

	return &emptypb.Empty{}, nil
}
