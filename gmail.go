package main

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gmail "google.golang.org/api/gmail/v1"
)

var (
	ErrTokenFileExists      = errors.New("token file does not exist")
	ErrCredentialFileExists = errors.New("credential file does not exists")
)

type GmailService struct {
	credentialFile string
	tokenFile      string
	query          string
	lastHistoryId  uint64
	config         *oauth2.Config
	token          *oauth2.Token
	service        *gmail.Service
}

func NewGmailService(credentialfile, tokenfile, query string) (srv *GmailService, err error) {

	srv = &GmailService{
		credentialFile: credentialfile,
		tokenFile:      tokenfile,
		query:          query,
	}

	if !srv.isFileExists(credentialfile) {
		return srv, ErrCredentialFileExists
	}

	cfg, err := srv.newConfigFromJSON()
	if err != nil {
		return srv, err
	}
	srv.config = cfg

	if !srv.isFileExists(tokenfile) {
		return srv, ErrTokenFileExists
	}

	token, err := srv.getTokenFromJSON()
	if err != nil {
		return srv, err
	}

	client := cfg.Client(context.Background(), &token)
	service, err := gmail.NewService(client)
	if err != nil {
		return srv, err
	}

	srv.token = &token
	srv.service = service

	return srv, nil
}

func (g *GmailService) isFileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func (g *GmailService) newConfigFromJSON() (cfg *oauth2.Config, err error) {

	filedata, err := ioutil.ReadFile(g.credentialFile)
	if err != nil {
		return cfg, err
	}

	cfg, err = google.ConfigFromJSON(filedata, gmail.GmailReadonlyScope)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func (g *GmailService) getTokenFromJSON() (token oauth2.Token, err error) {

	filedata, err := os.Open(g.tokenFile)
	if err != nil {
		return token, err
	}

	err = json.NewDecoder(filedata).Decode(&token)
	return token, err
}

func (g *GmailService) ReadSubjects(userID string) (subjects []string, err error) {

	var message *gmail.Message

	subjects = []string{}
	messageFormat := "metadata"

	rsp, err := g.service.Users.Messages.List(userID).Q(g.query).Do()
	if err != nil {
		return subjects, err
	}

	for _, id := range g.reverseMessageIDs(rsp.Messages) {

		message, err = g.service.Users.Messages.Get(userID, id).Format(messageFormat).Do()
		if err != nil {
			break
		}

		if g.lastHistoryId >= message.HistoryId {
			continue
		}

		subject := g.pickSubject(message)
		if subject == "" {
			continue
		}

		subjects = append(subjects, subject)
		if g.lastHistoryId < message.HistoryId {
			g.lastHistoryId = message.HistoryId
		}

		log.Printf("historyID: %d, Subject: %s, lastHistoryId: %d\n", message.HistoryId, subject, g.lastHistoryId)
	}

	return subjects, nil
}

func (g *GmailService) pickSubject(message *gmail.Message) (subject string) {

	subjectHeaderKeyName := "subject"

	if message.Payload == nil {
		return ""
	}

	for i := 0; i < len(message.Payload.Headers); i++ {
		key := message.Payload.Headers[i].Name
		if strings.ToLower(key) == subjectHeaderKeyName {
			subject = message.Payload.Headers[i].Value
			break
		}
	}

	return subject
}

func (g *GmailService) reverseMessageIDs(messages []*gmail.Message) (ids []string) {

	ids = []string{}

	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	for i := 0; i < len(messages); i++ {
		ids = append(ids, messages[i].Id)
	}

	return ids
}
