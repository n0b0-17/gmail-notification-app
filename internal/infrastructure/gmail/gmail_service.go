package gmail

import (
	"context"
	"encoding/base64"
	"fmt"
	"gmail-notification-app/internal/domain/models"
	"regexp"
	"strconv"
	"strings"
	"time"

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
	message, err := g.service.Users.Messages.Get("me", messageId).Format("full").Do()

	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	return message, nil
}

func ConvertToRecievedEmail(message *gmail.Message) *models.RecievedEmail {
	var fromEmail,toEmail string

	recivedAt := time.Unix(message.InternalDate/1000, 0)

	//ヘッダー情報の取得
	for _, header := range message.Payload.Headers {
		var emailRegex = regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
		switch header.Name{
		case "From":
			//正規表現を用いてEmailアドレスを取得する
			//memo:Type Messageの中にはEmailアドレスのみを値として持つフィールドがないらしい。。←そんなことある？
			if matches := emailRegex.FindString(header.Value); matches != "" {
				fromEmail = matches
			}
		case "To":
			if matches := emailRegex.FindString(header.Value); matches != "" {
				toEmail = matches
			}
		}
	}

	//メール本文の取得
	content := getMessageContent(message.Payload)
	//メールの本文から引落金額を取得
	amount := extractAmount(content)

	return &models.RecievedEmail{
		FromEmail: fromEmail,
		ToEmail: toEmail,
		RecievedAt: recivedAt,
		Contnet: content,
		Amount: int64(amount),
	}
}

func getMessageContent(payload *gmail.MessagePart) string {
	if payload.Body != nil && payload.Body.Data != ""{
		if data, err := base64.URLEncoding.DecodeString(payload.Body.Data); err == nil {
			return string(data)
		}
	}

	if payload.Parts != nil {
		for _, part := range payload.Parts {
			if part.MimeType == "text/plain" {
				if data,  err := base64.URLEncoding.DecodeString(part.Body.Data); err == nil {
					return string(data)
				}
			}
			if content := getMessageContent(part); content != ""{
				return content
			}
		}
	}
	return ""
}


func extractAmount(emailContent string)int{
	// 引落金額の行を探す正規表現
	re := regexp.MustCompile(`引落金額\s*：\s*([\d,]+\.?\d*)`)

	matches := re.FindStringSubmatch(emailContent)

	if len(matches) < 2 {
		return 0
	}

	amountStr := strings.ReplaceAll(matches[1],",","")

	amount,err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return 0
	}

	return int(amount)

}

