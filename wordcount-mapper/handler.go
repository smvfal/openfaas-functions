package function

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	handler "github.com/openfaas/templates-sdk/go-http"
)

// Handle a function invocation
func Handle(req handler.Request) (handler.Response, error) {

	type Message struct {
		Result map[string]int `json:"result"`
	}

	var errorResponse = handler.Response{
		Body:       nil,
		StatusCode: http.StatusInternalServerError,
	}

	chunk := string(req.Body)
	log.Printf("input received: %s", chunk)

	result := make(map[string]int)
	scanner := bufio.NewScanner(strings.NewReader(chunk))
	// Set the split function for the scanning operation.
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		word := scanner.Text()
		word = strings.Trim(word, ".,:;()[]{}!?'\"\"")
		wordLength := len(word)
		if wordLength != 0 {
			result[word] += 1
		}
	}
	if err := scanner.Err(); err != nil {
		return errorResponse, err
	}

	message := Message{result}
	jsonMessage, err := json.MarshalIndent(message, "", "  ")
	if err != nil {
		return errorResponse, err
	}

	return handler.Response{
		Body:       jsonMessage,
		StatusCode: http.StatusOK,
	}, err
}
