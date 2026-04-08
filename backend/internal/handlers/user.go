package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"puke-jiZhang/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{userService: service.NewUserService()}
}

// POST /api/v1/auth/login
// 微信授权登录（实际需调用微信API，用code换openid）
// 这里简化处理：直接传openid注册/登录
func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		OpenID    string `json:"openid" binding:"required"`
		Nickname  string `json:"nickname"`
		AvatarURL string `json:"avatar_url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.LoginOrRegister(req.OpenID, req.Nickname, req.AvatarURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// GET /api/v1/user/me
func (h *UserHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.userService.GetUser(userID.(int64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
