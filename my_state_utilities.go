package viki

import (
	"errors"
	"flag"
	"fmt"

	mailgun "github.com/mailgun/mailgun-go"
)

// getMochadState returns the string value of state of the object named name
// after asserting.
func (m *Viki) getMochadState(name string) string {
	_, o := m.ObjectManager.GetObjectByName(name)
	if o == nil {
		return ""
	}
	st, ok := o.State.(string)
	if !ok {
		return ""
	}
	return st
}

// getModeState returns the string value of state of the object named name
// after asserting.
func (m *Viki) getModeState(name string) string {
	_, o := m.ObjectManager.GetObjectByName(name)
	if o == nil {
		return ""
	}
	st, ok := o.State.(string)
	if !ok {
		return ""
	}
	return st
}

func (m *Viki) quickMail(recipients []string, msg string) error {
	domain := flag.Lookup("mg_domain")
	apikey := flag.Lookup("mg_apikey")
	pubkey := flag.Lookup("mg_pubkey")

	if domain == nil || apikey == nil || pubkey == nil {
		return errors.New("Failed to send email. Mailgun flags not set")
	}
	mg_domain := domain.Value.String()
	mg_apikey := apikey.Value.String()
	mg_pubkey := pubkey.Value.String()

	mg := mailgun.NewMailgun(mg_domain, mg_apikey, mg_pubkey)
	message := mailgun.NewMessage("home@fulton-ave", "Alert!", msg)
	for _, r := range recipients {
		if err := message.AddRecipient(r); err != nil {
			return err
		}
	}
	msg, id, err := mg.Send(message)

	if err != nil {
		return errors.New(fmt.Sprintf("Could not send message: %v, ID %v, %+v", err, id, msg))
	}
	return nil
}
