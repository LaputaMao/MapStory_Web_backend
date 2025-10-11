// Package handler my-web-app/internal/handler/user_handler.go
package handler

import (
	"StoryMap/internal/model"
	"StoryMap/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// RegisterRequest 定义注册请求的结构
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Role     string `json:"role" binding:"required"`
}

// LoginRequest 定义登录请求的结构
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

// Register 处理注册请求
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}

	createdUser, err := h.service.Register(user)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": createdUser})
}

// Login 处理登录请求
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}

// ListUsers 处理获取用户列表的请求
func (h *UserHandler) ListUsers(c *gin.Context) {
	// 解析分页参数，提供默认值
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	// 获取搜索参数
	username := c.Query("username")

	users, total, err := h.service.ListUsers(page, pageSize, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user list"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": users,
		"meta": gin.H{
			"total":      total,
			"page":       page,
			"pageSize":   pageSize,
			"totalPages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}
