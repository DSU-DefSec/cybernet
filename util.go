package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func errorOut(c *gin.Context, err error) {
	log.Println("error:", err)
	c.JSON(400, gin.H{"error": "Invalid request."})
	c.Abort()
}

func errorOutGraceful(c *gin.Context, err error) {
	log.Println("error:", err)
	c.Redirect(http.StatusSeeOther, "/")
	c.Abort()
}

func errorOutAnnoying(c *gin.Context, err error) {
	log.Println("error:", err)
	c.Redirect(http.StatusSeeOther, "/forbidden")
	c.Abort()
}

func parseTime(timeStr string) time.Time {
	timeStr += " " + locString
	parsedTime, err := time.Parse("01/02/06 3:04 MST", timeStr)
	if err != nil {
		log.Println("time parsing failed,", timeStr, "did not parse correctly:", err.Error())
	}
	return parsedTime
}

func formatTime(dur time.Duration) string {
	durSeconds := dur.Microseconds() / 1000000
	seconds := durSeconds % 60
	durSeconds -= seconds
	minutes := (durSeconds % (60 * 60)) / 60
	durSeconds -= minutes * 60
	hours := durSeconds / (60 * 60)
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
