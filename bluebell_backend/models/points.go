package models

import (
	"encoding/json"
	"errors"
	"time"
)

// UserPoints 用户积分模型
type UserPoints struct {
	ID          int64      `json:"id,string" db:"id"`
	UserID      int64      `json:"user_id,string" db:"user_id"`
	Points      int        `json:"points" db:"points"`
	PointsTotal int        `json:"points_total" db:"points_total"`
	CreateTime  time.Time  `json:"create_time" db:"create_time"`
	UpdateTime  time.Time  `json:"update_time" db:"update_time"`
	DeleteTime  *time.Time `json:"delete_time" db:"delete_time"`
}

// UserPointsTransaction 积分交易记录模型
type UserPointsTransaction struct {
	ID              int64     `json:"id,string" db:"id"`
	TransactionID   int64     `json:"transaction_id,string" db:"transaction_id"`
	UserID          int64     `json:"user_id,string" db:"user_id"`
	PointsChange    int       `json:"points_change" db:"points_change"`
	CurrentBalance  int       `json:"current_balance" db:"current_balance"`
	TransactionType int       `json:"transaction_type" db:"transaction_type"`
	Description     string    `json:"description" db:"description"`
	ExtJson         string    `json:"ext_json" db:"ext_json"`
	CreateTime      time.Time `json:"create_time" db:"create_time"`
	UpdateTime      time.Time `json:"update_time" db:"update_time"`
	DeleteTime      time.Time `json:"delete_time" db:"delete_time"`
}

// PointsChangeForm 积分变动请求参数
type PointsChangeForm struct {
	UserID      int64  `json:"user_id,string" binding:"required"`
	Points      int    `json:"points" binding:"required"`
	Description string `json:"description" binding:"required"`
	ExtType     string `json:"ext_type" binding:"required"`
	ExtID       string `json:"ext_id" binding:"required"`
}

// UnmarshalJSON 为PointsChangeForm类型实现自定义的UnmarshalJSON方法
func (p *PointsChangeForm) UnmarshalJSON(data []byte) (err error) {
	required := struct {
		UserID      int64  `json:"user_id"`
		Points      int    `json:"points"`
		Description string `json:"description"`
		ExtType     string `json:"ext_type"`
		ExtID       string `json:"ext_id"`
	}{}
	err = json.Unmarshal(data, &required)
	if err != nil {
		return
	} else if required.UserID == 0 {
		err = errors.New("缺少必填字段user_id")
	} else if required.Points == 0 {
		err = errors.New("缺少必填字段points")
	} else if len(required.Description) == 0 {
		err = errors.New("缺少必填字段description")
	} else if len(required.ExtType) == 0 {
		err = errors.New("缺少必填字段ext_type")
	} else if len(required.ExtID) == 0 {
		err = errors.New("缺少必填字段ext_id")
	} else {
		p.UserID = required.UserID
		p.Points = required.Points
		p.Description = required.Description
		p.ExtType = required.ExtType
		p.ExtID = required.ExtID
	}
	return
}