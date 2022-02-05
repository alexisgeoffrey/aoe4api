// Package aoe4api provides functions for retrieving
// player data from the Age of Empires 4 leaderboard API.
package aoe4api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type (
	Request interface {
		Query() ([]playerInfo, error)
		QueryElo(userId string) (string, error)
		QueryAllElo(userId string) (map[string]string, error)
	}

	request struct {
		client    *http.Client
		userAgent string
		payload   *payload
	}

	// A Payload represents the contents of a request to the AOE4 API.
	payload struct {
		Region       int    `json:"region"`
		Versus       string `json:"versus"`
		MatchType    string `json:"matchType"`
		TeamSize     string `json:"teamSize"`
		SearchPlayer string `json:"searchPlayer"`
		Page         int    `json:"page"`
		Count        int    `json:"count"`
	}

	response struct {
		Count int          `json:"count"`
		Items []playerInfo `json:"items"`
	}

	playerInfo struct {
		GameID       string      `json:"gameId"`
		UserID       string      `json:"userId"`
		RlUserID     int         `json:"rlUserId"`
		UserName     string      `json:"userName"`
		AvatarURL    string      `json:"avatarUrl"`
		PlayerNumber interface{} `json:"playerNumber"` // For now, will always be nil
		Elo          int         `json:"elo"`
		EloRating    int         `json:"eloRating"`
		Rank         int         `json:"rank"`
		Region       string      `json:"region"`
		Wins         int         `json:"wins"`
		WinPercent   float64     `json:"winPercent"`
		Losses       int         `json:"losses"`
		WinStreak    int         `json:"winStreak"`
		RankLevel    string      `json:"rankLevel"`
		RankIcon     string      `json:"rankIcon"`
	}

	safeMap struct {
		sync.Mutex
		respMap map[string]string
	}
)

const apiUrl = "https://api.ageofempires.com/api/ageiv/Leaderboard"

// Query queries the AOE4 API and returns the API response as a slice of playerInfo structs.
func (r *request) Query() ([]playerInfo, error) {
	result, err := query(r)
	if err != nil {
		return nil, fmt.Errorf("error querying aoe api: %w", err)
	}

	return result.Items, nil
}

// QueryElo queries the AOE4 API and returns the corresponding Elo value as a string.
func (r *request) QueryElo(userId string) (string, error) {
	response, err := query(r)
	if err != nil {
		return "", fmt.Errorf("error querying aoe api: %w", err)
	}

	if response.Count >= 0 {
		for _, item := range response.Items {
			if strings.Contains(item.UserID, userId) {
				return strconv.Itoa(item.Elo), nil
			}
		}
	}

	return "", fmt.Errorf("no Elo value found for match type %s for username %s", r.payload.MatchType, r.payload.SearchPlayer)
}

// QueryAllElo queries the AOE4 API and returns Elo values for all Elo types
// for a specific username as a map of Elo types and Elo values.
func (r *request) QueryAllElo(userId string) (map[string]string, error) {
	var wg sync.WaitGroup
	sm := &safeMap{respMap: map[string]string{}}

	for _, teamSize := range getEloTypes() {
		payloadCopy := *r.payload
		req := request{
			r.client,
			r.userAgent,
			&payloadCopy,
		}

		if teamSize == "custom" {
			req.payload.MatchType = teamSize
			req.payload.TeamSize = ""
		} else {
			req.payload.MatchType = "unranked"
			req.payload.TeamSize = teamSize
		}

		wg.Add(1)
		go func(ts string) {
			response, err := query(&req)
			if err != nil {
				log.Printf("error retrieving Elo from AOE api for %s: %v", req.payload.SearchPlayer, err)
			} else {
				if len(response.Items) > 0 {
					for _, item := range response.Items {
						if strings.Contains(item.UserID, userId) {
							sm.Lock()
							defer sm.Unlock()
							sm.respMap[ts] = strconv.Itoa(item.Elo)
							break
						}
					}
				}
			}

			wg.Done()
		}(teamSize)
	}
	wg.Wait()

	return sm.respMap, nil
}

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
