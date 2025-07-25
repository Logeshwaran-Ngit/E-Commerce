package background_job

import (
	"bytes"
	"html/template"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendOrderConfirmation(toEmail string, data map[string]string) error {
	tmpl, err := template.ParseFiles("templates/order_confirmation.html")
	if err != nil {
		log.Println("Template parse error:", err)
		return err
	}
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		log.Println("Template execution error:", err)
		return err
	}
	log.Println("Email body rendered:", body.String())
	from := mail.NewEmail("ECHAN Store", os.Getenv("FROM_EMAIL"))
	subject := "Your Order Confirmation - ECHAN Store"
	content := mail.NewContent("text/html", body.String())
	message := mail.NewV3Mail()
	message.SetFrom(from)
	message.Subject = subject
	message.AddContent(content)
	p := mail.NewPersonalization()
	customer := mail.NewEmail("Customer", toEmail)
	p.AddTos(customer)
	admin := mail.NewEmail("Store Admin", "nlogeshwaranece@gmail.com")
	p.AddCCs(admin)
	message.AddPersonalizations(p)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Println("SendGrid error:", err)
		return err
	}
	log.Println("SendGrid status code:", response.StatusCode)
	log.Println("SendGrid response body:", response.Body)
	log.Println("SendGrid response headers:", response.Headers)
	return nil
}
