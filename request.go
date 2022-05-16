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
		QueryElo(userId string) (int, error)
		QueryAllElo(userId string) (map[string]string, error)
	}

	request struct {
		client    *http.Client
		userAgent string
		payload   payload
	}

	// A Payload represents the contents of a request to the AOE4 API.
	payload struct {
		Versus       string `json:"versus"`
		MatchType    string `json:"matchType"`
		TeamSize     string `json:"teamSize"`
		SearchPlayer string `json:"searchPlayer"`
		Region       int    `json:"region"`
		Page         int    `json:"page"`
		Count        int    `json:"count"`
	}

	response struct {
		Items []playerInfo `json:"items"`
		Count int          `json:"count"`
	}

	playerInfo struct {
		PlayerNumber interface{} `json:"playerNumber"` // For now, will always be nil
		GameID       string      `json:"gameId"`
		UserID       string      `json:"userId"`
		UserName     string      `json:"userName"`
		AvatarURL    string      `json:"avatarUrl"`
		Region       string      `json:"region"`
		RankLevel    string      `json:"rankLevel"`
		RankIcon     string      `json:"rankIcon"`
		RlUserID     int         `json:"rlUserId"`
		Elo          int         `json:"elo"`
		EloRating    int         `json:"eloRating"`
		Rank         int         `json:"rank"`
		Wins         int         `json:"wins"`
		WinPercent   float64     `json:"winPercent"`
		Losses       int         `json:"losses"`
		WinStreak    int         `json:"winStreak"`
	}

	safeMap struct {
		respMap map[string]string
		sync.Mutex
	}
)

const apiUrl = "https://api.ageofempires.com/api/ageiv/Leaderboard"

var eloTypes = [...]enum{
	OneVOne,
	TwoVTwo,
	ThreeVThree,
	FourVFour,
	Custom,
}

// Query queries the AOE4 API and returns the API response as a slice of playerInfo structs.
func (r *request) Query() ([]playerInfo, error) {
	result, err := r.query()
	if err != nil {
		return nil, fmt.Errorf("error querying aoe api: %w", err)
	}

	return result.Items, nil
}

// QueryElo queries the AOE4 API and returns the corresponding Elo value as a string.
func (r *request) QueryElo(userId string) (int, error) {
	response, err := r.query()
	if err != nil {
		return 0, fmt.Errorf("error querying aoe api: %w", err)
	}

	if response.Count > 0 {
		for _, item := range response.Items {
			if strings.Contains(item.UserID, userId) {
				return item.Elo, nil
			}
		}
	}

	return 0, fmt.Errorf("no Elo value found for match type %s for username %s", r.payload.MatchType, r.payload.SearchPlayer)
}

// QueryAllElo queries the AOE4 API and returns Elo values for all Elo types
// for a specific username as a map of Elo types and Elo values.
func (r *request) QueryAllElo(userId string) (map[string]string, error) {
	var wg sync.WaitGroup
	sm := &safeMap{respMap: map[string]string{}}

	for _, et := range eloTypes {
		req := request{
			r.client,
			r.userAgent,
			r.payload,
		}

		if et == Custom {
			req.payload.MatchType = et.String()
			req.payload.TeamSize = ""
		} else {
			req.payload.MatchType = Unranked.String()
			req.payload.TeamSize = et.String()
		}

		wg.Add(1)
		go func(ts string) {
			response, err := req.query()
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
		}(et.String())
	}
	wg.Wait()

	return sm.respMap, nil
}

func (r *request) query() (response, error) {
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
