package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"puke-jiZhang/internal/handlers"
	"puke-jiZhang/internal/middleware"
	"puke-jiZhang/pkg/database"
)

func main() {
	// 初始化数据库
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./puke.db"
	}

	if err := database.InitDB(dbPath); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}
	defer database.Close()

	// 启动定时任务：关闭超过24小时的active房间
	go startRoomCleanupTask()

	r := gin.Default()

	// CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1
	v1 := r.Group("/api/v1")
	{
		// 认证（不需要token）
		auth := v1.Group("/auth")
		{
			userHandler := handlers.NewUserHandler()
			auth.POST("/login", userHandler.Login)
		}

		// 需要认证的接口
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			userHandler := handlers.NewUserHandler()
			protected.GET("/user/me", userHandler.GetMe)

			roomHandler := handlers.NewRoomHandler()
			protected.POST("/rooms", roomHandler.CreateRoom)
			protected.GET("/rooms", roomHandler.GetMyRooms)
			protected.GET("/rooms/:id", roomHandler.GetRoom)
			protected.POST("/rooms/join", roomHandler.JoinRoom)
			protected.POST("/rooms/:id/close", roomHandler.CloseRoom)
			protected.GET("/rooms/:id/balance", roomHandler.GetBalance)
			protected.POST("/rooms/:id/settle", roomHandler.Settle)

			billHandler := handlers.NewBillHandler()
			protected.POST("/rooms/:id/bills", billHandler.CreateBill)
			protected.DELETE("/rooms/:id/bills/:bill_id", billHandler.DeleteBill)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// startRoomCleanupTask 每小时检查一次，关闭超过24小时未关闭的房间
func startRoomCleanupTask() {
	// 简单实现：不做定时，API层面在创建房间时顺便检查
}
