package main

import (
	"io"
	"log"
	"regexp"

	"github.com/fatih/color"
	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	helpCmd = regexp.MustCompile(`^help$`)
	contact = regexp.MustCompile(`^contact$`)
	exitCmd = regexp.MustCompile(`^exit$`)
)

func main() {
	ssh.Handle(func(s ssh.Session) {
		color.New(color.FgYellow).Fprint(s, "\nSigfredo\n\n")

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
					io.WriteString(s, "  help\n")
					io.WriteString(s, "  contact\n")
					io.WriteString(s, "  exit\n\n")
				case contact.MatchString(line):
					term.SetPrompt("yru")
					io.WriteString(s, "  Email:")
					email, err := term.ReadLine()
					if err != nil {
						break
					}
					io.WriteString(s, "  Message:")
					message, err := term.ReadLine()
					if err != nil {
						break
					}

					color.New(color.FgYellow).Fprint(s, "\n\nThanks for contacting me!\n\n")
					io.WriteString(s, "  Email: "+email+"\n")
					io.WriteString(s, "  Message: "+message+"\n\n")
				}
			}
		}
	})

	log.Fatal(ssh.ListenAndServe(":2222", nil))
}
