package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	// Hardcoded CT timezone.
	location, _ = time.LoadLocation("America/Rainy_River")
	locString   = "CT"
	db          = &gorm.DB{}
)

func main() {

	// Initialize Gin router
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	initCookies(r)

	// Open database
	var err error
	db, err = gorm.Open(sqlite.Open("cybernet.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database!")
	}

	db.AutoMigrate(&User{}, &Secret{}, &Event{}, &SignUp{})
	// Add... add function
	r.SetFuncMap(template.FuncMap{
		"increment": func(x int) int {
			return x + 1
		},
	})

	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets")

	// 404 handler
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", pageData(c, "Page Not Found", nil))
	})

	// Route definitons
	publicRoutes := r.Group("/")
	{
		publicRoutes.GET("", func(c *gin.Context) {
			now := time.Now()

			comps := []Event{}
			active := []Event{}
			past := []Event{}

			result := db.Order("event_start asc").Find(&comps)
			if result.Error != nil {
				c.HTML(http.StatusOK, "index.html", pageData(c, "Events", gin.H{"error": result.Error}))
				return
			}

			for _, comp := range comps {
				if comp.EventEnd.Before(now) {
					past = append(active, comp)
				} else {
					active = append(active, comp)
				}

			}
			c.HTML(http.StatusOK, "index.html", pageData(c, "Events", gin.H{"active": active, "past": past}))
		})
		publicRoutes.GET("/login", func(c *gin.Context) {
			if getUser(c).IsValid() {
				c.Redirect(http.StatusSeeOther, "/")
			}
			c.HTML(http.StatusOK, "login.html", pageData(c, "Login", nil))
		})
		publicRoutes.GET("/setup", initialSetup)
		publicRoutes.POST("/setup", saveInitialSetup)
		publicRoutes.GET("/register", func(c *gin.Context) {
			if getUser(c).IsValid() {
				c.Redirect(http.StatusSeeOther, "/")
			}
			c.HTML(http.StatusOK, "register.html", pageData(c, "Register", nil))
		})
		publicRoutes.POST("/register", register)
		publicRoutes.POST("/login", login)
	}

	authRoutes := publicRoutes.Group("/")
	authRoutes.Use(authRequired)
	{
		authRoutes.GET("/logout", logout)
		authRoutes.GET("/users/:username", func(c *gin.Context) {
			username := c.Param("username")
			var userProfile = &User{}
			result := db.Limit(1).Find(userProfile, "username = ?", username)
			if result.Error != nil {
				c.HTML(http.StatusOK, "profile.html", pageData(c, "Profile", gin.H{"error": result.Error}))
				return
			}
			c.HTML(http.StatusOK, "profile.html", pageData(c, "Profile", gin.H{"userProfile": userProfile}))
		})
		authRoutes.POST("/users/:id", editProfile)
		authRoutes.GET("/events/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.PostForm("id"))
			if err != nil {
				c.HTML(http.StatusBadRequest, "index.html", pageData(c, "Events", gin.H{"error": err}))
				return
			}
			event := &Event{}
			result := db.First(event, "id = ?", id)
			if result.Error != nil {
				c.HTML(http.StatusOK, "index.html", pageData(c, "Events", gin.H{"error": result.Error}))
				return
			}

			attendees := &[]User{}
			result = db.Where("event_id = ?", id).Find(event)
			if result.Error != nil {
				c.HTML(http.StatusOK, "index.html", pageData(c, "Events", gin.H{"error": result.Error}))
				return
			}

			c.HTML(http.StatusOK, "index.html", pageData(c, "Events", gin.H{"event": event, "attendees": attendees}))
		})
		authRoutes.GET("/join/:id", joinEvent)
		authRoutes.GET("/export/:id", exportProfile)
	}

	apiRoutes := publicRoutes.Group("/api/")
	apiRoutes.Use(apiRequired)
	{
		apiRoutes.GET("/score", scoreInput)
		apiRoutes.GET("/comps", scoreInput)
	}

	adminRoutes := publicRoutes.Group("/")
	adminRoutes.Use(adminRequired)
	{
		adminRoutes.GET("/config", func(c *gin.Context) {
			secrets := []Secret{}
			result := db.Order("time desc").Find(&secrets)
			if result.Error != nil {
				c.HTML(http.StatusOK, "config.html", pageData(c, "Config", gin.H{"error": result.Error}))
				return
			}
			c.HTML(http.StatusOK, "config.html", pageData(c, "Config", gin.H{"secrets": secrets}))
		})
		adminRoutes.POST("/config", setConfig)
		adminRoutes.POST("/new", editCompetition)
		adminRoutes.GET("/new", func(c *gin.Context) {
			c.HTML(http.StatusOK, "compdata.html", pageData(c, "New Competition", gin.H{"secrets": nil}))
		})
		adminRoutes.GET("/add", manualScore)
		adminRoutes.POST("/add", processManualScore)
	}

	r.Run(":1337")
}
