package middleware

import (
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/order/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func IsAdminAuth(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	currentUser := claims.(*model.CustomClaims)

	if currentUser.AuthorityId != 2 { // 403
		ctx.JSON(http.StatusForbidden, gin.H{"msg": "用户无权限"})
		ctx.Abort()
		return
	} else {
		zap.S().Info("admin有权限查看用户列表")
	}

	ctx.Next()
}
