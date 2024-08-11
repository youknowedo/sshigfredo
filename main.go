package main

import (
	"fmt"
	"io"
	"log"
	"net/smtp"
	"os"
	"regexp"

	"github.com/fatih/color"
	"github.com/gliderlabs/ssh"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	helpCmd = regexp.MustCompile(`^help$`)
	contact = regexp.MustCompile(`^contact$`)
	exitCmd = regexp.MustCompile(`^exit$`)

	validEmail = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	auth := smtp.PlainAuth("", os.Getenv("SMTP_EMAIL"), os.Getenv("SMTP_PASSWORD"), os.Getenv("SMTP_HOST"))

	ssh.Handle(func(s ssh.Session) {
		color.New(color.FgYellow).Fprint(s, "\nSigfredo\n\n")
		io.WriteString(s, "type help to get started\n")

		// TODO: Make something fun here!

		term := terminal.NewTerminal(s, ("=> "))
		for {
			line, err := term.ReadLine()
			if err != nil {
				break
			}

			if len(line) > 0 {
				switch {
				case exitCmd.MatchString(line):
					s.Exit(0)
				case helpCmd.MatchString(line):
					color.New(color.FgYellow).Fprint(s, "\nCommands:\n\n")
					io.WriteString(s, "  help - Shows this message\n")
					io.WriteString(s, "  contact - Hits me up. On mail but hits me up nonetheless.\n")
					io.WriteString(s, "  exit - Exits this session\n\n")
				case contact.MatchString(line):
					term.SetPrompt(": ")

					// Aske for email until it's valid
					var email string
					for {
						io.WriteString(s, "   Email")
						email, err = term.ReadLine()
						if err != nil {
							break
						}

						if validEmail.MatchString(email) {
							break
						}

						color.New(color.FgRed).Fprint(s, "Invalid email\n\n")
					}
					io.WriteString(s, "   Message")
					message, err := term.ReadLine()
					if err != nil {
						break
					}

					err = smtp.SendMail(os.Getenv("SMTP_HOST")+":"+os.Getenv("SMTP_PORT"), auth, os.Getenv("SMTP_EMAIL"), []string{"hello@sigfredo.fun"}, []byte("Subject: SSH from "+email+"\nTo: hello@sigfredo.fun\n\n"+message))
					if err != nil {
						color.New(color.FgRed).Fprint(s, "\nError sending message\n\n")
						fmt.Println(err)
					} else {
						color.New(color.FgGreen).Fprint(s, "\nMessage sent\n\n")
					}

					term.SetPrompt("=> ")
				}
			}
		}
	})

	log.Fatal(ssh.ListenAndServe(":2222", nil))
}
