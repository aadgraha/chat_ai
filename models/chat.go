package models

type ChatRequest struct {
	Prompt string `json:"prompt"`
	ChatID string `json:"chat_id"`
}

type ChatResponse struct {
	Message string `json:"message"`
}