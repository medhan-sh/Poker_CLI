package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	store := NewInMemoryPlayerStore()
	server := PlayerServer{store}
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), PostWinReqest(player))
	server.ServeHTTP(httptest.NewRecorder(), PostWinReqest(player))
	server.ServeHTTP(httptest.NewRecorder(), PostWinReqest(player))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetScoreRequest(player))
	asserStatusCode(t, response.Code, http.StatusOK)

	assertResponse(t, response.Body.String(), "3")

}
