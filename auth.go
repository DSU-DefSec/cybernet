package main

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

// getUUID returns a randomly generated UUID
func getUUID() string {
	return uuid.New().String()
}

// initCookies use gin-contrib/sessions{/cookie} to initalize a cookie store.
func initCookies(r *gin.Engine) {
	//r.Use(sessions.Sessions("cybernet", cookie.NewStore([]byte(getUUID()))))
	r.Use(sessions.Sessions("cybernet", cookie.NewStore([]byte("abcd"))))
}

// login is a handler that parses a form and checks for specific data
func login(c *gin.Context) {
	session := sessions.Default(c)

	username := c.PostForm("username")
	password := c.PostForm("password")

	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.HTML(http.StatusBadRequest, "login.html", pageData(c, "login", gin.H{"error": "Username or password can't be empty."}))
		return
	}

	var user = &User{}
	db.First(user, "username = ?", username)
	if !user.IsValid() {
		c.HTML(http.StatusBadRequest, "login.html", pageData(c, "login", gin.H{"error": errors.New("Invalid username or password.")}))
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.HashedPw), []byte(password))
	if err != nil {
		c.HTML(http.StatusBadRequest, "login.html", pageData(c, "login", gin.H{"error": errors.New("Invalid username or password.")}))
		return
	}

	session.Set("user", username)
	if err := session.Save(); err != nil {
		c.HTML(http.StatusBadRequest, "login.html", pageData(c, "login", gin.H{"error": "Failed to save session."}))
		return
	}

	c.Redirect(http.StatusSeeOther, "/")
}

func register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	secret := c.PostForm("secret")

	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.HTML(http.StatusBadRequest, "register.html", pageData(c, "Regster", gin.H{"error": "Username or password can't be empty.", "username": username, "secret": secret}))
		return
	}

	var testUser = &User{}
	db.First(testUser, "username = ?", username)
	if testUser.IsValid() {
		c.HTML(http.StatusBadRequest, "register.html", pageData(c, "login", gin.H{"error": errors.New("That username is taken."), "username": username, "secret": secret}))
		return
	}

	var foundSecret = &Secret{}
	db.Order("time desc").First(foundSecret)
	if foundSecret.Secret != secret || foundSecret.Time.IsZero() {
		db.Find(foundSecret, "secret = ?", secret)
		if foundSecret.Time.IsZero() {
			c.HTML(http.StatusBadRequest, "register.html", pageData(c, "Register", gin.H{"error": errors.New("Invalid secret!"), "username": username, "secret": secret}))
		} else {
			c.HTML(http.StatusBadRequest, "register.html", pageData(c, "Register", gin.H{"error": errors.New("Sorry, that secret is expired!"), "username": username, "secret": secret}))
		}
		return
	}

	givenPw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.HTML(http.StatusBadRequest, "register.html", pageData(c, "Register", gin.H{"error": err, "username": username, "secret": secret}))
		return
	}

	user := User{
		Username: username,
		HashedPw: string(givenPw),
	}

	result := db.Create(&user)
	if result.Error != nil {
		c.HTML(http.StatusBadRequest, "register.html", pageData(c, "Register", gin.H{"error": err, "username": username, "secret": secret}))
	}

	c.HTML(http.StatusOK, "login.html", pageData(c, "Login", gin.H{"message": "Registration successfull! Please log in."}))
}

func getUser(c *gin.Context) *User {
	session := sessions.Default(c)
	username := session.Get("user")
	var user = &User{}
	db.Limit(1).Find(user, "username = ?", username)
	return user
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	if user == nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}
	session.Delete("user")
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.Redirect(http.StatusSeeOther, "/")
}

// authRequired provides authentication middleware for ensuring that a user is logged in.
func authRequired(c *gin.Context) {
	user := getUser(c)
	if !user.IsValid() {
		c.Redirect(http.StatusSeeOther, "/login")
		c.Abort()
	}
	c.Next()
}

func adminRequired(c *gin.Context) {
	user := getUser(c)
	if !user.Admin {
		c.Redirect(http.StatusSeeOther, "/login")
		errorOutAnnoying(c, errors.New("Non-admin user tried to access admin endpoint."))
		return
	}
	c.Next()
}

func apiRequired(c *gin.Context) {
	c.Next()
}
