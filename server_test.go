package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type stubPlayerStore struct {
	score map[string]string
}

func (s *stubPlayerStore) GetPlayerScore(name string) string {
	score := s.score[name]
	return score
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

func TestGetPlayer(t *testing.T) {
	store := stubPlayerStore{
		map[string]string{
			"Pepper": "20",
			"Floyd":  "10",
		},
	}
	server := &PlayerServer{&store}
	t.Run("Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "20"

		assertResponse(t, got, want)
	})
	t.Run("Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "10"
		assertResponse(t, got, want)
	})
}
