package voiceflow

import (
	"basicScraper/internal/schemas"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type ChatbotClient struct {
	Url string
}

func NewChatbotClient(url string) *ChatbotClient {
	return &ChatbotClient{
		Url: url,
	}
}

func (c *ChatbotClient) AskQuestions(questions []string) ([]schemas.VoiceFlowExcelRow, error) {
	payload, _ := json.Marshal(schemas.VoiceFlowQuestion{
		Action: schemas.VoiceFlowAction{
			Type: schemas.ActionLaunch,
		},
	})

	req, _ := http.NewRequest("POST", c.Url, bytes.NewBuffer(payload))
	_, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	excelRow := make([]schemas.VoiceFlowExcelRow, 0)
	for _, question := range questions {
		answer, err := c.addQuestion(question)
		if err != nil {
			log.Println("Error with the question: ", err)
			excelRow = append(excelRow, schemas.VoiceFlowExcelRow{Question: question, Response: err.Error()})
			continue
		}
		excelRow = append(excelRow, schemas.VoiceFlowExcelRow{Question: question, Response: answer})
	}

	return excelRow, nil
}

func (c *ChatbotClient) addQuestion(q string) (string, error) {
	payload, _ := json.Marshal(schemas.VoiceFlowQuestion{
		Action: schemas.VoiceFlowAction{
			Type:    schemas.ActionText,
			Payload: q,
		},
	})

	req, _ := http.NewRequest("POST", c.Url, bytes.NewBuffer(payload))

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", os.Getenv("VOICEFLOWAPIKEY"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	rsp := []schemas.VoiceFlowResponse{}
	json.Unmarshal(body, &rsp)

	for _, trace := range rsp {
		if trace.Type == "text" {
			return trace.Payload.Message, nil
		}
	}

	return "", fmt.Errorf("response could not be made")
}
