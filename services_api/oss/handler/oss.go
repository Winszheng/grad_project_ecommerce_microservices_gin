package handler

import (
	"fmt"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/oss/global"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"strings"

	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/oss/utils"
	"github.com/gin-gonic/gin"
)

func Token(c *gin.Context) {
	response := utils.Get_policy_token()
	c.Header("Access-Control-Allow-Methods", "POST")
	c.Header("Access-Control-Allow-Origin", "*")
	c.String(200, response)
}

func HandlerRequest(ctx *gin.Context) {
	zap.S().Info("调用[oss.POST]请求")

	bytePublicKey, err := utils.GetPublicKey(ctx)
	if err != nil {
		utils.ResponseFailed(ctx)
		return
	}

	byteAuthorization, err := utils.GetAuthorization(ctx)
	if err != nil {
		utils.ResponseFailed(ctx)
		return
	}

	byteMD5, bodyStr, err := utils.GetMD5FromNewAuthString(ctx)
	if err != nil {
		utils.ResponseFailed(ctx)
		return
	}

	decodeurl, err := url.QueryUnescape(bodyStr)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(decodeurl)
	params := make(map[string]string)
	datas := strings.Split(decodeurl, "&")
	for _, v := range datas {
		sdatas := strings.Split(v, "=")
		fmt.Println(v)
		params[sdatas[0]] = sdatas[1]
	}
	fileName := params["filename"]
	fileUrl := fmt.Sprintf("%s/%s", global.ServerConfig.OssInfo.Host, fileName)

	if utils.VerifySignature(bytePublicKey, byteMD5, byteAuthorization) {

		ctx.JSON(http.StatusOK, gin.H{
			"url": fileUrl, // 上传成功把url返回给前端
		})
	} else {
		utils.ResponseFailed(ctx) // response FAILED : 400
	}
}
