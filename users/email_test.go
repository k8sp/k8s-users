package users

import (
	"testing"

	"github.com/go-gomail/gomail"
	"github.com/stretchr/testify/assert"
)

func TestNewSmtpInfo(t *testing.T) {
	e := SmtpInfo{
		ESMTPServer: "email.test.com",
		AdminEmail:  "admin@domain.com",
		AdminSecrt:  "admin",
		subject:     "k8s key and crt",
		text:        "Successful: your crt and key for k8s are in attachment!\n",
		ch:          make(chan *gomail.Message),
	}

	smtp := NewSmtpInfo("email.test.com", "admin@domain.com", "admin")
	assert.Equal(t, e.ESMTPServer, smtp.ESMTPServer, "should be equal")
	assert.Equal(t, e.AdminEmail, smtp.AdminEmail, "should be equal")
	assert.Equal(t, e.AdminSecrt, smtp.AdminSecrt, "should be equal")
}
