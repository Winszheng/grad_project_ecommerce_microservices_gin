package model

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	ID          uint
	Nickname    string
	AuthorityId uint // 角色、权限
	jwt.StandardClaims
}
