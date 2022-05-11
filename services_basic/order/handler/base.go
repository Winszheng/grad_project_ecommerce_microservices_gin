package handler

import (
	"fmt"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

func GenerateNo(userID int32) string {
	now := time.Now()
	rand.Seed(time.Now().UnixNano())
	no := fmt.Sprintf("%d%d%d%d%d%d%d%d",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Nanosecond(),
		userID, rand.Intn(90)+10)
	return no
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
