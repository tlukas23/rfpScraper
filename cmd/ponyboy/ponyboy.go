package main

import (
	"basicScraper/internal/rfp"
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

func main() {

	err := loadEnvFile(".env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}
	token := os.Getenv("BOT_TOKEN")
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running. Press Ctrl+C to exit.")
	// Wait here until CTRL+C is entered
	select {
	case <-waitExit():
		return
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Check if the message is from a bot
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!RFP") {

		s.ChannelMessageSend(m.ChannelID, "Scraping rfp.com ..")
		// Replace "path/to/your/file.txt" with the path to your file

		filename, err := rfp.MakeRfpSpreadsheet()
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Unable to scrape rfp.com right now")
			return
		}
		log.Println("Sending file to discord server")

		file, err := os.Open(filename)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		currentTime := time.Now()

		// Format the date as a string
		currentDateString := currentTime.Format("2006-01-02") // YYYY-MM-DD
		fileName := "RFPLeads-" + currentDateString + ".xlsx"
		// Send the file to the channel
		s.ChannelFileSend(m.ChannelID, fileName, file)
	}
}

func waitExit() <-chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
	return c
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
