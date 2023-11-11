package mailing

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func EmailVerify(email string) error {
	// TODO: implement email verification
	err := godotenv.Load(".env")
	key := os.Getenv("SENDGRID_API_KEY")

	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(99999) + 10000

	from := mail.NewEmail("Securelee", "noreply@securelee.tech")
	subject := "Welcome to Securelee Vault! [Verification]"
	to := mail.NewEmail("", email)
	plainTextContent := "Verify your email address. Hello, You have selected this email address as your Securelee ID. To verify that it's you, enter the code below on the email verification terminal : " + strconv.Itoa(code) + " Best Regards, Securelee Team"
	htmlContent := "<strong><p align=center>Verify your email address.</p></strong> <p>Hello,</p> <p>You have selected this email address as your Securelee ID. To verify that it's you, enter the code below on the email verification terminal : </p> <p><strong>" + strconv.Itoa(code) + "</strong></p> <br> <p>Best Regards,</p> Securelee Team"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(key)
	_, err = client.Send(message)
	if err != nil {
		return err
	}

	var check int
	fmt.Print("\n> To Verify the Email Address enter the code sent on the entered email address : ")
	fmt.Scanf("%d\n", &check)
	if check != code {
		fmt.Println("Invalid Code Entered!")
		os.Exit(1)
	}

	return nil
}

func SendMail(name string, email string) error {
	err := godotenv.Load(".env")
	from := mail.NewEmail("Securelee", "noreply@securelee.tech")
	key := os.Getenv("SENDGRID_API_KEY")

	subject := "Welcome to Securelee Vault!"

	to := mail.NewEmail(name, email)

	plainTextContent := "Hey " + name + ", Welcome to Securelee! We're glad you're here."

	htmlContent := "<strong>Hey " + name + "</strong>,<br> <p> Welcome to Securelee! </p> <p> We're glad you're here.<p> If you have any queries regarding Securelee CLI, feel free to reach out to us at enquiry@securelee.tech . </p> <p> We're here with you every step of the way.</p> <br> <p> Best Regards,</p> <p>Securelee Team</p>"

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient(key)

	response, err := client.Send(message)
	_ = response
	if err != nil {
		log.Println(err)
	}

	return nil
}
