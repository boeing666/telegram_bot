package bot

import (
	protos "tg_reader_bot/internal/protobufs"
	"time"

	"github.com/gotd/td/telegram/message/markup"
	"github.com/gotd/td/tg"
	"google.golang.org/protobuf/proto"
)

func buildInitalMenu() tg.ReplyMarkupClass {
	return markup.InlineRow(
		CreateButton(
			"Добавить канал",
			uint32(protos.MessageID_AddNewChannel),
			nil,
		),
		CreateButton(
			"Мои каналы",
			uint32(protos.MessageID_MyChannels),
			nil,
		),
	)
}

func CreateButton(name string, msgID uint32, data proto.Message) *tg.KeyboardButtonCallback {
	msg := []byte{}
	if data != nil {
		msg, _ = proto.Marshal(data)
	}

	header := protos.MessageHeader{Time: uint64(time.Now().Unix()), Msgid: msgID, Msg: msg}
	result, _ := proto.Marshal(&header)
	return markup.Callback(
		name,
		result,
	)
}

func CreateBackButton(name string, backMenuID uint32, msg proto.Message) tg.KeyboardButtonRow {
	button := &protos.ButtonMenuBack{Newmenu: backMenuID}
	if msg != nil {
		bytes, _ := proto.Marshal(msg)
		button.Msg = bytes
	}

	return tg.KeyboardButtonRow{
		Buttons: []tg.KeyboardButtonClass{
			CreateButton(
				name,
				uint32(protos.MessageID_Back),
				button,
			),
		},
	}
}
