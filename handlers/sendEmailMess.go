package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"math/rand"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)
type EmailAndPassword struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func generateRandomCode() int {
	rand.NewSource(time.Now().UnixNano())
	return rand.Intn(90000) + 10000
}


func (h *BaseHandler) sendConfirmationEmail(reqData *EmailAndPassword, code int) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	emailAdress := os.Getenv("EMAIL_ADDRESS")
	emailPass := os.Getenv("EMAIL_PASSWORDCONF")
	smtpName := os.Getenv("SMTP")
	port, err := strconv.Atoi(os.Getenv("PORT"))
	m := gomail.NewMessage()
	m.SetHeader("From", emailAdress)
	m.SetHeader("To", reqData.Email)
	m.SetHeader("Subject", "Confirmation Email")
	m.SetBody("text/html", fmt.Sprintf("Спасибо за регистрацию, вот ваша ссылка на подтверждение: <a href=\"http://localhost:8000/register?code=%d\">http://localhost:8000/register?code=%d</a>", code, code))
	h.Code[reqData.Email] = code
	log.Printf("Code for user %s: %d\n", reqData.Email, code)
	d := gomail.NewDialer(smtpName, port, emailAdress, emailPass)

	if err := d.DialAndSend(m); err != nil {
		return "", err
	}

	return "отправилось", nil
}



func (h BaseHandler) SendMail(c *gin.Context) {
	var reqData *EmailAndPassword
	code := generateRandomCode()
	if err := c.BindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.sendConfirmationEmail(reqData, code)

	c.JSON(http.StatusOK, gin.H{
		"result": "чет вышло",
	})
}
