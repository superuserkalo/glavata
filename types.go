package main

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

