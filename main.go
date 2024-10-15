package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Event struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
}

func fetchUserActivity(username string) ([]Event, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/events", username)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing response body:", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: failed to fetch data, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var events []Event
	err = json.Unmarshal(body, &events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	log.Info().Msg("This is an info message. You are using Zerolog!")
	os.Stderr.Sync()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-sigChan
		log.Info().Msg("Received interrupt signal, closing program...")
		cancel()
	}()

	var username string

	if len(os.Args) < 2 {
		log.Info().Msg("Please provide a username")
		log.Info().Msg("Please enter the GitHub username. Example: github-activity <username>")

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("github-activity ")
		input, _ := reader.ReadString('\n')

		username = strings.TrimSpace(input)
	} else {
		username = os.Args[1]
	}

	log.Info().Msgf("Fetching activity for user: %s", username)

	select {
	case <-ctx.Done():
		log.Info().Msg("Program Closed!")
		return
	default:
		events, err := fetchUserActivity(username)
		if err != nil {
			log.Error().Err(err).Msg("Failed to fetch events")
			return
		}

		for _, event := range events {
			fmt.Printf("%s event at %s\n", event.Type, event.Repo.Name)
		}
	}

	log.Info().Msg("Program Closed!")
}
