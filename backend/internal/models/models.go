package models

import "time"

// User 用户
type User struct {
	ID        int64     `json:"id"`
	OpenID    string    `json:"openid"`   // 微信openid
	Nickname  string    `json:"nickname"`  // 昵称
	AvatarURL string    `json:"avatar_url"`// 头像
	CreatedAt time.Time `json:"created_at"`
}

// Room 房间
type Room struct {
	ID        int64     `json:"id"`
	RoomCode  string    `json:"room_code"` // 6位房间码
	Name      string    `json:"name"`       // 房间名称
	GameType  string    `json:"game_type"`  // poker/mahjong/generic
	CreatorID int64     `json:"creator_id"`
	Status    string    `json:"status"`     // active/closed
	CreatedAt time.Time `json:"created_at"`
	ClosedAt  *time.Time `json:"closed_at,omitempty"`
}

// RoomMember 房间成员
type RoomMember struct {
	ID       int64     `json:"id"`
	RoomID   int64     `json:"room_id"`
	UserID   int64     `json:"user_id"`
	User     *User     `json:"user,omitempty"`
	JoinedAt time.Time `json:"joined_at"`
}

// Bill 账单
type Bill struct {
	ID         int64     `json:"id"`
	RoomID     int64     `json:"room_id"`
	PayerID    int64     `json:"payer_id"`   // 付款人
	ReceiverID int64     `json:"receiver_id"`// 收款人
	Payer      *User     `json:"payer,omitempty"`
	Receiver   *User     `json:"receiver,omitempty"`
	Amount     int64     `json:"amount"`     // 金额（分，整数）
	Note       string    `json:"note"`       // 备注
	CreatedAt  time.Time `json:"created_at"`
}

// Settlement 结算记录
type Settlement struct {
	ID        int64     `json:"id"`
	RoomID    int64     `json:"room_id"`
	SettledAt time.Time `json:"settled_at"`
}

// Balance 成员余额（欠款计算结果）
type Balance struct {
	UserID  int64  `json:"user_id"`
	Balance int64  `json:"balance"` // 正=应收，负=应付
	User    *User  `json:"user,omitempty"`
}

// SettlementItem 还款方案项
type SettlementItem struct {
	FromUserID   int64  `json:"from_user_id"`
	ToUserID     int64  `json:"to_user_id"`
	Amount       int64  `json:"amount"`
	FromUserName string `json:"from_user_name,omitempty"`
	ToUserName   string `json:"to_user_name,omitempty"`
}
