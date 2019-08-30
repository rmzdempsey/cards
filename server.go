package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)

type CardServerConfig struct {
	port        int
	gameConfigs []CardGameConfig
}

type CardGameConfig struct {
	Name string
}

var serverConfig CardServerConfig

func startCardServer(config CardServerConfig) {

	shutDownChannel := make(chan string)
	defer close(shutDownChannel)

	m := cardRouterWithShutdown(config, shutDownChannel)

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

func cardRouter(config CardServerConfig) *mux.Router {

	serverConfig = config
	
	r := mux.NewRouter()

	r.HandleFunc("/available-games/", getAvailableGamesHandler).Methods("GET")
	r.HandleFunc("/tables/", createTableHandler ).Methods("POST")

	return r
}

func cardRouterWithShutdown(config CardServerConfig, shutdownChannel chan string) *mux.Router {

	m := cardRouter(config)

	m.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
		shutdownChannel <- "shutdown please"
	})

	return m
}

func getAvailableGamesHandler(w http.ResponseWriter, r *http.Request) {
	var availableGames []string
	for _, config := range serverConfig.gameConfigs {
		availableGames = append(availableGames, config.Name)
	}
	jsonGames, _ := json.Marshal(availableGames)
	w.Write(jsonGames)
}
