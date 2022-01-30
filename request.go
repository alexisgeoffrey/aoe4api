// Package aoe4api provides functions for retrieving
// player data from the Age of Empires 4 leaderboard API.
package aoe4api

import (
	"fmt"
	"net/http"
	"sync"
)

type (
	Request interface {
		Query() (PlayerInfo, error)
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
		Count int              `json:"count"`
		Items []playerInfoData `json:"items"`
	}

	playerInfoData struct {
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

	PlayerInfo struct {
		GameID       string
		UserID       string
		RlUserID     int
		UserName     string
		AvatarURL    string
		PlayerNumber interface{} // For now, will always be nil
		Elo          int
		EloRating    int
		Rank         int
		Region       string
		Wins         int
		WinPercent   float64
		Losses       int
		WinStreak    int
		RankLevel    string
		RankIcon     string
	}

	safeMap struct {
		sync.Mutex
		respMap map[string]string
	}
)

// Query queries the AOE4 API and returns the API response as a Response struct.
func (r *request) Query() (PlayerInfo, error) {
	return PlayerInfo{}, nil // TODO
}

// QueryElo queries the AOE4 API and returns the corresponding Elo value as a string.
func (r *request) QueryElo() (string, error) {
	sm := &safeMap{respMap: make(map[string]string)}

	if err := queryEloToMap(r, sm); err != nil {
		return "", fmt.Errorf("error querying aoe api or inserting in map: %w", err)
	}

	if elo, ok := sm.respMap[r.payload.MatchType]; ok {
		return elo, nil
	}
	return "", fmt.Errorf("no Elo value found for match type %s for username %s", r.payload.MatchType, r.payload.SearchPlayer)
}

// QueryAllElo queries the AOE4 API and returns Elo values for all Elo types
// for a specific username as a map of Elo types and Elo values.
func (r *request) QueryAllElo() (map[string]string, error) {
	var wg sync.WaitGroup
	sm := &safeMap{respMap: make(map[string]string)}

	for _, teamSize := range getEloTypes() {
		req := *r
		if teamSize == "custom" {
			req.payload.MatchType = teamSize
			req.payload.TeamSize = ""
		} else {
			req.payload.MatchType = "unranked"
			req.payload.TeamSize = teamSize
		}

		wg.Add(1)
		go func() {
			if err := queryEloToMap(&req, sm); err != nil {
				fmt.Printf("error retrieving Elo from AOE api for %s: %v", req.payload.SearchPlayer, err)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	return sm.respMap, nil
}
