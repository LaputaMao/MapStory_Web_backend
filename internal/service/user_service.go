// Package service my-web-app/internal/service/user_service.go
package service

import (
	"StoryMap/internal/model"
	"StoryMap/internal/store"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

// JWT 签名密钥，在生产环境中应该从配置文件读取
var jwtSecret = []byte("your-very-secret-key")

// JWTClaims 定义了 JWT 中存储的数据
type JWTClaims struct {
	UserID   uint   `json:"userId"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type UserService struct {
	store *store.UserStore
}

func NewUserService(s *store.UserStore) *UserService {
	return &UserService{store: s}
}

// Register 处理用户注册逻辑
func (s *UserService) Register(user *model.User) (*model.User, error) {
	// 检查用户或邮箱是否已存在
	_, err := s.store.FindByEmail(user.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err // 其他数据库错误
	}

	// 创建用户（密码哈希会在 BeforeCreate Hook 中自动完成）
	if err := s.store.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

// Login 处理用户登录逻辑
func (s *UserService) Login(email, password string) (string, error) {
	// 1. 查找用户
	user, err := s.store.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("invalid email or password")
		}
		return "", err
	}

	// 2. 验证密码
	if !user.CheckPasswordHash(password) {
		return "", errors.New("invalid email or password")
	}

	// 3. 更新最后登录时间
	now := time.Now()
	user.LastLoginTime = &now
	if err := s.store.Update(user); err != nil {
		// 即使更新失败，登录也算成功，可以只记录日志
	}

	// 4. 生成 JWT
	token, err := s.generateJWT(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

// ListUsers 处理获取用户列表的逻辑
func (s *UserService) ListUsers(page, pageSize int, username string) ([]model.User, int64, error) {
	return s.store.List(page, pageSize, username)
}

// generateJWT 生成一个 JWT
func (s *UserService) generateJWT(user *model.User) (string, error) {
	claims := JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token 有效期 24 小时
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
