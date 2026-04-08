package service

import (
	"database/sql"
	"errors"

	"puke-jiZhang/internal/models"
	"puke-jiZhang/pkg/database"
)

type BillService struct{}

func NewBillService() *BillService {
	return &BillService{}
}

// 记一笔账
func (s *BillService) CreateBill(roomID int64, payerID int64, receiverID int64, amount int64, note string) (*models.Bill, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	// 验证两人都是房间成员
	memberRow := database.DB.QueryRow(
		"SELECT id FROM room_members WHERE room_id = ? AND user_id = ?",
		roomID, payerID,
	)
	var memberID int64
	if err := memberRow.Scan(&memberID); err != nil {
		return nil, errors.New("payer is not a room member")
	}

	memberRow = database.DB.QueryRow(
		"SELECT id FROM room_members WHERE room_id = ? AND user_id = ?",
		roomID, receiverID,
	)
	if err := memberRow.Scan(&memberID); err != nil {
		return nil, errors.New("receiver is not a room member")
	}

	result, err := database.DB.Exec(
		"INSERT INTO bills (room_id, payer_id, receiver_id, amount, note) VALUES (?, ?, ?, ?, ?)",
		roomID, payerID, receiverID, amount, note,
	)
	if err != nil {
		return nil, err
	}

	billID, _ := result.LastInsertId()
	return s.GetBill(billID)
}

func (s *BillService) GetBill(billID int64) (*models.Bill, error) {
	row := database.DB.QueryRow(`
		SELECT b.id, b.room_id, b.payer_id, b.receiver_id, b.amount, b.note, b.created_at,
			   p.id, p.openid, p.nickname, p.avatar_url, p.created_at,
			   r.id, r.openid, r.nickname, r.avatar_url, r.created_at
		FROM bills b
		INNER JOIN users p ON b.payer_id = p.id
		INNER JOIN users r ON b.receiver_id = r.id
		WHERE b.id = ?
	`, billID)

	bill := &models.Bill{}
	var p, r models.User
	err := row.Scan(&bill.ID, &bill.RoomID, &bill.PayerID, &bill.ReceiverID, &bill.Amount, &bill.Note, &bill.CreatedAt,
		&p.ID, &p.OpenID, &p.Nickname, &p.AvatarURL, &p.CreatedAt,
		&r.ID, &r.OpenID, &r.Nickname, &r.AvatarURL, &r.CreatedAt)
	if err != nil {
		return nil, err
	}
	bill.Payer = &p
	bill.Receiver = &r
	return bill, nil
}

func (s *BillService) DeleteBill(billID int64, userID int64) error {
	// 只能删除自己的账单
	row := database.DB.QueryRow("SELECT payer_id FROM bills WHERE id = ?", billID)
	var payerID int64
	if err := row.Scan(&payerID); err != nil {
		return errors.New("bill not found")
	}
	if payerID != userID {
		return errors.New("only payer can delete the bill")
	}

	_, err := database.DB.Exec("DELETE FROM bills WHERE id = ?", billID)
	return err
}
