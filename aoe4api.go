// Package aoe4api provides functions for retrieving
// player data from the Age of Empires 4 leaderboard API.
package aoe4api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const apiUrl = "https://api.ageofempires.com/api/ageiv/Leaderboard"

func query(r *request) (response, error) {
	payloadBytes, err := json.Marshal(r.payload)
	if err != nil {
		return response{}, fmt.Errorf("error marshaling json payload: %w", err)
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", apiUrl, body)
	if err != nil {
		return response{}, fmt.Errorf("error creating POST request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if r.userAgent != "" {
		req.Header.Set("User-Agent", r.userAgent)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return response{}, fmt.Errorf("error sending POST to API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNoContent {
			return response{}, nil
		}
		return response{}, fmt.Errorf("error from API, received status code %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return response{}, fmt.Errorf("error reading API response: %w", err)
	}

	var respBodyJson response
	if err := json.Unmarshal(respBody, &respBodyJson); err != nil {
		return response{}, fmt.Errorf("error unmarshaling json API response: %w", err)
	}

	return respBodyJson, nil
}

func getEloTypes() [5]string {
	return [...]string{
		"1v1",
		"2v2",
		"3v3",
		"4v4",
		"custom",
	}
}
