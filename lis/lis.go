package lis

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
)

type LISConfig struct {
	endpoint  string
	username  string
	password  string
	groupname string
	day       string
	time      string
	details   string
}

func retConfig() *LISConfig {
	parser := argparse.NewParser("Let It Sleep", "Initial config for the tool")
	endpoint := parser.String("e", "endpoint", &argparse.Options{Help: "Endpoint of the API", Required: true})
	username := parser.String("u", "username", &argparse.Options{Required: true})
	password := parser.String("p", "password", &argparse.Options{Required: true})
	groupname := parser.String("g", "group", &argparse.Options{Required: true, Help: "Group ID in login form"})
	day := parser.String("d", "day", &argparse.Options{Required: true, Help: "Day to try book the slot"})
	time := parser.String("t", "time", &argparse.Options{Required: true, Help: "Time Slot to try book"})
	details := parser.String("s", "description", &argparse.Options{Help: "Comment for your booking", Default: "To Play"})
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Printf("Error on parsing: %s", err)
		os.Exit(1)
	}
	return &LISConfig{
		endpoint:  *endpoint,
		username:  *username,
		password:  *password,
		groupname: *groupname,
		day:       *day,
		time:      *time,
		details:   *details,
	}
}

func Book() {
	config := retConfig()
	instance := NewInstance(
		config.endpoint,
		config.username,
		config.password,
		config.groupname,
	)
	err := instance.Authorise()
	if err != nil {
		fmt.Printf("Failed to authorise user: %s\n", err)
		os.Exit(1)
	}
	session, err := NewSchedule(instance)
	if err != nil {
		fmt.Printf("Failed on making new session: %s", err.Error())
		os.Exit(1)
	}
	session.Refresh()

	booked := session.BookIfPossible(config.day, config.time, config.details)
	if booked != nil {
		fmt.Printf("Booked: %s", *booked)
		os.Exit(0)
	} else {
		fmt.Print("Failed with booking")
		os.Exit(2)
	}
}
