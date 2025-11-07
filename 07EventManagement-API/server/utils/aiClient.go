package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/models"
)

// ✅ Generate Event Description (Full context-aware)
func GenerateEventDescription(event models.Event) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=%s", apiKey)

	prompt := fmt.Sprintf(`
You are an intelligent AI writer. Based on the event details below,
write a short, modern, and engaging event description (under 100 words).

Event Name: %s
Type: %s
Expected Attendance: %d
Location: %s
Status: %s
Visibility: %s
Image URL: %s
Created At: %s

Make it sound realistic and professional, suitable for event listings or promotional pages.
`, event.EventName, event.EventtType, event.EventAttendence, event.Location, event.Status, event.IsPublic, event.ImageUrl, event.CreatedAt.Format("02 Jan 2006"))

	reqBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{"parts": []map[string]string{{"text": prompt}}},
		},
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 40 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Gemini raw response (event):", string(body))

	var data map[string]interface{}
	_ = json.Unmarshal(body, &data)

	// Extract clean text
	if candidates, ok := data["candidates"].([]interface{}); ok && len(candidates) > 0 {
		candidate := candidates[0].(map[string]interface{})
		if content, ok := candidate["content"].(map[string]interface{}); ok {
			if parts, ok := content["parts"].([]interface{}); ok && len(parts) > 0 {
				if part, ok := parts[0].(map[string]interface{}); ok {
					if text, ok := part["text"].(string); ok {
						return strings.TrimSpace(text), nil
					}
				}
			}
		}
	}
	return "AI failed to generate event description 😔", nil
}

// ✅ Generate Function Description (Full context-aware)
func GenerateFunctionDescription(function models.Function) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=%s", apiKey)

	prompt := fmt.Sprintf(`
You are an AI assistant that writes elegant function descriptions.
Using the following details, write a lively and beautiful short description (under 80 words).

Function Name: %s
Type: %s
Location: %s
Status: %s
Visibility: %s
Image URL: %s
Created At: %s

Make it sound natural, inviting, and event-like — as if written by a professional event planner.
`, function.FuncName, function.FuncType, function.Location, function.Status, function.IsPublic, function.ImageUrl, function.CreatedAt.Format("02 Jan 2006"))

	reqBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{"parts": []map[string]string{{"text": prompt}}},
		},
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 40 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Gemini raw response (function):", string(body))

	var data map[string]interface{}
	_ = json.Unmarshal(body, &data)

	// Extract clean text
	if candidates, ok := data["candidates"].([]interface{}); ok && len(candidates) > 0 {
		candidate := candidates[0].(map[string]interface{})
		if content, ok := candidate["content"].(map[string]interface{}); ok {
			if parts, ok := content["parts"].([]interface{}); ok && len(parts) > 0 {
				if part, ok := parts[0].(map[string]interface{}); ok {
					if text, ok := part["text"].(string); ok {
						return strings.TrimSpace(text), nil
					}
				}
			}
		}
	}
	return "AI failed to generate function description 😔", nil
}

// 🧠 Generic Gemini AI Response Generator (for recommend/assistant)
func GenerateAIResponse(prompt string) (map[string]interface{}, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=%s", apiKey)

	reqBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{"parts": []map[string]string{{"text": prompt}}},
		},
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 40 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Gemini raw response (assistant):", string(body))

	var data map[string]interface{}
	_ = json.Unmarshal(body, &data)

	if candidates, ok := data["candidates"].([]interface{}); ok && len(candidates) > 0 {
		candidate := candidates[0].(map[string]interface{})
		if content, ok := candidate["content"].(map[string]interface{}); ok {
			if parts, ok := content["parts"].([]interface{}); ok && len(parts) > 0 {
				if part, ok := parts[0].(map[string]interface{}); ok {
					if text, ok := part["text"].(string); ok {
						cleanText := strings.TrimSpace(text)
						cleanText = strings.TrimPrefix(cleanText, "```json")
						cleanText = strings.TrimSuffix(cleanText, "```")
						cleanText = strings.TrimSpace(cleanText)

						var parsed map[string]interface{}
						if err := json.Unmarshal([]byte(cleanText), &parsed); err == nil {
							return parsed, nil
						}
						return map[string]interface{}{"text": cleanText}, nil
					}
				}
			}
		}
	}

	return map[string]interface{}{"error": "AI failed to generate structured output 😔"}, nil
}
