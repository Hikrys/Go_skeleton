package biz

// 这个文件只是示例，你需要根据自己的需要具体实现
import (
	"Go_skeleton/pkg/logger"
	"context"
)

// User 这是业务领域的核心数据结构
type User struct {
	ID       int64
	Username string
	Password string // 实际项目中不要存明文
}

// 定义数据接口
// biz 层只定义“我需要什么”，不关心“怎么存”。
// 具体怎么存（用 MySQL 还是 Redis），是由 data 层实现的（依赖倒置原则）。
type UserRepo interface {
	Save(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id int64) (*User, error)
}

// 定义业务用例
// 这里面写具体的业务逻辑，比如“注册用户”、“修改密码”
type UserUsecase struct {
	repo UserRepo // 它依赖上面的接口
}

// 构造函数 这是 Wire 需要使用的函数
// Wire 会自动发现：要想创建 UserUsecase，必须先给一个 UserRepo
func NewUserUsecase(repo UserRepo) *UserUsecase {
	return &UserUsecase{repo: repo}
}

// 编写一个具体的业务方法（这里是示例）
func (uc *UserUsecase) CreateUser(ctx context.Context, u *User) error {
	// 这里可以写业务逻辑，比如检查密码强度、用户名是否重复
	logger.Log.Info("正在创建用户业务逻辑...")
	return uc.repo.Save(ctx, u)
}
