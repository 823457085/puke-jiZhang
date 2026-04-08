package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"puke-jiZhang/internal/service"
)

type RoomHandler struct {
	roomService *service.RoomService
}

func NewRoomHandler() *RoomHandler {
	return &RoomHandler{roomService: service.NewRoomService()}
}

func getUserID(c *gin.Context) (int64, bool) {
	uid, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return 0, false
	}
	return uid.(int64), true
}

// POST /api/v1/rooms
func (h *RoomHandler) CreateRoom(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var req struct {
		Name     string `json:"name" binding:"required"`
		GameType string `json:"game_type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.GameType == "" {
		req.GameType = "generic"
	}

	room, err := h.roomService.CreateRoom(req.Name, req.GameType, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"room": room})
}

// GET /api/v1/rooms
func (h *RoomHandler) GetMyRooms(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	rooms, err := h.roomService.GetUserRooms(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rooms": rooms})
}

// GET /api/v1/rooms/:id
func (h *RoomHandler) GetRoom(c *gin.Context) {
	_, ok := getUserID(c)
	if !ok {
		return
	}

	roomID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}

	room, err := h.roomService.GetRoom(roomID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}

	members, err := h.roomService.GetRoomMembers(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	bills, err := h.roomService.GetRoomBills(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	balances, err := h.roomService.CalculateBalances(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 填充用户名到balance
	userMap := make(map[int64]string)
	for _, b := range balances {
		if _, ok := userMap[b.UserID]; !ok {
			var nickname string
			service.NewUserService().GetUser(b.UserID)
			userMap[b.UserID] = nickname
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"room":    room,
		"members": members,
		"bills":   bills,
		"balances": balances,
	})
}

// POST /api/v1/rooms/join
func (h *RoomHandler) JoinRoom(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room, err := h.roomService.JoinRoom(req.Code, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"room": room})
}

// POST /api/v1/rooms/:id/close
func (h *RoomHandler) CloseRoom(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	roomID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}

	err = h.roomService.CloseRoom(roomID, userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "room closed"})
}

// GET /api/v1/rooms/:id/balance
func (h *RoomHandler) GetBalance(c *gin.Context) {
	_, ok := getUserID(c)
	if !ok {
		return
	}

	roomID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}

	balances, err := h.roomService.CalculateBalances(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balances": balances})
}

// POST /api/v1/rooms/:id/settle
func (h *RoomHandler) Settle(c *gin.Context) {
	_, ok := getUserID(c)
	if !ok {
		return
	}

	roomID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}

	balances, err := h.roomService.CalculateBalances(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	settlements := h.roomService.CalculateSettlement(balances)
	settlements = h.roomService.EnrichSettlements(settlements)

	c.JSON(http.StatusOK, gin.H{"settlements": settlements})
}
