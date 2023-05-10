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
const getServerMethod string = "groups.getLongPollServer"
const sendMessageMethod string = "messages.send"
const apiVersion string = "5.131"

type Bot struct {
	session     LPSession
	accessToken string
	wait        int
}

// TODO: internet
func (b *Bot) sendMessage(to int, text string, k Keyboard) {
	keyboardParam, _ := json.Marshal(k)

	parameters := url.Values{}
	parameters.Set("user_id", strconv.Itoa(to))
	parameters.Set("message", text)
	//parameters.Set("group_id", "220417305")
	parameters.Set("access_token", b.accessToken)
	parameters.Set("v", apiVersion)
	parameters.Set("random_id", strconv.FormatInt(time.Now().UnixMilli(), 10))
	parameters.Set("keyboard", string(keyboardParam))
	resp, err := http.PostForm(apiAddress+sendMessageMethod, parameters)
	if err != nil {
		log.Fatal(err.Error())
	}
	body, _ := io.ReadAll(resp.Body)
	log.Println(body)

}

func (b *Bot) poll() ([]Update, error) {
	// building query
	parameters := url.Values{}
	parameters.Set("act", "a_check")
	parameters.Set("key", b.session.Key)
	parameters.Set("ts", b.session.EventNum)
	parameters.Set("wait", string(b.wait))

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
		case "message_event":
			err := b.handleChoice(u.Object)
			if err != nil {
				log.Printf("unable to handle button click: %s", err.Error())
			}
		}
	}
}

func NewVkBot(token string, groupId string) (*Bot, error) {
	// building long poll server request
	parameters := url.Values{}
	parameters.Set("access_token", token)
	parameters.Set("v", apiVersion)
	parameters.Set("group_id", groupId)

	// requesting
	resp, err := http.PostForm(apiAddress+getServerMethod, parameters)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("getting long poll server error: %w", err)
	}

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
	}, nil
}
