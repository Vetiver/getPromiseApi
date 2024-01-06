package handlers

import (
	"getPromiseApi/db"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func sendConfirmationEmail(email string, code int) (string, error) {
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
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Confirmation Email")
	m.SetBody("text/html", "Спасибо за регистрацию, вот ваша ссылка на подтверждение: <a href=\"http://localhost:8000/register?code=попа\">http://localhost:8000/register?code=555</a>")

	d := gomail.NewDialer(smtpName, port, emailAdress, emailPass)

	if err := d.DialAndSend(m); err != nil {
		return "", err
	}

	return "отправилось", nil
}

func (h BaseHandler) SendMail(c *gin.Context) {
	sendConfirmationEmail("vetiverdev@gmail.com", 555)
	var user *db.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.Email =  "vetiverdev@gmail.com"
	user.Password = "3123123"
	user.ConfirmCode = 555
	

	c.JSON(http.StatusOK, gin.H{
        "result": "чет вышло",
    })
}
