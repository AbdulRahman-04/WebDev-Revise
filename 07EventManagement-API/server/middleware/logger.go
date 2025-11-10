package middleware

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// Define color codes for console log readability
var (
	green  = "\033[32m"
	yellow = "\033[33m"
	red    = "\033[31m"
	reset  = "\033[0m"
)

// statusColor returns color based on HTTP status code
func statusColor(status int) string {
	switch {
	case status >= 200 && status < 300:
		return green // Success responses
	case status >= 300 && status < 400:
		return yellow // Redirects
	default:
		return red // Errors
	}
}

// CustomLogger is a Gin middleware for logging each HTTP request
func CustomLogger() gin.HandlerFunc {
	// 🧩 Ensure logs directory exists before creating log file
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.Mkdir("logs", 0755)
		if err != nil {
			log.Fatalf("❌ Could not create logs directory: %v", err)
		}
	}

	// 🧩 Create or open a log file for persistent logging
	logFile, err := os.OpenFile("logs/server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("❌ Could not open log file: %v", err)
	}

	// 🧩 Create a multi-writer: output logs to both console and file
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	logger := log.New(multiWriter, "[IVENTS LOG] ", log.LstdFlags)

	return func(c *gin.Context) {
		// 🧩 Record request start time
		start := time.Now()

		// 🧩 Generate a unique trace ID for each request
		traceID := fmt.Sprintf("%d", time.Now().UnixNano())
		c.Set("traceID", traceID)

		// 🧩 Process the request
		c.Next()

		// 🧩 Calculate request duration
		duration := time.Since(start).Round(time.Millisecond)

		// 🧩 Extract response and request info
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()
		color := statusColor(status)

		// 🧩 Format and log request details (both console + file)
		logger.Printf("%s[%s]%s %s | %s%3d%s | %-7s | %-15s | %s | Took: %v",
			green, start.Format("2006-01-02 15:04:05"), reset,
			traceID, color, status, reset,
			method, clientIP, path, duration,
		)
	}
}
