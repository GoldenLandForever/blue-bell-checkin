package models

type Checkin struct {
	UserID      uint64 `db:"user_id" json:"user_id,string"`
	TimeStamp   int64  `db:"time_stamp" json:"time_stamp,string"`
	CheckinType string `db:"checkin_type" json:"checkin_type"`
}
