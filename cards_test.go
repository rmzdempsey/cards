package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCardsAvailableGamesHandler(t *testing.T) {

	var testCases = []struct {
		name      string
		gameNames []string
	}{
		{"empty", []string{}},
		{"single game", []string{"poker"}},
		{"two games", []string{"poker", "blackjack"}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			config := CardServerConfig{
				port: 8000,
			}
			for _, gameName := range testCase.gameNames {
				config.gameConfigs = append(config.gameConfigs, CardGameConfig{Name: gameName})
			}
			ts := httptest.NewServer(cardRouter(config))
			defer ts.Close()

			res, err := http.Get(ts.URL + "/available-games/" )
			if err != nil {
				t.Fatal(err)
			}
			body, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				t.Fatal(err)
			}
			var resultArr []string
			err = json.Unmarshal(body, &resultArr)
			if err != nil {
				t.Fatal(err)
			}

			if len(resultArr) != len(testCase.gameNames) {
				t.Fatalf("expected %d elements, got %d", len(testCase.gameNames), len(resultArr))
			}
			for idx := 0; idx < len(testCase.gameNames); idx++ {
				if !contains(resultArr, testCase.gameNames[idx]) {
					t.Fatalf("expected to contain %s in %v", testCase.gameNames[idx], resultArr)
				}
			}
		})

	}
}

func contains(sl []string, s string) bool {
	for _, a := range sl {
		if a == s {
			return true
		}
	}
	return false
}
