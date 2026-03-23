package data

import (
	"Go_skeleton/internal/biz" // 引用 biz 定义的实体
	"Go_skeleton/pkg/logger"
	"context"
)

type userRepo struct {
	// 这里可以放 db *gorm.DB
}

// NewUserRepo 构造函数
func NewUserRepo() biz.UserRepo {
	return &userRepo{}
}

func (r *userRepo) Save(ctx context.Context, user *biz.User) error {
	logger.Log.Info("Data层：正在模拟保存用户到数据库...")
	return nil
}

func (r *userRepo) FindByID(ctx context.Context, id int64) (*biz.User, error) {
	return &biz.User{ID: id, Username: "模拟用户"}, nil
}
