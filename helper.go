package aoe4api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func queryEloToMap(r *request, sm *safeMap) error {
	payloadBytes, err := json.Marshal(r.pLoad)
	if err != nil {
		return fmt.Errorf("error marshaling json payload: %w", err)
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://api.ageofempires.com/api/ageiv/Leaderboard", body)
	if err != nil {
		return fmt.Errorf("error creating POST request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", r.userAgent)

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending POST to API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNoContent {
			return nil
		}
		return fmt.Errorf("error from API, received status code %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading API response: %w", err)
	}

	var respBodyJson response
	if err := json.Unmarshal(respBody, &respBodyJson); err != nil {
		return fmt.Errorf("error unmarshaling json API response: %w", err)
	}
	if respBodyJson.Count < 1 {
		return nil
	}

	sm.Lock()
	defer sm.Unlock()
	sm.respMap[r.pLoad.MatchType] = strconv.Itoa(respBodyJson.Items[0].Elo)

	return nil
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
