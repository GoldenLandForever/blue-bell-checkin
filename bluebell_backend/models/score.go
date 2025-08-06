package models

// ScoreChange 积分变动记录
// @Description 记录用户积分变动的信息
// @property UserID int64 用户ID
// @property Change int 积分变动值（正数为增加，负数为减少）
// @property Reason string 变动原因（如comment, post, vote等）
// @property CreatedAt int64 创建时间戳

type ScoreChange struct {
	UserID    int64  `json:"user_id"`
	Change    int    `json:"change"`
	Reason    string `json:"reason"`
	CreatedAt int64  `json:"created_at"`
}