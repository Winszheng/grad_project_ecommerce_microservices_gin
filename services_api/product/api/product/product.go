package product

import (
	"context"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/product/api"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/product/form"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/product/global"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/product/proto"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func FilterList(ctx *gin.Context) {
	zap.S().Info("调用了[product.FilterList]")
	req := &proto.ProductByFilterRequest{}

	priceMin := ctx.DefaultQuery("pmin", "0")
	priceMinInt, _ := strconv.Atoi(priceMin)
	req.MinPrice = int32(priceMinInt)

	priceMax := ctx.DefaultQuery("pmax", "0")
	priceMaxInt, _ := strconv.Atoi(priceMax)
	req.MaxPrice = int32(priceMaxInt)

	isHot := ctx.DefaultQuery("ih", "0")
	if isHot == "1" {
		req.IsHot = true
	}
	isNew := ctx.DefaultQuery("in", "0")
	if isNew == "1" {
		req.IsNew = true
	}

	isTab := ctx.DefaultQuery("it", "0")
	if isTab == "1" {
		req.IsTab = true
	}

	categoryId := ctx.DefaultQuery("c", "0")
	categoryIdInt, _ := strconv.Atoi(categoryId)
	req.TopItem = int32(categoryIdInt)

	pages := ctx.DefaultQuery("p", "0")
	pagesInt, _ := strconv.Atoi(pages)
	req.Page = int32(pagesInt)

	perNums := ctx.DefaultQuery("pnum", "0")
	perNumsInt, _ := strconv.Atoi(perNums)
	req.PagePerNum = int32(perNumsInt)

	keywords := ctx.DefaultQuery("q", "")
	req.Keyword = keywords

	brandId := ctx.DefaultQuery("b", "0")
	brandIdInt, _ := strconv.Atoi(brandId)
	req.Brand = int32(brandIdInt)

	rsp, err := global.ProductClient.GetProductByFilter(context.WithValue(context.Background(), "ginContext", ctx), req)
	if err != nil {
		zap.S().Errorw("[List] 查询 【商品列表】失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	//e.Exit()
	reMap := map[string]interface{}{
		"total": rsp.Total,
	}

	productList := make([]interface{}, 0)
	for _, value := range rsp.Data {
		productList = append(productList, map[string]interface{}{
			"id":    value.Id,
			"name":  value.Name,
			"brief": value.Brief,
			//"desc":        value.Description,
			"is_ship_free": value.IsShipFree,
			"images":       value.Images,
			"desc_images":  value.DescImages,
			"front_image":  value.FrontImage,
			"pro_price":    value.ProPrice,
			"normal_price": value.NormalPrice,
			"category": map[string]interface{}{
				"id":   value.ItemID,
				"name": value.Item.Name,
			},
			"brand": map[string]interface{}{
				"id":   value.Brand.Id,
				"name": value.Brand.Name,
				"logo": value.Brand.Logo,
			},
			"is_new":     value.IsNew,
			"is_on_sale": value.OnSale,
		})
	}
	reMap["data"] = productList

	ctx.JSON(http.StatusOK, reMap)
}

func CreateProduct(ctx *gin.Context) {
	product := form.ProductInfoForm{}
	if err := ctx.ShouldBindJSON(&product); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	rsp, err := global.ProductClient.CreateProduct(context.Background(), &proto.CreateProductInfo{
		Name:        product.Name,
		ArticleNum:  product.ArticleNum,
		NormalPrice: product.NormalPrice,
		ProPrice:    product.ProPrice,
		Brief:       product.Brief,

		IsShipFree: *product.IsShipFree,
		Images:     product.Images,
		DescImages: product.DescImages,
		FrontImage: product.FrontImage,
		ItemID:     product.ItemId,
		BrandID:    product.Brand,
		Stocks:     product.Stock,
	})
	if err != nil {

		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	zap.S().Info("前端获取的库存为：", product.Stock)

	ctx.JSON(http.StatusOK, rsp)
}

func GetProduct(ctx *gin.Context) {
	zap.S().Info("调用[product.GetProduct]")
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	r, err := global.ProductClient.GetProductDetailByID(context.Background(), &proto.ProductID{Id: int32(id)})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	result := map[string]interface{}{
		"id":           r.Id,
		"name":         r.Name,
		"brief":        r.Brief,
		"description":  r.Description,
		"is_ship_free": r.IsShipFree,
		"images":       r.Images,
		"desc_images":  r.DescImages,
		"front_image":  r.FrontImage,
		"pro_price":    r.ProPrice,
		"normal_price": r.NormalPrice,
		"item": map[string]interface{}{
			"id":   r.Item.Id,
			"name": r.Item.Name,
		},
		"brand": map[string]interface{}{
			"id":   r.Brand.Id,
			"name": r.Brand.Name,
			"logo": r.Brand.Logo,
		},

		"is_new":      r.IsNew,
		"on_sale":     r.OnSale,
		"article_num": r.ArticleNum,
	}

	ctx.JSON(http.StatusOK, result)
}

func DeleteProduct(ctx *gin.Context) {}
func GetStock(ctx *gin.Context)      {}

func Update(ctx *gin.Context) {
	zap.S().Info("调用[product.Update]")

	product := form.ProductInfoForm{}
	if err := ctx.ShouldBindJSON(&product); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	zap.S().Info("操作商品id为:", i)
	zap.S().Infof("product.DescImages:%+v", product.DescImages)

	if _, err = global.ProductClient.UpdateProduct(context.Background(), &proto.CreateProductInfo{
		Id:          int32(i),
		Name:        product.Name,
		ArticleNum:  product.ArticleNum,
		Stocks:      product.Stock,
		NormalPrice: product.NormalPrice,
		ProPrice:    product.ProPrice,
		Brief:       product.Brief,
		IsShipFree:  *product.IsShipFree,
		Images:      product.Images,
		DescImages:  product.DescImages,
		FrontImage:  product.FrontImage,
		ItemID:      product.ItemId,
		BrandID:     product.Brand,
	}); err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "商品信息更新成功",
	})
}
func UpdateStatus(ctx *gin.Context) {}
