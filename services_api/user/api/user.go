package api

import (
	"context"
	"fmt"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/user/forms"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/user/global"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/user/global/response"
	middlewares "github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/user/middleware"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/user/model"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_api/user/proto"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{"msg": e.Message()})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "内部错误"})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "用户服务不可用"})
			default:
				msg := fmt.Sprintf("%d %s", e.Code(), e.Message())
				c.JSON(http.StatusInternalServerError, gin.H{"msg": msg})
			}
		}
		return
	}
}

func GetUserList(ctx *gin.Context) {
	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("pnum", "10")
	pSizeInt, _ := strconv.Atoi(pSize)
	rsp, err := global.UserSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询【用户服务】失败")
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	result := make([]interface{}, 0)
	for _, val := range rsp.Data {
		user := response.UserResponse{
			Id:       val.Id,
			Nickname: val.Nickname,
			Birthday: response.JsonTime(time.Unix(int64(val.Birthday), 0)),
			Gender:   val.Gender,
			Mobile:   val.Mobile,
		}

		result = append(result, user)
	}

	ctx.JSON(http.StatusOK, result)
}

func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

func HandleValidatorError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})
	return
}

func PasswordLogin(ctx *gin.Context) {

	passwordLoginForm := forms.PasswordLoginForm{}
	if err := ctx.ShouldBindJSON(&passwordLoginForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	zap.S().Info("CaptchaId:", passwordLoginForm.CaptchaId, "Captcha:", passwordLoginForm.CaptchaAnswer)
	if !store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.CaptchaAnswer, true) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"captcha": "验证码错误",
		})
		return
	}

	if rsp, err := global.UserSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: passwordLoginForm.Mobile,
	}); err != nil {

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusBadRequest, gin.H{"msg": "用户不存在"})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "登录失败"}) // 其余错误都放在default好了
			}
		}
	} else {

		if passRsp, passErr := global.UserSrvClient.CheckPassword(context.Background(), &proto.CheckPasswordInfo{
			Password:         passwordLoginForm.Password, // 提交过来的raw密码
			EncrytedPassword: rsp.Password,               // 数据库里查出来的密码
		}); passErr != nil {

			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "登录错误"})
		} else {
			if passRsp.Success {

				j := middlewares.NewJWT()
				zap.S().Info("role:", rsp.Role)
				claims := model.CustomClaims{
					ID:          uint(rsp.Id),
					Nickname:    rsp.Nickname,
					AuthorityId: uint(rsp.Role),
					StandardClaims: jwt.StandardClaims{
						NotBefore: time.Now().Unix(),               // 签名的生效时间
						ExpiresAt: time.Now().Unix() + 60*60*24*30, // 30天过期，已经是数字就没必要用time.Second了
						Issuer:    "Yuno",
					},
				}
				token, err := j.CreateToken(claims)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "生成token失败"})
					return
				}

				ctx.JSON(http.StatusOK, gin.H{"msg": "登录成功",
					"id":         rsp.Id,
					"nickname":   rsp.Nickname,
					"role":       rsp.Role,
					"token":      token,
					"expired at": (time.Now().Unix() + 60*60*24*30) * 1000,
				})

			} else {
				ctx.JSON(http.StatusNotFound, gin.H{"msg": "登录失败"})
			}

		}

	}
}

func RegisterUser(ctx *gin.Context) {

	registerForm := forms.RegisterUserForm{}
	if err := ctx.ShouldBindJSON(&registerForm); err != nil {
		HandleValidatorError(ctx, err)
	}

	user, err := global.UserSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		Nickname: registerForm.Nickname,
		Password: registerForm.Password,
		Mobile:   registerForm.Mobile,
	})
	if err != nil {
		zap.S().Errorf("[RegisterUser] 查询 【新建用户失败】失败: %s", err.Error())
		HandleValidatorError(ctx, err)
	}

	j := middlewares.NewJWT()
	zap.S().Info("role:", user.Role)
	claims := model.CustomClaims{
		ID:          uint(user.Id),
		Nickname:    user.Nickname,
		AuthorityId: uint(user.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),               // 签名的生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*30, // 30天过期，已经是数字就没必要用time.Second了
			Issuer:    "Yuno",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "生成token失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"id":         user.Id,
		"nick_name":  user.Nickname, // 后端只管返回，前端怎么用是前端的事
		"token":      token,
		"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
	})
}

func GetUserDetail(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	currentUser := claims.(*model.CustomClaims)
	zap.S().Infof("访问用户: %d", currentUser.ID)

	rsp, err := global.UserSrvClient.GetUserById(context.Background(), &proto.IdRequest{
		Id: int32(currentUser.ID),
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"name":     rsp.Nickname,
		"birthday": time.Unix(int64(rsp.Birthday), 0).Format("2006-01-02"),
		"gender":   rsp.Gender,
		"mobile":   rsp.Mobile,
	})
}

func UpdateUser(ctx *gin.Context) {
	updateUserForm := forms.UpdateUserForm{}
	if err := ctx.ShouldBind(&updateUserForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	claims, _ := ctx.Get("claims")
	currentUser := claims.(*model.CustomClaims)
	zap.S().Infof("访问用户: %d", currentUser.ID)

	loc, _ := time.LoadLocation("Local")
	birthDay, _ := time.ParseInLocation("2006-01-02", updateUserForm.Birthday, loc)
	_, err := global.UserSrvClient.UpdateUser(context.Background(), &proto.UpdateUserInfo{
		Id:       int32(currentUser.ID),
		Nickname: updateUserForm.Name,
		Gender:   updateUserForm.Gender,
		Birthday: uint64(birthDay.Unix()),
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}
