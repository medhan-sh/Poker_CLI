package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type stubPlayerStore struct {
	score    map[string]int
	winCalls []string
	league   league
}

func (s *stubPlayerStore) GetPlayerScore(name string) int {
	score := s.score[name]
	return score
}
func (s *stubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}
func (s *stubPlayerStore) GetLeague() league {
	return s.league
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
		nil, nil,
	}
	server := NewPlayerServer(&store)
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
		nil, nil,
	}
	server := NewPlayerServer(&store)

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

// league
func NewLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
}
func assertLeague(t testing.TB, got, want []Player) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
func getLeagueFromResponse(t testing.TB, body io.Reader) (league []Player) {
	t.Helper()
	league, err := NewLeague(body)

	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", body, err)
	}

	return
}
func assertContentType(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}
func TestLeague(t *testing.T) {

	t.Run("return 200 on /league", func(t *testing.T) {
		store := stubPlayerStore{}
		server := NewPlayerServer(&store)
		request := NewLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		var got []Player

		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("Unable to pasrse response from server %v to slice Player: %v", response.Body, err)
		}
		asserStatusCode(t, response.Code, http.StatusOK)
	})
	t.Run("Return league tabl as JSON", func(t *testing.T) {
		wantedLeague := []Player{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}
		store := stubPlayerStore{nil, nil, wantedLeague}
		server := NewPlayerServer(&store)

		request := NewLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getLeagueFromResponse(t, response.Body)
		assertLeague(t, got, wantedLeague)
		asserStatusCode(t, response.Code, http.StatusOK)
		assertContentType(t, response, jsonContentType)
	})
}
