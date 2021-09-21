package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)


var (
	// Hardcoded CT timezone
	location, _    = time.LoadLocation("America/Rainy_River")
	locString = "CT"
)

func main() {
	err := readConfig(dwConf)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "illegal config"))
	}

	// Initialize Gin router
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Add... add function
	r.SetFuncMap(template.FuncMap{
		"inc": func(x int) int {
			return x + 1
		},
	})

	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets")
	initCookies(r)

	// 404 handler
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", nil)
	})

	// Route definitons
	publicRoutes := r.Group("/")
	{
		routes.GET("", viewHome)
		routes.GET("/login", func(c *gin.Context) {
			if getUserOptional(c).IsValid() {
				c.Redirect(http.StatusSeeOther, "/")
			}
			c.HTML(http.StatusOK, "login.html", pageData(c, "Login", nil))
		})
		routes.POST("/login", login)
	}

	authRoutes := routes.Group("/")
	authRoutes.Use(authRequired)
	{
		authRoutes.GET("/logout", logout)
		authRoutes.GET("/profile/:id", viewProfile)
		authRoutes.GET("/attend/:id", attendEvent)
		authRoutes.GET("/export/:id", exportProfile)
	}

	apiRoutes := routes.Group("/")
	apiRoutes.Use(apiRequired)
	{
		apiRoutes.GET("/score", scoreInput)
	}

	adminRoutes := routes.Group("/")
	adminRoutes.Use(adminRequired)
	{
		adminRoutes.GET("/add", manualScore)
		adminRoutes.POST("/add", processManualScore)
		adminRoutes.GET("/export/:id", exportTeamData)
		adminRoutes.GET("/profile/:id", viewPCR)
	}

	r.Run(":80")
}
