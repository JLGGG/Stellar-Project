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

const KEYPAIR_COMMAND string = "/keypair"
const MAKE_ACCOUNT_COMMAND string = "/make_account"
const START_COMMAND string = "/start"
const WALLET_COMMAND string = "/wallet"

const TelegramApiUrl string = "https://api.telegram.org/bot"
const TelegramToken string = "BOT_TOKEN"
const TelegramSendMessage string = "/sendMessage"

var globalTelegramApi string = TelegramApiUrl + os.Getenv(TelegramToken) + TelegramSendMessage

type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}

type Chat struct {
	Id int `json:"id"`
}

//Update is a Telegram object that the handler receives every time an user interacts with the bot
type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}

func parseTelegramRequest(r *http.Request) (*Update, error) {
	var update Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("Could not decode incoming update %s", err.Error())
		return nil, err
	}
	return &update, nil
}

func TelegramWebHook(w http.ResponseWriter, r *http.Request) {
	// Parse incoming request
	var update, err = parseTelegramRequest(r)
	if err != nil {
		log.Printf("error parsing update, %s", err.Error())
		return
	}

	// echo for test.
	var response, errTelegram = sendTextToTelegramBot(update.Message.Chat.Id, "test")
	if errTelegram != nil {
		log.Printf("Error: %s, response: %s", errTelegram.Error(), response)
	} else {
		log.Printf("Test: %s, chat id: %d", "test", update.Message.Chat.Id)
	}
}

func sendTextToTelegramBot(chatId int, txt string) (string, error) {
	log.Printf("Send %s to the chat id: %d", txt, chatId)
	var telegramAPI string = globalTelegramApi
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
