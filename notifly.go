package main

import (
	"log"
	"time"
	"flag"
	"github.com/0xAX/notificator"
	"golang.org/x/net/context"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var notify *notificator.Notificator

type signal struct {
	ID           int64  `json:"id"`           // Not used when creating a signal
	Pid          string `json:"pid"`          // DK5QPID
	ZoneID       string `json:"zoneId"`       // KEY_A, KEY_B, etc...
	Name         string `json:"name"`         // message title
	Message      string `json:"message"`      // message body
	Effect       string `json:"effect"`       // e.g. SET_COLOR, BLINK, etc...
	Color        string `json:"color"`        // color in hex format. E.g.: "#FF0044"
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getAuth() *github.Client {
	var accesstoken string

	flag.StringVar(&accesstoken, "oauth", "", "What's your OAuth Token?")

	flag.Parse()

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accesstoken})

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return client
}

func getNotif(client *github.Client) []*github.Notification {
	ctx := context.Background()
	notifs, _, err := client.Activity.ListNotifications(ctx, nil)
	checkErr(err)

	return notifs
}

func isNotification(client *github.Client) []*github.Notification {
	notifs := getNotif(client)

	return notifs
}

func sendSignal(notification *github.Notification) {
	notify = notificator.New(notificator.Options{
		DefaultIcon: "icon/default.png",
		AppName:     "Github Notifications",
	})

	notify.Push(*notification.Subject.Title, *notification.Reason, "/home/hrouille/.bgs/bg2.jpg", notificator.UR_NORMAL)
}

func main() {
	client := getAuth()
	isSignalSent := false

	for true {
		if notifications := isNotification(client); len(notifications) > 0 && !isSignalSent {
			sendSignal(notifications[len(notifications)-1])
			isSignalSent = true
		} else if !(len(notifications) == 0) && isSignalSent {
			isSignalSent = false
		}

		time.Sleep(time.Duration(2) * time.Second)
	}

}
