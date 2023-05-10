package vk

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const apiAddress string = "https://api.vk.com/method/"
const apiVersion string = "5.131"

const getServerMethod string = "groups.getLongPollServer"
const sendMessageMethod string = "messages.send"
const editGroupMethod = "groups.edit"

type Bot struct {
	session     LPSession
	accessToken string
	wait        int
	groupId     int
}

func newTitleButton(label string, color string) (cbb Button) {
	body, _ := json.Marshal(map[string]string{
		"command": changeCmd,
		"data":    label,
		"title":   label,
	})

	cbb.Color = color
	cbb.Action = Action{
		Type:    "text",
		Label:   label,
		Payload: string(body),
	}
	return
}

func newConfirmButton(label string, color string, title string) (cbb Button) {
	body, _ := json.Marshal(map[string]string{
		"command": confirmCmd,
		"data":    label,
		"title":   title,
	})

	cbb.Color = color
	cbb.Action = Action{
		Type:    "text",
		Label:   label,
		Payload: string(body),
	}
	return
}

func newTitlesKeyboard() *Keyboard {
	k := &Keyboard{
		OneTime: true,
		Inline:  false,
	}
	k.Buttons = make([][]Button, 4)
	k.Buttons[0] = append(k.Buttons[0], newTitleButton("Шрек (фан-страница)", "primary"))
	k.Buttons[1] = append(k.Buttons[1], newTitleButton("Клуб любителей сусликов", "primary"))
	k.Buttons[2] = append(k.Buttons[2], newTitleButton("Веду душу к богу", "primary"))
	k.Buttons[3] = append(k.Buttons[3], newTitleButton("Фантазия закончилась", "primary"))
	return k
}

func newConfirmKeyboard(title string) *Keyboard {
	k := &Keyboard{
		OneTime: true,
		Inline:  false,
	}
	k.Buttons = make([][]Button, 1)
	k.Buttons[0] = append(k.Buttons[0], newConfirmButton(yesOption, "positive", title))
	k.Buttons[0] = append(k.Buttons[0], newConfirmButton(noOption, "negative", title))
	return k
}

func (b *Bot) changeGroupTitle(title string) (err error) {
	parameters := url.Values{}
	parameters.Set("group_id", strconv.Itoa(b.groupId))
	parameters.Set("title", title)
	parameters.Set("access_token", b.accessToken)
	parameters.Set("v", apiVersion)
	resp, err := http.PostForm(apiAddress+editGroupMethod, parameters)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	fmt.Println(body)
	return
}

func (b *Bot) sendKeyboard(to int, text string, k *Keyboard) (err error) {
	keyboardParam, err := json.Marshal(k)
	if err != nil {
		return errors.New("incorrect keyboard format")
	}

	parameters := url.Values{}
	parameters.Set("user_id", strconv.Itoa(to))
	parameters.Set("message", text)
	parameters.Set("access_token", b.accessToken)
	parameters.Set("v", apiVersion)
	parameters.Set("random_id", strconv.FormatInt(time.Now().UnixMilli(), 10))
	parameters.Set("keyboard", string(keyboardParam))
	_, err = http.PostForm(apiAddress+sendMessageMethod, parameters)
	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) sendMessage(to int, text string) (err error) {
	parameters := url.Values{}
	parameters.Set("user_id", strconv.Itoa(to))
	parameters.Set("message", text)
	parameters.Set("access_token", b.accessToken)
	parameters.Set("v", apiVersion)
	parameters.Set("random_id", strconv.FormatInt(time.Now().UnixMilli(), 10))
	_, err = http.PostForm(apiAddress+sendMessageMethod, parameters)
	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) poll() ([]Update, error) {
	// building query
	parameters := url.Values{}
	parameters.Set("act", "a_check")
	parameters.Set("key", b.session.Key)
	parameters.Set("ts", b.session.EventNum)
	parameters.Set("wait", strconv.Itoa(b.wait))

	// making request
	resp, err := http.DefaultClient.PostForm(b.session.Server, parameters)
	if err != nil {
		return nil, errors.New("http request error")
	}
	defer resp.Body.Close()

	// reading
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("unable to read updates")
	}

	// unmarshalling
	var lpResponse LPResponse
	err = json.Unmarshal(data, &lpResponse)
	if err != nil {
		return nil, errors.New("updates format is incorrect")
	}
	b.session.EventNum = lpResponse.Ts
	return lpResponse.Updates, nil
}

func (b *Bot) getUpdatesChan() <-chan Update {
	ch := make(chan Update)
	go func() {
		for {
			updates, err := b.poll()
			if err != nil {
				log.Fatal(fmt.Errorf("unable to poll server: %w", err))
			}

			for _, u := range updates {
				ch <- u
			}
		}
	}()
	return ch
}

/*
func parsePayload(p string) (string, error) {
	var res string
	if n, err := fmt.Sscanf(p, payloadFormat, &res); err != nil {
		return "", err
	} else if n < 1 {
		return "", errors.New("unable to parse command")
	}
	return res, nil
}
*/

func (b *Bot) PollAndServe() {
	for u := range b.getUpdatesChan() {
		switch u.Type {
		case "message_new":
			err := b.handleMessage(u.Object)
			if err != nil {
				log.Printf("unable to handle message: %s", err.Error())
			}
		default:
			log.Printf("unsupported update: %s\n", u.Type)
		}
	}
}

func NewVkBot(token string, groupId int) (*Bot, error) {
	// building long poll server request
	parameters := url.Values{}
	parameters.Set("access_token", token)
	parameters.Set("v", apiVersion)
	parameters.Set("group_id", strconv.Itoa(groupId))

	// requesting
	resp, err := http.PostForm(apiAddress+getServerMethod, parameters)
	if err != nil {
		return nil, fmt.Errorf("getting long poll server error: %w", err)
	}
	defer resp.Body.Close()

	// processing response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("incorrect long poll server settings: %w", err)
	}
	var reqResponse LPServerRequest
	err = json.Unmarshal(body, &reqResponse)
	if err != nil {
		return nil, fmt.Errorf("incorrect long poll server settings: %w", err)
	}
	return &Bot{
		session:     reqResponse.Session,
		accessToken: token,
		groupId:     groupId,
	}, nil
}
