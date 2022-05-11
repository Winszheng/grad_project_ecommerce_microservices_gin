package response

import (
	"fmt"
	"time"
)

type UserResponse struct {
	Id       int32    `json:"id"`
	Nickname string   `json:"name"`
	Birthday JsonTime `json:"birthday"` // 生日作为日期，肯定是time比较方便
	Gender   string   `json:"gender"`
	Mobile   string   `json:"mobile"`
}

type JsonTime time.Time

func (j JsonTime) MarshalJSON() ([]byte, error) {
	stmp := time.Time(j).Format("2006-01-02")
	stmp = fmt.Sprintf("\"%s\"", stmp)
	return []byte(stmp), nil
}
