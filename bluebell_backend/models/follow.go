package models

type Follow struct {
	FollowedID uint64 `db:"followed_id" json:"followed_id,string"` // 要关注的用户ID
}
