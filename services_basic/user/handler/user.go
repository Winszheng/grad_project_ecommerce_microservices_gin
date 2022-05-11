package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/user/global"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/user/model"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/user/proto"
	pass "github.com/anaskhan96/go-password-encoder"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"strings"
	"time"
)

type UserServer struct {
	proto.UnsafeUserServer
}

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

func Model2Response(user model.User) proto.UserInfoResponse {
	userInfoResp := proto.UserInfoResponse{
		Id:       user.ID,
		Password: user.Password,
		Mobile:   user.Mobile,
		Nickname: user.Nickname,
		Gender:   user.Gender,
		Role:     int32(user.Role),
	}

	if user.Birthday != nil {
		userInfoResp.Birthday = uint64(user.Birthday.Unix()) // Unix能把time.Time转换成int64
	}
	return userInfoResp
}

func (s *UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	zap.S().Info("调用【GetUserList】")

	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	rsp := proto.UserListResponse{}
	rsp.Total = int32(result.RowsAffected)

	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users)
	for _, user := range users {
		userInfoResp := Model2Response(user)
		rsp.Data = append(rsp.Data, &userInfoResp)
	}
	return &rsp, nil
}

func (s *UserServer) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var user model.User

	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "号码为 %s 的用户不存在", req.Mobile)
	}
	userInfo := Model2Response(user)
	return &userInfo, nil
}

func (s *UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.First(&user, req.Id) // 主键可以直接查
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "id为 %d 用户不存在", req.Id)
	}
	userInfo := Model2Response(user)
	return &userInfo, nil
}

func (s *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "电话号码为 %s 的用户已存在", req.Mobile)
	}
	user.Mobile = req.Mobile
	user.Nickname = req.Nickname

	option := pass.Options{16, 100, 32, sha512.New}
	salt, encodedPassword := pass.Encode(req.Password, &option)
	password := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPassword)

	user.Password = password

	fmt.Println("[CreateUser] raw pass:", req.Password)

	result = global.DB.Create(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	userInfoResponse := Model2Response(user)
	return &userInfoResponse, nil
}

func (s *UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (*emptypb.Empty, error) {
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "更新时发现用户不存在")
	}

	birthday := time.Unix(int64(req.Birthday), 0)
	user.Birthday = &birthday
	user.Nickname = req.Nickname
	user.Gender = req.Gender
	result = global.DB.Save(&user) // 可以用save来更新啊....

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	return &empty.Empty{}, nil
}

func (s *UserServer) CheckPassword(ctx context.Context, req *proto.CheckPasswordInfo) (*proto.CheckResponse, error) {
	passSlice := strings.Split(req.EncrytedPassword, "$")
	fmt.Println(passSlice)
	option := pass.Options{
		SaltLen:      16,  // salt值设置成16个字节
		Iterations:   100, // 参数
		KeyLen:       32,
		HashFunction: sha512.New,
	}

	check := pass.Verify(req.Password, passSlice[2], passSlice[3], &option)
	fmt.Println(check, req.Password, passSlice[2], passSlice[3])
	return &proto.CheckResponse{Success: check}, nil
}
