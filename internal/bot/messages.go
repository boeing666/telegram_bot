package bot

import (
	"tg_reader_bot/internal/protobufs"
	"time"

	"github.com/gotd/td/telegram/message/markup"
	"github.com/gotd/td/tg"
	"google.golang.org/protobuf/proto"
)

func buildInitalMenu() tg.ReplyMarkupClass {
	return markup.InlineRow(
		CreateButton(
			"Добавить канал",
			protobufs.MessageID_AddNewChannel,
			nil,
		),
		CreateButton(
			"Мои каналы",
			protobufs.MessageID_MyChannels,
			nil,
		),
	)
}

func CreateButton(name string, msgID protobufs.MessageID, data proto.Message) *tg.KeyboardButtonCallback {
	var msg []byte
	if data != nil {
		msg, _ = proto.Marshal(data)
	}

	header := protobufs.MessageHeader{Time: uint64(time.Now().Unix()), Msgid: msgID, Msg: msg}
	result, _ := proto.Marshal(&header)
	return markup.Callback(
		name,
		result,
	)
}

func CreateButtonRow(name string, msgID protobufs.MessageID, data proto.Message) tg.KeyboardButtonRow {
	return tg.KeyboardButtonRow{
		Buttons: []tg.KeyboardButtonClass{
			CreateButton(name, msgID, data),
		},
	}
}

func CreateSpaceButton() *tg.KeyboardButtonCallback {
	return CreateButton(" ", protobufs.MessageID_Spacer, nil)
}

func CreateSpaceButtonRow() tg.KeyboardButtonRow {
	return tg.KeyboardButtonRow{
		Buttons: []tg.KeyboardButtonClass{
			CreateSpaceButton(),
		},
	}
}

func CreateBackButton(name string, backMenuID protobufs.MessageID, msg proto.Message) tg.KeyboardButtonRow {
	button := &protobufs.ButtonMenuBack{Newmenu: backMenuID}
	if msg != nil {
		bytes, _ := proto.Marshal(msg)
		button.Msg = bytes
	}

	return tg.KeyboardButtonRow{
		Buttons: []tg.KeyboardButtonClass{
			CreateButton(
				name,
				protobufs.MessageID_Back,
				button,
			),
		},
	}
}
