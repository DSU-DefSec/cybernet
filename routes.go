package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func editProfile(c *gin.Context) {
	user := getUser(c)
	username := c.PostForm("username")

	userProfile := &User{}
	result := db.First(&userProfile, "username = ?", username)

	if result.Error != nil {
		c.HTML(http.StatusBadRequest, "profile.html", pageData(c, "Edit Profile", gin.H{"error": result.Error}))
		return
	}

	if !user.Admin && user.Username != userProfile.Username {
		c.HTML(http.StatusBadRequest, "profile.html", pageData(c, "Edit Profile", gin.H{"error": errors.New("Nice try bub.")}))
		return
	}

	c.ShouldBind(userProfile)
	result = db.Where("username = ?", username).Save(userProfile)
	if result.Error != nil {
		c.HTML(http.StatusBadRequest, "profile.html", pageData(c, "Edit Profile", gin.H{"error": result.Error}))
		return
	}

	c.HTML(http.StatusBadRequest, "profile.html", pageData(c, "Edit Profile", gin.H{"message": "Successfully saved.", "userProfile": userProfile}))
}

func joinEvent(c *gin.Context) {
	user := getUser(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "index.html", pageData(c, "Events", gin.H{"error": err}))
		return
	}

	var event = &Event{}
	result := db.First(event, "id = ?", id)
	if result.Error != nil {
		c.HTML(http.StatusBadRequest, "index.html", pageData(c, "Events", gin.H{"error": err}))
		return
	}

	signup := &SignUp{}
	result = db.Limit(1).Find(signup, "event_id = ? and user = ?", id, user.Username)
	if result.Error != nil {
		c.HTML(http.StatusBadRequest, "index.html", pageData(c, "Events", gin.H{"error": result.Error}))
		return
	}

	if signup.User != "" {
		c.HTML(http.StatusBadRequest, "index.html", pageData(c, "Events", gin.H{"error": errors.New("You already signed up!")}))
		return
	}

	if user.IALab == "" {
		c.HTML(http.StatusBadRequest, "index.html", pageData(c, "Events", gin.H{"error": errors.New("You need to configure your IALab username!")}))
		return
	}

	signup = &SignUp{
		Time:    time.Now(),
		User:    user.Username,
		EventID: uint(id),
	}
	result = db.Create(signup)
	if result.Error != nil {
		c.HTML(http.StatusBadRequest, "index.html", pageData(c, "Events", gin.H{"error": result.Error}))
		return
	}

	apiDeploy(event.VApp, "DefSec_Lessons", []string{user.IALab})

	c.HTML(http.StatusOK, "index.html", pageData(c, "Scoreboard", gin.H{"message": "Sign up successful!"}))
}

func exportProfile(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", pageData(c, "Scoreboard", gin.H{}))
}

func scoreInput(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", pageData(c, "Scoreboard", gin.H{}))
}

func manualScore(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", pageData(c, "Scoreboard", gin.H{}))
}
func processManualScore(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", pageData(c, "Scoreboard", gin.H{}))
}

func editCompetition(c *gin.Context) {
	id, err := strconv.Atoi(c.PostForm("id"))
	var comp *Event
	if err != nil || strings.TrimSpace(c.PostForm("id")) == "" {
		comp = &Event{}
	} else {
		// fetch comp
		comp = &Event{}
		print(id)
	}
	c.ShouldBind(comp)
	result := db.Create(comp)
	if result.Error != nil {
		c.HTML(http.StatusBadRequest, "compdata.html", pageData(c, "Edit Competition", gin.H{"error": result.Error}))
	}
	c.Redirect(http.StatusSeeOther, "/")
}

func setConfig(c *gin.Context) {
	secrets := []Secret{}
	result := db.Order("time desc").Find(&secrets)
	if result.Error != nil {
		c.HTML(http.StatusOK, "config.html", pageData(c, "Config", gin.H{"error": result.Error}))
		return
	}

	newSecret := c.PostForm("secret")
	var secret = &Secret{}
	result = db.First(secret, "secret = ?", newSecret)
	if result.Error == nil {
		c.HTML(http.StatusOK, "config.html", pageData(c, "Configuration", gin.H{"error": "Secret is not unique!", "secret": newSecret, "secrets": secrets}))
		return
	}

	createdSecret := Secret{
		Secret: newSecret,
		Time:   time.Now(),
	}

	result = db.Create(&createdSecret)
	if result.Error != nil {
		c.HTML(http.StatusBadRequest, "config.html", pageData(c, "Configuration", gin.H{"error": result.Error, "secret": newSecret, "secrets": secrets}))
		return
	}

	secrets = append(secrets, createdSecret)
	c.HTML(http.StatusOK, "config.html", pageData(c, "Configuration", gin.H{"secrets": secrets}))
}

func initialSetup(c *gin.Context) {
	var user User
	if res := db.First(&user); errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.HTML(http.StatusOK, "setup.html", pageData(c, "Initial Setup", gin.H{}))
		return
	}
	c.Redirect(http.StatusSeeOther, "/login")
}

func saveInitialSetup(c *gin.Context) {
	var user User
	if res := db.First(&user); errors.Is(res.Error, gorm.ErrRecordNotFound) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
			c.HTML(http.StatusBadRequest, "setup.html", pageData(c, "Setup", gin.H{"error": "Username or password can't be empty."}))
			return
		}

		givenPw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			c.HTML(http.StatusBadRequest, "setup.html", pageData(c, "Setup", gin.H{"error": err}))
			return
		}

		user := User{
			Admin:    true,
			Username: username,
			HashedPw: string(givenPw),
		}

		result := db.Create(&user)
		if result.Error != nil {
			c.HTML(http.StatusBadRequest, "setup.html", pageData(c, "Setup", gin.H{"error": err}))
		}
	}
	c.Redirect(http.StatusSeeOther, "/login")
}
func allAttendance(c *gin.Context) {
	var allAttendance []Attendance

	var events []Event
	result := db.Find(&events)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	for _, event := range events {
		var signups []SignUp
		result := db.Find(&signups, "event_id = ?", event.ID)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
			return
		}

		attendance := Attendance{
			Event: event.Title,
		}

		for _, signup := range signups {
			var user User
			result = db.First(&user, "username = ?", signup.User)
			if result.Error != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
				return
			}
			attendance.Attendees = append(attendance.Attendees, Attendee{
				Time:     signup.Time,
				Username: user.IALab,
			})
		}
		allAttendance = append(allAttendance, attendance)
	}

	c.JSON(http.StatusOK, gin.H{"events": allAttendance})
}

func eventAttendance(c *gin.Context) {
	var signups []SignUp
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	result := db.Find(&signups, "event_id = ?", id)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	var event Event
	result = db.Find(&event, "id = ?", id)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	attendance := Attendance{
		Event: event.Title,
	}

	for _, signup := range signups {
		var user User
		result = db.First(&user, "username = ?", signup.User)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
			return
		}
		attendance.Attendees = append(attendance.Attendees, Attendee{
			Time:     signup.Time,
			Username: user.IALab,
		})
	}
	c.JSON(http.StatusOK, gin.H{"attendance": attendance})
}

func pageData(c *gin.Context, title string, ginMap gin.H) gin.H {
	newGinMap := gin.H{}
	newGinMap["title"] = title
	newGinMap["user"] = getUser(c)
	for key, value := range ginMap {
		newGinMap[key] = value
	}
	return newGinMap
}
