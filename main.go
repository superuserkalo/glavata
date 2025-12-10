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
