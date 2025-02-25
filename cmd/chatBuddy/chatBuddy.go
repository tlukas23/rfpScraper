package main

import (
	"basicScraper/internal/voiceflow"
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	err := loadEnvFile(".env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	filename := flag.String("file", "ChatbotQs,xlsx", "Name of the file (.xlsx)")
	flag.Parse()

	log.Println(*filename)

	questions, err := voiceflow.ExtractQuestions(*filename)
	if err != nil {
		log.Fatal("Unable to read file: ", err)
	}

	chatbot := voiceflow.NewChatbotClient("https://general-runtime.voiceflow.com/state/user/1/interact?logs=off")
	rows, err := chatbot.AskQuestions(questions)
	if err != nil {
		log.Fatal("Error while chattings: ", err)
	}

	if err := voiceflow.MakeResultExcelDoc(rows); err != nil {
		log.Fatal("Unable to make result sheet: ", err)
	}
}

func loadEnvFile(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		os.Setenv(key, value)
	}

	return scanner.Err()
}
