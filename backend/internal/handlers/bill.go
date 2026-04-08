package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"puke-jiZhang/internal/service"
)

type BillHandler struct {
	billService *service.BillService
}

func NewBillHandler() *BillHandler {
	return &BillHandler{billService: service.NewBillService()}
}

func getUserIDFromContext(c *gin.Context) (int64, bool) {
	uid, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return 0, false
	}
	return uid.(int64), true
}

// POST /api/v1/rooms/:id/bills
func (h *BillHandler) CreateBill(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	roomID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}

	var req struct {
		PayerID    int64  `json:"payer_id" binding:"required"`
		ReceiverID int64  `json:"receiver_id" binding:"required"`
		Amount     int64  `json:"amount" binding:"required"`
		Note       string `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bill, err := h.billService.CreateBill(roomID, req.PayerID, req.ReceiverID, req.Amount, req.Note)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 记录操作人
	_ = userID

	c.JSON(http.StatusOK, gin.H{"bill": bill})
}

// DELETE /api/v1/rooms/:id/bills/:bill_id
func (h *BillHandler) DeleteBill(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	billID, err := strconv.ParseInt(c.Param("bill_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bill id"})
		return
	}

	err = h.billService.DeleteBill(billID, userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "bill deleted"})
}
