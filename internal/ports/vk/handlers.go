package vk

import (
	"encoding/json"
	"errors"
	"fmt"
)

const helloMessage = "Здравствуйте! Пожалуйста, выберите следующее название для этого паблика"

func newCBButton(label string, color string) (cbb CBButton) {
	cbb.Color = color
	//body, _ := json.Marshal(ChangeNamePayload{cmd: "change_name", newTitle: label})

	cbb.Action = CallBack{
		Type:    "callback",
		Label:   label,
		Payload: fmt.Sprintf("{\"name\":\"%s\"}", label),
	}
	return
}

func (b *Bot) handleStartCmd(obj MessageObject) (err error) {
	k := Keyboard{
		OneTime: false,
		Inline:  true,
	}
	k.Buttons = make([][]CBButton, 4)
	k.Buttons[0] = append(k.Buttons[0], newCBButton("(фан-страница)", "primary"))
	k.Buttons[1] = append(k.Buttons[1], newCBButton("Клуб любителей сусликов", "primary"))
	k.Buttons[2] = append(k.Buttons[2], newCBButton("Веду душу к богу", "primary"))
	k.Buttons[3] = append(k.Buttons[3], newCBButton("Фантазии пришёл конец", "primary"))
	b.sendMessage(obj.Message.FromId, helloMessage, k)
	return nil
}

func (b *Bot) handleTextMessage(obj MessageObject) (err error) {
	// TODO: sendMessage
	return nil
}

const startCmd string = "{\"command\":\"start\"}"

func (b *Bot) handleMessage(rawObject json.RawMessage) (err error) {
	var obj MessageObject
	err = json.Unmarshal(rawObject, &obj)
	if err != nil {
		return errors.New("unable to marshal object field")
	}

	// if there is a command
	if obj.Message.Payload == startCmd {
		//cmd, err := parsePayload(obj.Message.Payload)
		if err != nil {
			return errors.New("unable to parse payload")
		}
		switch obj.Message.Payload {
		case startCmd:
			if err = b.handleStartCmd(obj); err != nil {
				return fmt.Errorf("unable to handle start message: [%w]", err)
			}
		default:
			return errors.New("unsupported command provided")
		}
	} else { // if there is no command
		if err = b.handleTextMessage(obj); err != nil {
			return fmt.Errorf("unable to handle text message: [%w]", err)
		}
	}
	return nil
}

func (b *Bot) AskConfirmation(name string, userId int) {
	k := Keyboard{
		OneTime: false,
		Inline:  true,
	}
	k.Buttons = make([][]CBButton, 2)
	k.Buttons[0] = append(k.Buttons[0], newCBButton("Да, разумеется", "positive"))
	k.Buttons[1] = append(k.Buttons[1], newCBButton("Нет, я случайно нажал", "negative"))

	b.sendMessage(userId, fmt.Sprintf("Вы уверены, что вам нравится название \"%s\"?", name), k)
}

func (b *Bot) handleChoice(rawObject json.RawMessage) (err error) {
	var obj EventObject
	err = json.Unmarshal(rawObject, &obj)
	if err != nil {
		return errors.New("unable to marshal object field")
	}
	var nameWrapper map[string]string
	err = json.Unmarshal(obj.Payload, &nameWrapper)
	if err != nil {
		return errors.New("unable to get name wrapper")
	}

	name := nameWrapper["name"]

	b.AskConfirmation(name, obj.UserId)

	return nil
}
