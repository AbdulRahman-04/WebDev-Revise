package main

import (
	// "context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/config"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/controllers/private"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/controllers/public"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/middleware"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/routes"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	ginlimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
	memory "github.com/ulule/limiter/v3/drivers/store/memory"
)

func main() {
	// 🚀 Use all available CPU cores for max concurrency
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Printf("🚀 Using %d CPU cores for concurrency\n", runtime.NumCPU())

	// ----------------- DB + Redis -----------------
	utils.DBConnect()
	utils.ConnectRedis()

	// ----------------- Gin Engine -----------------
	router := gin.New()
	router.Use(gin.Recovery()) // built-in panic recovery
	router.Use(middleware.CustomLogger())
	router.Use(SecureHeaders())

	// ----------------- Global Rate Limiter -----------------
	store := memory.NewStore()
	rate, _ := limiter.NewRateFromFormatted("100-S") // 100 req/sec global
	instance := limiter.New(store, rate)
	router.Use(ginlimiter.NewMiddleware(instance))

	// ----------------- CORS -----------------
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 🔥 you can restrict later
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// ----------------- Init Mongo Collections -----------------
	public.UserCollect()
	public.AdminCollect()
	private.UserAccessCollect()
	private.EventsCollect()
	private.FunctionCollect()
	private.AdminAccessCollect()

	// ----------------- Test Route -----------------
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "🚀 Event Management AI Backend is Live!"})
	})

	// ----------------- Register Routes -----------------
	routes.PublicRoutes(router)
	routes.PrivateRoutes(router)

	// ----------------- Graceful Shutdown -----------------
	srvPort := fmt.Sprintf(":%d", config.AppConfig.Port)
	serverExit := make(chan os.Signal, 1)
	signal.Notify(serverExit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("🌍 Server running on http://localhost%s", srvPort)
		if err := router.Run(srvPort); err != nil {
			log.Fatalf("❌ Server crashed: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-serverExit
	log.Println("🧹 Shutting down gracefully...")

	// (Optional) Add cleanup logic here
	if utils.RedisClient != nil {
		_ = utils.RedisClient.Close()
	}
	log.Println("✅ Redis closed successfully")

	log.Println("✅ Server stopped cleanly")
	os.Exit(0)

}

// ----------------- Secure Headers Middleware -----------------
func SecureHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		c.Next()
	}
}
