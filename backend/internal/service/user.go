package service

import (
	"database/sql"
	"errors"

	"puke-jiZhang/internal/models"
	"puke-jiZhang/pkg/database"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

// 微信授权登录（模拟：实际需要调用微信API）
func (s *UserService) LoginOrRegister(openid string, nickname string, avatarURL string) (*models.User, error) {
	// 先查是否已存在
	row := database.DB.QueryRow("SELECT id, openid, nickname, avatar_url, created_at FROM users WHERE openid = ?", openid)

	user := &models.User{}
	err := row.Scan(&user.ID, &user.OpenID, &user.Nickname, &user.AvatarURL, &user.CreatedAt)
	if err == nil {
		// 已存在，更新昵称和头像
		if nickname != "" {
			database.DB.Exec("UPDATE users SET nickname = ?, avatar_url = ? WHERE id = ?", nickname, avatarURL, user.ID)
			user.Nickname = nickname
			user.AvatarURL = avatarURL
		}
		return user, nil
	}

	if err != sql.ErrNoRows {
		return nil, err
	}

	// 不存在，创建
	result, err := database.DB.Exec(
		"INSERT INTO users (openid, nickname, avatar_url) VALUES (?, ?, ?)",
		openid, nickname, avatarURL,
	)
	if err != nil {
		return nil, err
	}

	userID, _ := result.LastInsertId()
	return s.GetUser(userID)
}

func (s *UserService) GetUser(userID int64) (*models.User, error) {
	row := database.DB.QueryRow("SELECT id, openid, nickname, avatar_url, created_at FROM users WHERE id = ?", userID)
	user := &models.User{}
	err := row.Scan(&user.ID, &user.OpenID, &user.Nickname, &user.AvatarURL, &user.CreatedAt)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}
