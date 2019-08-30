package main

import (
	"net/http"
)

type Table struct {
	gameConfig CardGameConfig
	players Player
}

func createTableHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}