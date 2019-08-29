package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type CardServerConfig struct {
	port        int
	gameConfigs []CardGameConfig
}

type CardGameConfig struct {
	Name string
}

func startCardServer(config CardServerConfig) {

	shutDownChannel := make(chan string)
	defer close(shutDownChannel)

	m, err := cardRouterWithShutdown(config, shutDownChannel)
	if err != nil {

	}

	s := http.Server{Addr: ":" + strconv.Itoa(config.port), Handler: m}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	select {
	case code := <-shutDownChannel:
		log.Printf("Got shutdown request: %s", code)
		s.Shutdown(context.Background())

	}
	log.Printf("Finished")
}

func cardRouter(config CardServerConfig) *http.ServeMux {
	m := http.NewServeMux()

	m.Handle("/available-games/", getAvailableGamesHandler(config))

	return m
}

func cardRouterWithShutdown(config CardServerConfig, shutdownChannel chan string) (*http.ServeMux, error) {

	m := cardRouter(config)

	m.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
		shutdownChannel <- "shutdown please"
	})

	return m, nil
}

func getAvailableGamesHandler(config CardServerConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var availableGames []string
		for _, config := range config.gameConfigs {
			availableGames = append(availableGames, config.Name)
		}
		jsonGames, _ := json.Marshal(availableGames)
		w.Write(jsonGames)
	})
}
