package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

const keypairCommand string = "/keypair"
const createAccountCommand string = "/create_account"
const startCommand string = "/start"

const telegramApiURL string = "https://api.telegram.org/bot"
const telegramToken string = "BOT_TOKEN" // No external exposure
const telegramSendMessage string = "/sendMessage"

var telegramApi string = telegramApiURL + os.Getenv(telegramToken) + telegramSendMessage

type Message struct {
	text string `json:"text"`
	chat Chat   `json:"chat"`
}

type Chat struct {
	id int `json:"id"`
}

//Update is a Telegram object that the handler receives every time an user interacts with the bot
type Update struct {
	updateID int     `json:"update_id"`
	message  Message `json:"message"`
}

func parseTelegramRequest(r *http.Request) (*Update, error) {
	var update Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("Could not decode incoming update %s", err.Error())
		return nil, err
	}
	return &update, nil
}

func telegramWebHook(w http.ResponseWriter, r *http.Request) {
	var update, err = parseTelegramRequest(r)
	if err != nil {
		log.Printf("error parsing update, %s", err.Error())
		return
	}

	var response, errTelegram = sendTextToTelegramBot(update.message.chat.id, update.message.text)
	if errTelegram != nil {
		log.Printf("Error: %s, response: %s", errTelegram.Error(), response)
	} else {
		log.Printf("Echo: %s, chat id: %d", update.message, update.message.chat.id)
	}
}

func sendTextToTelegramBot(chatId int, txt string) (string, error) {
	log.Printf("Send %s to the chat id: %d", txt, chatId)
	var telegramAPI string = telegramApi
	response, err := http.PostForm(
		telegramAPI,
		url.Values{
			"chat_id": {strconv.Itoa(chatId)},
			"text":    {txt},
		})
	if err != nil {
		log.Printf("Error is %s when sending text to the chat bot", err.Error())
		return "", err
	}
	defer response.Body.Close()

	var bodyBytes, errReadBody = ioutil.ReadAll(response.Body)
	if errReadBody != nil {
		log.Printf("Error in parsing telegram answer %s", errReadBody.Error())
		return "", err
	}
	bodyString := string(bodyBytes)
	log.Printf("Body of Telegram Response: %s", bodyString)

	return bodyString, nil
}
