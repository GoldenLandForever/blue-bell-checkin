package mysql

import (
	"bluebell_backend/models"
	"database/sql"
	"errors"
	"go.uber.org/zap"
	"time"
)

// 错误变量
var (
	ErrUserPointsNotExist      = errors.New("user points not exist")
	ErrUpdatePointsFailed      = errors.New("update points failed")
	ErrCreateTransactionFailed = errors.New("create transaction failed")
)

// GetUserPointsByUserID 根据用户ID查询用户积分
func GetUserPointsByUserID(userID int64) (*models.UserPoints, error) {
	userPoints := new(models.UserPoints)
	sqlStr := `select id, user_id, points, points_total, create_time, update_time, delete_time from user_points where user_id = ? and delete_time is null`
	err := db.Get(userPoints, sqlStr, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserPointsNotExist
		}
		return nil, err
	}
	return userPoints, nil
}

// UpdateUserPoints 更新用户积分
func UpdateUserPoints(userID int64, points, pointsTotal int) error {
	// 先检查用户积分记录是否存在
	exists := false
	sqlCheck := `select exists(select 1 from user_points where user_id = ? and delete_time is null)`
	if err := db.Get(&exists, sqlCheck, userID); err != nil {
		return err
	}

	now := time.Now()
	if exists {
		// 更新现有记录
		sqlUpdate := `update user_points set points = ?, points_total = ?, update_time = ? where user_id = ? and delete_time is null`
		_, err := db.Exec(sqlUpdate, points, pointsTotal, now, userID)
		if err != nil {
			return ErrUpdatePointsFailed
		}
	} else {
		// 创建新记录
		sqlInsert := `insert into user_points(user_id, points, points_total, create_time, update_time) values(?, ?, ?, ?, ?)`
		_, err := db.Exec(sqlInsert, userID, points, pointsTotal, now, now)
		if err != nil {
			return ErrUpdatePointsFailed
		}
	}
	return nil
}

// CreateUserPointsTransaction 创建积分交易记录
func CreateUserPointsTransaction(transaction *models.UserPointsTransaction) error {
	sqlStr := `insert into user_points_transactions(transaction_id, user_id, points_change, current_balance, transaction_type, description, ext_json, create_time, update_time) values(?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	_, err := db.Exec(sqlStr,
		transaction.TransactionID,
		transaction.UserID,
		transaction.PointsChange,
		transaction.CurrentBalance,
		transaction.TransactionType,
		transaction.Description,
		transaction.ExtJson,
		now,
		now,
	)
	if err != nil {
		zap.L().Error("创建积分交易记录出错", zap.Error(err))
		return ErrCreateTransactionFailed
	}
	return nil
}

//func ErrUserPointsNotExist(user *models.UserPoints) error {
//	sqlStr := `insert into  user_points(user_id, points, points_total, create_time, update_time) values(?, ?, ?, ?, ?)`
//	now := time.Now()
//	_, err := db.Exec(sqlStr,
//		user.UserID,
//		user.Points,
//		user.PointsTotal,
//		now,
//		now)
//	if err != nil {
//		return ErrUserPointsNotExist
//	}
//	return nil
//}
