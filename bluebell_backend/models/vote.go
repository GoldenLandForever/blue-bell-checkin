package models

import "time"

type Vote struct {
	UserID     uint64    `json:"user_id,string" db:"user_id"` // 指定json序列化/反序列化时使用小写user_id
	PostID     uint64    `json:"post_id,string" db:"post_id"` // 指定json序列化/反序列化时使用小写post_id
	VoteType   int8      `json:"vote_type" db:"vote_type"`
	CreateTime time.Time `json:"create_time" db:"create_time"`
	UpdateTime time.Time `json:"update_time" db:"update_time"`
	DeleteTime time.Time `json:"delete_time" db:"delete_time"`
}
