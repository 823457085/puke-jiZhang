package service

import (
	"database/sql"
	"errors"
	"math/rand"
	"sort"
	"time"

	"puke-jiZhang/internal/models"
	"puke-jiZhang/pkg/database"
)

type RoomService struct{}

func NewRoomService() *RoomService {
	return &RoomService{}
}

func (s *RoomService) GenerateRoomCode() string {
	for i := 0; i < 10; i++ {
		code := ""
		for j := 0; j < 6; j++ {
			code += string(rune('0' + rand.Intn(10)))
		}
		row := database.DB.QueryRow("SELECT id FROM rooms WHERE room_code = ? AND status = 'active'", code)
		var id int64
		if row.Scan(&id) == nil {
			continue
		}
		return code
	}
	return ""
}

func (s *RoomService) CreateRoom(name string, gameType string, creatorID int64) (*models.Room, error) {
	code := s.GenerateRoomCode()
	if code == "" {
		return nil, errors.New("failed to generate unique room code")
	}

	result, err := database.DB.Exec(
		"INSERT INTO rooms (room_code, name, game_type, creator_id) VALUES (?, ?, ?, ?)",
		code, name, gameType, creatorID,
	)
	if err != nil {
		return nil, err
	}

	roomID, _ := result.LastInsertId()

	_, err = database.DB.Exec(
		"INSERT OR IGNORE INTO room_members (room_id, user_id) VALUES (?, ?)",
		roomID, creatorID,
	)
	if err != nil {
		return nil, err
	}

	return s.GetRoom(roomID)
}

func (s *RoomService) GetRoomByCode(code string) (*models.Room, error) {
	row := database.DB.QueryRow(
		"SELECT id, room_code, name, game_type, creator_id, status, created_at, closed_at FROM rooms WHERE room_code = ? AND status = 'active'",
		code,
	)
	room := &models.Room{}
	var closedAt sql.NullTime
	err := row.Scan(&room.ID, &room.RoomCode, &room.Name, &room.GameType, &room.CreatorID, &room.Status, &room.CreatedAt, &closedAt)
	if err != nil {
		return nil, err
	}
	if closedAt.Valid {
		room.ClosedAt = &closedAt.Time
	}
	return room, nil
}

func (s *RoomService) GetRoom(roomID int64) (*models.Room, error) {
	row := database.DB.QueryRow(
		"SELECT id, room_code, name, game_type, creator_id, status, created_at, closed_at FROM rooms WHERE id = ?",
		roomID,
	)
	room := &models.Room{}
	var closedAt sql.NullTime
	err := row.Scan(&room.ID, &room.RoomCode, &room.Name, &room.GameType, &room.CreatorID, &room.Status, &room.CreatedAt, &closedAt)
	if err != nil {
		return nil, err
	}
	if closedAt.Valid {
		room.ClosedAt = &closedAt.Time
	}
	return room, nil
}

func (s *RoomService) JoinRoom(code string, userID int64) (*models.Room, error) {
	room, err := s.GetRoomByCode(code)
	if err != nil {
		return nil, errors.New("room not found")
	}

	_, err = database.DB.Exec(
		"INSERT OR IGNORE INTO room_members (room_id, user_id) VALUES (?, ?)",
		room.ID, userID,
	)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (s *RoomService) GetUserRooms(userID int64) ([]models.Room, error) {
	rows, err := database.DB.Query(`
		SELECT r.id, r.room_code, r.name, r.game_type, r.creator_id, r.status, r.created_at, r.closed_at
		FROM rooms r
		INNER JOIN room_members rm ON r.id = rm.room_id
		WHERE rm.user_id = ?
		ORDER BY r.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []models.Room
	for rows.Next() {
		var room models.Room
		var closedAt sql.NullTime
		rows.Scan(&room.ID, &room.RoomCode, &room.Name, &room.GameType, &room.CreatorID, &room.Status, &room.CreatedAt, &closedAt)
		if closedAt.Valid {
			room.ClosedAt = &closedAt.Time
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func (s *RoomService) CloseRoom(roomID int64, userID int64) error {
	room, err := s.GetRoom(roomID)
	if err != nil {
		return errors.New("room not found")
	}
	if room.CreatorID != userID {
		return errors.New("only creator can close the room")
	}

	now := time.Now()
	_, err = database.DB.Exec(
		"UPDATE rooms SET status = 'closed', closed_at = ? WHERE id = ?",
		now, roomID,
	)
	return err
}

func (s *RoomService) GetRoomMembers(roomID int64) ([]models.RoomMember, error) {
	rows, err := database.DB.Query(`
		SELECT rm.id, rm.room_id, rm.user_id, rm.joined_at,
			   u.id, u.openid, u.nickname, u.avatar_url, u.created_at
		FROM room_members rm
		INNER JOIN users u ON rm.user_id = u.id
		WHERE rm.room_id = ?
		ORDER BY rm.joined_at ASC
	`, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.RoomMember
	for rows.Next() {
		var m models.RoomMember
		var u models.User
		rows.Scan(&m.ID, &m.RoomID, &m.UserID, &m.JoinedAt, &u.ID, &u.OpenID, &u.Nickname, &u.AvatarURL, &u.CreatedAt)
		m.User = &u
		members = append(members, m)
	}
	return members, nil
}

func (s *RoomService) GetRoomBills(roomID int64) ([]models.Bill, error) {
	rows, err := database.DB.Query(`
		SELECT b.id, b.room_id, b.payer_id, b.receiver_id, b.amount, b.note, b.created_at,
			   p.id, p.openid, p.nickname, p.avatar_url, p.created_at,
			   r.id, r.openid, r.nickname, r.avatar_url, r.created_at
		FROM bills b
		INNER JOIN users p ON b.payer_id = p.id
		INNER JOIN users r ON b.receiver_id = r.id
		WHERE b.room_id = ?
		ORDER BY b.created_at DESC
	`, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bills []models.Bill
	for rows.Next() {
		var b models.Bill
		var p, r models.User
		rows.Scan(&b.ID, &b.RoomID, &b.PayerID, &b.ReceiverID, &b.Amount, &b.Note, &b.CreatedAt,
			&p.ID, &p.OpenID, &p.Nickname, &p.AvatarURL, &p.CreatedAt,
			&r.ID, &r.OpenID, &r.Nickname, &r.AvatarURL, &r.CreatedAt)
		b.Payer = &p
		b.Receiver = &r
		bills = append(bills, b)
	}
	return bills, nil
}

func (s *RoomService) CalculateBalances(roomID int64) ([]models.Balance, error) {
	rows, err := database.DB.Query(`
		SELECT user_id FROM room_members WHERE room_id = ?
	`, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	balanceMap := make(map[int64]int64)
	var userIDs []int64
	for rows.Next() {
		var userID int64
		rows.Scan(&userID)
		balanceMap[userID] = 0
		userIDs = append(userIDs, userID)
	}

	billRows, err := database.DB.Query(`
		SELECT payer_id, receiver_id, amount FROM bills WHERE room_id = ?
	`, roomID)
	if err != nil {
		return nil, err
	}
	defer billRows.Close()

	for billRows.Next() {
		var payerID, receiverID, amount int64
		billRows.Scan(&payerID, &receiverID, &amount)
		balanceMap[payerID] -= amount
		balanceMap[receiverID] += amount
	}

	var balances []models.Balance
	for _, userID := range userIDs {
		balance := models.Balance{UserID: userID, Balance: balanceMap[userID]}
		balances = append(balances, balance)
	}
	return balances, nil
}

func (s *RoomService) CalculateSettlement(balances []models.Balance) []models.SettlementItem {
	var creditors []models.Balance
	var debtors []models.Balance

	for _, b := range balances {
		if b.Balance > 0 {
			creditors = append(creditors, b)
		} else if b.Balance < 0 {
			debtors = append(debtors, b)
		}
	}

	sort.Slice(creditors, func(i, j int) bool { return creditors[i].Balance > creditors[j].Balance })
	sort.Slice(debtors, func(i, j int) bool { return debtors[i].Balance < debtors[j].Balance })

	var settlements []models.SettlementItem
	i, j := 0, 0

	for i < len(creditors) && j < len(debtors) {
		credit := creditors[i].Balance
		debt := -debtors[j].Balance

		amount := min(credit, debt)

		if amount > 0 {
			settlements = append(settlements, models.SettlementItem{
				FromUserID: debtors[j].UserID,
				ToUserID:   creditors[i].UserID,
				Amount:     amount,
			})
		}

		credit -= amount
		debt -= amount

		creditors[i].Balance = credit
		debtors[j].Balance = -debt

		if credit == 0 {
			i++
		}
		if debt == 0 {
			j++
		}
	}

	return settlements
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func (s *RoomService) EnrichSettlements(settlements []models.SettlementItem) []models.SettlementItem {
	for i := range settlements {
		var fromNick, toNick string
		database.DB.QueryRow("SELECT nickname FROM users WHERE id = ?", settlements[i].FromUserID).Scan(&fromNick)
		database.DB.QueryRow("SELECT nickname FROM users WHERE id = ?", settlements[i].ToUserID).Scan(&toNick)
		settlements[i].FromUserName = fromNick
		settlements[i].ToUserName = toNick
	}
	return settlements
}
