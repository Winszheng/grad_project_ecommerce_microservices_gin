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
	"gorm.io/gorm"
)

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func (s *ProductServer) BrandList(ctx context.Context, req *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	zap.S().Info("调用[handler.BrandList]")
	var rsp proto.BrandListResponse
	var brand []model.Brand
	result := global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&brand)
	if result.Error != nil {
		return nil, result.Error
	}

	var total int64
	global.DB.Model(&model.Brand{}).Count(&total)
	rsp.Total = int32(total)

	for _, b := range brand {
		rsp.Data = append(rsp.Data, &proto.BrandInfoResponse{ // 没毛病
			Id:   b.ID,
			Name: b.Name,
			Logo: b.LogoUrl,
		})
	}

	return &rsp, nil
}

func (s *ProductServer) CreateBrand(ctx context.Context, req *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	brand := model.Brand{}
	if result := global.DB.Where("name=?", req.Name).Find(&brand); result.RowsAffected != 0 {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("品牌【%s】已存在", req.Name))
	}

	brand = model.Brand{
		Name:    req.Name,
		LogoUrl: req.Logo,
	}

	global.DB.Save(&brand)

	return &proto.BrandInfoResponse{
		Id:   brand.ID,
		Name: brand.Name,
		Logo: brand.LogoUrl,
	}, nil

}

func (s *ProductServer) DeleteBrand(ctx context.Context, req *proto.BrandRequest) (*emptypb.Empty, error) {
	// 根据 表名&id 删除数据
	if result := global.DB.Delete(&model.Brand{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("品牌【%s】不存在", req.Name))
	}
	return &emptypb.Empty{}, nil
}

// UpdateBrand 更新品牌
func (s *ProductServer) UpdateBrand(ctx context.Context, req *proto.BrandRequest) (*emptypb.Empty, error) {
	brand := model.Brand{} // 因为save会更新整条数据，为了避免没有更新的字段被覆盖，所以要先把整条数据查出来，修改完再save回去
	if result := global.DB.Where("name=?", req.Name).Find(&brand); result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("品牌【%s】不存在, 无法更新", req.Name))
	}

	if req.Name != "" {
		brand.Name = req.Name
	}

	if req.Logo != "" {
		brand.LogoUrl = req.Logo
		zap.S().Info("更新logo为%s", req.Logo)
	}

	global.DB.Save(&brand)

	return &emptypb.Empty{}, nil
}

func (s *ProductServer) GetBrand(ctx context.Context, req *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	brand := model.Brand{} // 读取数据到model
	if result := global.DB.Where("id=?", req.Id).Find(&brand); result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("品牌【%s】不存在, 无法更新", req.Name))
	}
	rsp := proto.BrandInfoResponse{
		Id:   brand.ID,
		Name: brand.Name,
		Logo: brand.LogoUrl,
	}
	return &rsp, nil
}
