package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type stubPlayerStore struct {
	score    map[string]int
	winCalls []string
}

func (s *stubPlayerStore) GetPlayerScore(name string) int {
	score := s.score[name]
	return score
}
func (s *stubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func newGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}
func assertResponse(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
func asserStatusCode(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("Got Error code %v, wanted %v", got, want)
	}
}

func TestGetPlayer(t *testing.T) {
	store := stubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		nil,
	}
	server := &PlayerServer{&store}
	t.Run("Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "20"
		asserStatusCode(t, response.Code, http.StatusOK)
		assertResponse(t, got, want)
	})
	t.Run("Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "10"
		asserStatusCode(t, response.Code, http.StatusOK)
		assertResponse(t, got, want)
	})
	t.Run("Unknown input", func(t *testing.T) {
		request := newGetScoreRequest("Dembelle")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Code
		want := http.StatusNotFound
		asserStatusCode(t, got, want)
	})
}

func PostWinReqest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func TestStoreWins(t *testing.T) {
	store := stubPlayerStore{
		map[string]int{},
		nil,
	}
	server := &PlayerServer{&store}

	t.Run("returns accepted on POST", func(t *testing.T) {
		player := "Pepper"
		request := PostWinReqest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		asserStatusCode(t, response.Code, http.StatusAccepted)

		if len(store.winCalls) != 1 {
			t.Errorf("Got %d wincalls, wanted %d", len(store.winCalls), 1)
		}
		if store.winCalls[0] != player {
			t.Errorf("did not store correct winner got %s wanted %s", store.winCalls[0], player)
		}
	})
}
