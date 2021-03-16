package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
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
}

type Update struct {
	updateID int     `json:"update_id"`
	message  Message `json:"message"`
}
