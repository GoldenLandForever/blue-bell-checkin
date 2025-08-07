package mysql

import (
	"bluebell_backend/models"

	"go.uber.org/zap"
)

func CreateComment(comment *models.Comment) (err error) {
	sqlStr := `insert into comment(
	comment_id, content, post_id, author_id, parent_id)
	values(?,?,?,?,?)`
	_, err = Db.Exec(sqlStr, comment.CommentID, comment.Content, comment.PostID,
		comment.AuthorID, comment.ParentID)
	if err != nil {
		zap.L().Error("insert comment failed", zap.Error(err))
		err = ErrorInsertFailed
		return
	}
	return
}

func GetCommentListByIDs(CommentList models.CommentList) (comments []*models.Comment, err error) {
	// 使用字符串拼接处理排序关键字（注意防范SQL注入！）
	sqlStr := `SELECT comment_id, content, post_id, author_id, parent_id, create_time
               FROM comment
               WHERE post_id = ?
               ORDER BY create_time ` + CommentList.Order + ` limit ?,?`

	comments = make([]*models.Comment, 0, CommentList.Size)
	// 注意：现在只有一个占位符 (?)，所以只传入 PostID
	err = Db.Select(&comments, sqlStr, CommentList.PostID, (CommentList.Page-1)*CommentList.Size, CommentList.Size)
	if err != nil {
		zap.L().Error("查询帖子下评论失败", zap.Error(err))
		return nil, err
	}
	return comments, nil
}
