package vk

import (
	"encoding/json"
	"errors"
	"fmt"
)

const helloMessage = "Здравствуйте! Пожалуйста, выберите следующее название для этого паблика"
const confirmationMessageFormat = "Вы уверены, что вам нравится название \"%s\"?"
const congratsMessage = "Спасибо большое, что поменяли название! Но вы можете сделать это ещё раз!"
const againMessage = "Пожалуйста, помогите выбрать название."
const changeCmd = "changeName"
const confirmCmd = "confirm"
const startCmd = "start"
const yesOption = "Да, разумеется"
const noOption = "Нет, я случайно"

func (b *Bot) handleMessage(rawObject json.RawMessage) (err error) {
	var obj MessageObject
	err = json.Unmarshal(rawObject, &obj)
	if err != nil {
		return errors.New("unable to unmarshal object field")
	}

	// handling plain message
	if obj.Message.Payload == "" {
		return b.handleTextMessage(obj)
	}

	// handling messages with commands
	payload := make(map[string]string)
	err = json.Unmarshal(json.RawMessage(obj.Message.Payload), &payload)
	if err != nil {
		return err
	}

	switch payload["command"] {
	case startCmd:
		err = b.handleStartCmd(obj)
	case changeCmd:
		err = b.handleChoice(obj, payload["data"])
	case confirmCmd:
		err = b.handleConfirm(obj, payload["data"], payload["title"])
	default:
		return errors.New("unsupported command")
	}
	if err != nil {
		return
	}
	return nil
}

func (b *Bot) handleStartCmd(obj MessageObject) (err error) {
	return b.sendKeyboard(obj.Message.FromId, helloMessage, newTitlesKeyboard())
}

func (b *Bot) handleTextMessage(obj MessageObject) (err error) {
	return b.sendMessage(obj.Message.FromId, againMessage)
}

func (b *Bot) handleConfirm(obj MessageObject, option string, title string) (err error) {
	switch option {
	case yesOption:
		err = b.changeGroupTitle(title)
		if err != nil {
			return
		}
		err = b.sendKeyboard(obj.Message.FromId, congratsMessage, newTitlesKeyboard())
	case noOption:
		err = b.sendKeyboard(obj.Message.FromId, againMessage, newTitlesKeyboard())
	default:
		return errors.New("incorrect data in payload")
	}
	if err != nil {
		return
	}
	return
}

func (b *Bot) handleChoice(obj MessageObject, choice string) (err error) {
	err = b.sendKeyboard(
		obj.Message.FromId,
		fmt.Sprintf(confirmationMessageFormat, choice),
		newConfirmKeyboard(choice),
	)
	if err != nil {
		return err
	}
	return nil
}
