package gmail

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type GmailService struct {
	service *gmail.Service
}

func NewGmailService(accessToken string) (*GmailService, error) {
	ctx := context.Background()
	service, err := gmail.NewService(ctx, option.WithTokenSource(
		oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken}),
	))
	if err != nil {
		return nil, err 
	}

	return &GmailService{service: service}, nil

}

func (g *GmailService) ListMessages(targetEmail string) ([]*gmail.Message, error) {
 query:= fmt.Sprintf("from:%s", targetEmail)
 messages, err := g.service.Users.Messages.List("me").Q(query).MaxResults(10).Do()
 if err != nil {
	return nil, fmt.Errorf("failed to list messages: %w", err)
 }

 var fullMesages []*gmail.Message

 for _, msg := range messages.Messages {
	messages, err := g.getMessage(msg.Id)
	if err != nil {
		continue
	}
	fullMesages = append(fullMesages, messages)
 }

 return fullMesages,nil
}


func (g *GmailService)getMessage(messageId string) (*gmail.Message, error) {
	message, err := g.service.Users.Messages.Get("me", messageId).Do()

	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	return message, nil
}

