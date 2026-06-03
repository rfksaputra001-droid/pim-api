package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"pim-api-go/internal/config"
	"pim-api-go/internal/handler"
	"pim-api-go/internal/middleware"
	"pim-api-go/internal/repository"
	"pim-api-go/internal/service"
)

func main() {
	config.Load()
	config.ConnectDB()

	db := config.DB

	adminRepo := repository.NewAdminRepo(db)
	contribRepo := repository.NewContributionRepo(db)
	expenseRepo := repository.NewExpenseRepo(db)

	authSvc := service.NewAuthService(adminRepo)
	contribSvc := service.NewContributionService(contribRepo, expenseRepo)
	expenseSvc := service.NewExpenseService(expenseRepo)
	userSvc := service.NewUserService(adminRepo)

	authH := handler.NewAuthHandler(authSvc)
	contribH := handler.NewContributionHandler(contribSvc)
	expenseH := handler.NewExpenseHandler(expenseSvc)
	userH := handler.NewUserHandler(userSvc)

	r := gin.Default()
	r.SetTrustedProxies(nil)

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", frontendURL)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.MaxMultipartMemory = 10 << 20

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true, "ts": "now"})
	})

	pub := r.Group("/api/public")
	{
		pub.GET("/dashboard", contribH.Dashboard)
		pub.GET("/contributions", contribH.ListPublic)
		pub.POST("/contributions", contribH.Submit)
		pub.GET("/expenses", expenseH.ListPublic)
		pub.GET("/expenses/:id/receipt", expenseH.GetReceipt)
	}

	admin := r.Group("/api/admin")
	{
		admin.POST("/auth/login", authH.Login)
		admin.POST("/auth/logout", authH.Logout)
		admin.GET("/auth/me", middleware.RequireAuth(), authH.Me)

		auth := admin.Group("", middleware.RequireAuth())
		// static route sebelum parametrik untuk hindari routing conflict
		auth.GET("/contributions/export", contribH.ExportCSV)
		auth.GET("/contributions", contribH.ListAdmin)
		auth.PATCH("/contributions/:id/verify", contribH.Verify)
		auth.PATCH("/contributions/:id/reject", contribH.Reject)
		auth.POST("/expenses", expenseH.Create)
		auth.GET("/expenses", expenseH.ListAdmin)

		superAdmin := auth.Group("", middleware.RequireSuperAdmin())
		superAdmin.GET("/users", userH.List)
		superAdmin.POST("/users", userH.Create)
		superAdmin.PATCH("/users/:id/role", userH.UpdateRole)
		superAdmin.PATCH("/users/:id/password", userH.UpdatePassword)
		superAdmin.DELETE("/users/:id", userH.Delete)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3002"
	}
	log.Printf("Server jalan di port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
