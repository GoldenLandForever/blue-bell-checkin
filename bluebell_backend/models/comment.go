package models

import "time"

type Comment struct {
	PostID     uint64    `db:"post_id" json:"post_id,string"`
	ParentID   uint64    `db:"parent_id" json:"parent_id"`
	CommentID  uint64    `db:"comment_id" json:"comment_id"`
	AuthorID   uint64    `db:"author_id" json:"author_id"`
	Content    string    `db:"content" json:"content"`
	CreateTime time.Time `db:"create_time" json:"create_time"`
}

type CommentList struct {
	PostID      uint64 `db:"post_id" json:"post_id,string" form:"post_id"`
	CommunityID uint64 `json:"community_id" form:"community_id"`   // 可以为空
	Page        int64  `json:"page" form:"page"`                   // 页码
	Size        int64  `json:"size" form:"size"`                   // 每页数量
	Order       string `json:"order" form:"order" example:"score"` // 排序依据
}
