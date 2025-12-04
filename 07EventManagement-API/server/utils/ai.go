package utils

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func AI(prompt string) (string, error) {

	// Escape user input to avoid JSON break
	safe := strings.ReplaceAll(prompt, `"`, `'`)
	safe = strings.ReplaceAll(safe, "\n", "\\n")

	apiKey := os.Getenv("GEMINI_API_KEY")
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=" + apiKey

	body := `{
		"contents": [
			{
				"parts": [
					{"text": "` + safe + `" }
				]
			}
		]
	}`

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	raw, _ := ioutil.ReadAll(resp.Body)
	respStr := string(raw)

	// Extract clean text
	marker := `"text": "`
	start := strings.Index(respStr, marker)

	if start != -1 {
		start += len(marker)
		end := strings.Index(respStr[start:], `"`)
		if end != -1 {
			return strings.TrimSpace(respStr[start : start+end]), nil
		}
	}

	return respStr, nil
}
