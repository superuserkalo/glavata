package main

import (
	"net/http"
	"net/url"
	"encoding/json"
	"os"
	"log"
	"strconv"

	"github.com/joho/godotenv"
)

type From struct {
	FirstName string `json:"first_name"`
}

type Chat struct {
	ID int64 `json:"id"`
}

type Message struct {
	MessageID int `json:"message_id"`
	From *From `json:"from"`
	Chat *Chat `json:"chat"`
	Text string `json:"text"`
}

type Update struct {
	UpdateID int `json:"update_id"`
	Message *Message `json:"message"`
}

type UpdateResp struct {
	Ok bool `json:"ok"`
	Res []Update `json:"result"`
}

func main() {

	err := godotenv.Load(".env")
  if err != nil {
    log.Fatal("Error loading .env file")
  }

	baseUrl := "https://api.telegram.org/bot"
	token := os.Getenv("telegram_key")
	apiBase := baseUrl + token + "/"

	offset := 0
	getUpdate(apiBase, offset)
}

func sendMessage(apiBase string, chatID int64, text string) error {
	values := url.Values{}
	values.Set("chat_id", strconv.FormatInt(chatID, 10))
	values.Set("text", text)

	resp, err := http.PostForm(apiBase+"sendMessage", values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func getUpdate(apiBase string, offset int) {
	for {
		resp, err := http.Get(apiBase + "getUpdates?timeout=30&offset=" + strconv.Itoa(offset))
		if err != nil {
			log.Fatal("getUpdates error:", err)
		}

		var upRes UpdateResp
		if err := json.NewDecoder(resp.Body).Decode(&upRes); err != nil {
			resp.Body.Close()
			log.Println("decode error:", err)
			continue
  		}
		resp.Body.Close()

		if !upRes.Ok {
			log.Println("telegram returned ok=false")
			continue
		}

		if len(upRes.Res) == 0 {
			continue
		}

		for _, u := range(upRes.Res) {
			offset = u.UpdateID + 1

			if u.Message == nil || u.Message.Text == "" {
				continue
			}

			name := u.Message.From.FirstName
			if name == "" {
				name = u.Message.Text
			}

			if err := sendMessage(apiBase, u.Message.Chat.ID, "Hey "+name); err != nil {
				log.Println("sendMessage error:", err)
			}
		}
	}
}
