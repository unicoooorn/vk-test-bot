package vk

import "encoding/json"

type CallBack struct {
	Type    string `json:"type"`
	Label   string `json:"label"`
	Payload string `json:"payload"`
}

type ChangeNamePayload struct {
	cmd      string `json:"command"`
	newTitle string `json:"new_title"`
}

type CBButton struct {
	Action CallBack `json:"action"`
	Color  string   `json:"color"`
}

type Keyboard struct {
	OneTime bool         `json:"one_time"`
	Buttons [][]CBButton `json:"buttons"`
	Inline  bool         `json:"inline"`
}

type LPSession struct {
	Key      string `json:"key"`
	Server   string `json:"server"`
	EventNum string `json:"ts"`
}

type LPServerRequest struct {
	Session LPSession `json:"response"`
}

type Update struct {
	GroupId int             `json:"group_id"`
	Type    string          `json:"type"`
	EventId string          `json:"event_id"`
	V       string          `json:"v"`
	Object  json.RawMessage `json:"object"`
}

type LPResponse struct {
	Ts      string   `json:"ts"`
	Updates []Update `json:"updates"`
}

type MessageObject struct {
	Message struct {
		Date                  int           `json:"date"`
		FromId                int           `json:"from_id"`
		Id                    int           `json:"id"`
		Out                   int           `json:"out"`
		Attachments           []interface{} `json:"attachments"`
		ConversationMessageId int           `json:"conversation_message_id"`
		FwdMessages           []interface{} `json:"fwd_messages"`
		Important             bool          `json:"important"`
		IsHidden              bool          `json:"is_hidden"`
		Payload               string        `json:"payload"`
		PeerId                int           `json:"peer_id"`
		RandomId              int           `json:"random_id"`
		Text                  string        `json:"text"`
	} `json:"message"`
	ClientInfo struct {
		ButtonActions  []string `json:"button_actions"`
		Keyboard       bool     `json:"keyboard"`
		InlineKeyboard bool     `json:"inline_keyboard"`
		Carousel       bool     `json:"carousel"`
		LangId         int      `json:"lang_id"`
	} `json:"client_info"`
}

type EventObject struct {
	ConversationMessageId int             `json:"conversation_message_id"`
	UserId                int             `json:"user_id"`
	PeerId                int             `json:"peer_id"`
	EventId               string          `json:"event_id"`
	Payload               json.RawMessage `json:"payload"`
}
