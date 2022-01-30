package aoe4api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type (
	Request interface {
		Query() ([]playerInfo, error)
		QueryElo() (string, error)
		QueryAllElo() (map[string]string, error)
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

// Query queries the AOE4 API and returns the API response as a Response struct.
func (r *request) Query() ([]playerInfo, error) {
	result, err := query(r)
	if err != nil {
		return nil, fmt.Errorf("error querying aoe api: %w", err)
	}

	return result.Items, nil
}

// QueryElo queries the AOE4 API and returns the corresponding Elo value as a string.
func (r *request) QueryElo() (string, error) {
	response, err := query(r)
	if err != nil {
		return "", fmt.Errorf("error querying aoe api: %w", err)
	}

	if response.Count >= 0 {
		return strconv.Itoa(response.Items[0].Elo), nil
	}

	return "", fmt.Errorf("no Elo value found for match type %s for username %s", r.payload.MatchType, r.payload.SearchPlayer)
}

// QueryAllElo queries the AOE4 API and returns Elo values for all Elo types
// for a specific username as a map of Elo types and Elo values.
func (r *request) QueryAllElo() (map[string]string, error) {
	var wg sync.WaitGroup
	sm := &safeMap{respMap: make(map[string]string)}

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
				sm.Lock()
				defer sm.Unlock()
				if len(response.Items) > 0 {
					sm.respMap[ts] = strconv.Itoa(response.Items[0].Elo)
				}
			}

			wg.Done()
		}(teamSize)
	}
	wg.Wait()

	return sm.respMap, nil
}
