package bot

import (
	"fmt"
	"tg_reader_bot/internal/serializer"

	"github.com/gotd/td/telegram/message/markup"
	"github.com/gotd/td/tg"
)

const (
	AddNewChannel = iota
	MyChannels
	AddNewKeyWord
	RemoveKeyWord
	NextChannels
	PrevChannels
	NextKeyWords
	PrevKeyWords
	ChannelInfo
	Back
	MainPage
)

type buttonChannelInfo struct {
	Name string
}

type buttonMenuBack struct {
	BackMenu int
}

type buttonRemoveKeyWord struct {
	ID int32
}

func CreateButton(name string, msgID uint32, data any) *tg.KeyboardButtonCallback {
	bytes, err := serializer.EncodeMessage(msgID, data)
	if err != nil {
		fmt.Println("CreateInlineButton error %v", err)
		return nil
	}

	return markup.Callback(
		name,
		bytes,
	)
}
