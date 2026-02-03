package service

import (
	"Go_skeleton/internal/biz"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UserService 持有业务逻辑 Usecase
type UserService struct {
	uc *biz.UserUsecase
}

func NewUserService(uc *biz.UserUsecase) *UserService {
	return &UserService{uc: uc}
}

// RegisterRoutes 注册路由
func (s *UserService) RegisterRoutes(r *gin.Engine) {
	r.POST("/users", func(c *gin.Context) {
		err := s.uc.CreateUser(c.Request.Context(), &biz.User{Username: "test_user"})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "用户创建成功"})
	})
}
