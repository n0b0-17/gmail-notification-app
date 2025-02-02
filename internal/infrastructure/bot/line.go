package bot

import (
	"fmt"
	"log"

	"github.com/line/line-bot-sdk-go/linebot"
)



func NewLineBot (channelSecret string, channelToken string) (*linebot.Client, error) {
if channelSecret == "" || channelToken== "" {
		log.Printf("LINEBOT Client parameters - channelSecret empty:%v, ChannelSecret empty:%v",
			channelSecret == "",
			channelToken == "",
		)
		return nil,fmt.Errorf("invalid linebot parameters")
	}

	client, err := linebot.New(channelSecret, channelToken)
	if err != nil {
		return nil,fmt.Errorf("failed to create linebot client: %v",err)
	}

	return client,nil 
}
