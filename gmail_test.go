package main

import (
	"testing"

	gmail "google.golang.org/api/gmail/v1"
)

func Test_pickSubject(t *testing.T) {

	srv := GmailService{}
	subject := srv.pickSubject(&gmail.Message{
		Payload: &gmail.MessagePart{
			Headers: []*gmail.MessagePartHeader{
				{Name: "X-Server", Value: "remote"},
				{Name: "X-Antivirus", Value: "portshield"},
				{Name: "Subject", Value: "cloud watch"},
			},
		},
	})

	if subject == "" {
		t.Error("subject is empty")
	}

	if subject != "cloud watch" {
		t.Errorf("want: 'cloud watch', got: '%s'", subject)
	}
}

func Test_messageIDs(t *testing.T) {

	srv := GmailService{}
	ids := srv.reverseMessageIDs([]*gmail.Message{
		{Id: "1"},
		{Id: "2"},
		{Id: "3"},
		{Id: "4"},
		{Id: "5"},
	})

	if len(ids) != 5 {
		t.Errorf("lenght error want: 5, got: %d", len(ids))
	}

	if ids[0] != "5" {
		t.Errorf("reverse error want: 5, got: %s", ids[0])
	}
}
