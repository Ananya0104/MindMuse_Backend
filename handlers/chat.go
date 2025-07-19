package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"lambda-server/helpers"
	"lambda-server/models"

	"github.com/gin-gonic/gin"
)

// Placeholder for Hugging Face API key
var HuggingFaceAPIKey = "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"

// ChatRequest represents the incoming chat request from frontend
// Includes userId, sessionId, and the user's message
//
type ChatRequest struct {
	UserId    string `json:"userId" binding:"required"`
	SessionId string `json:"sessionId" binding:"required"`
	Message   string `json:"message" binding:"required"`
}

type ChatResponse struct {
	AIResponse string `json:"aiResponse"`
}

// HandleChat handles the chat POST endpoint
func HandleChat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	// Fetch last N messages for context (e.g., last 10)
	const contextLimit = 10
	chatHistory, err := helpers.GetChatHistoryBySession(req.UserId, req.SessionId, contextLimit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chat history", "details": err.Error()})
		return
	}

	// Build prompt for Hugging Face (append all previous messages)
	prompt := ""
	for _, msg := range chatHistory {
		if msg.Sender == "user" {
			prompt += "User: " + msg.Message + "\n"
		} else {
			prompt += "AI: " + msg.Message + "\n"
		}
	}
	prompt += "User: " + req.Message + "\nAI:"

	// Call Hugging Face API
	aiResponse, err := callHuggingFaceAPI(prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get AI response", "details": err.Error()})
		return
	}

	timestamp := time.Now().Unix()
	// Store user message
	userMsg := &models.ChatMessage{
		UserId:    req.UserId,
		SessionId: req.SessionId,
		Timestamp: timestamp,
		Sender:    "user",
		Message:   req.Message,
	}
	helpers.StoreChatMessage(userMsg)

	// Store AI response
	aiMsg := &models.ChatMessage{
		UserId:    req.UserId,
		SessionId: req.SessionId,
		Timestamp: timestamp + 1, // ensure ordering
		Sender:    "ai",
		Message:   aiResponse,
	}
	helpers.StoreChatMessage(aiMsg)

	c.JSON(http.StatusOK, ChatResponse{AIResponse: aiResponse})
}

// callHuggingFaceAPI sends the prompt to Hugging Face and returns the AI's response
func callHuggingFaceAPI(prompt string) (string, error) {
	apiURL := "https://router.huggingface.co/v1/chat/completions" // Inference Providers router endpoint

	// Prepare messages array for chat format
	messages := []map[string]string{
		{"role": "user", "content": prompt},
	}
	body, _ := json.Marshal(map[string]interface{}{
		"messages": messages,
		"model":   "moonshotai/Kimi-K2-Instruct:novita",
	})

	req, err := http.NewRequestWithContext(context.Background(), "POST", apiURL, ioutil.NopCloser(bytes.NewReader(body)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+HuggingFaceAPIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Debug: print raw Hugging Face response
	fmt.Println("Hugging Face raw response:", string(respBody))

	// Parse Inference Providers response
	var hfResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &hfResp); err == nil && len(hfResp.Choices) > 0 {
		return hfResp.Choices[0].Message.Content, nil
	}
	return "", os.ErrInvalid
} 