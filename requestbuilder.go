package aoe4api

import (
	"errors"
	"fmt"
	"net/http"
)

type requestBuilder struct{ request }

func NewRequestBuilder() *requestBuilder {
	return &requestBuilder{
		request{
			client: http.DefaultClient,
			payload: payload{
				Versus:    "players",
				MatchType: "unranked",
				TeamSize:  "1v1",
				Region:    int(Global),
				Page:      1,
				Count:     100,
			},
		},
	}
}

func (r *requestBuilder) SetHttpClient(client *http.Client) *requestBuilder {
	r.client = client
	return r
}

func (r *requestBuilder) SetUserAgent(userAgent string) *requestBuilder {
	r.userAgent = userAgent
	return r
}

func (r *requestBuilder) SetRegion(reg Region) *requestBuilder {
	r.payload.Region = int(reg)
	return r
}

func (r *requestBuilder) SetVersus(vs Versus) *requestBuilder {
	r.payload.Versus = vs.String()
	return r
}

func (r *requestBuilder) SetMatchType(mt MatchType) *requestBuilder {
	r.payload.MatchType = mt.String()
	return r
}

func (r *requestBuilder) SetTeamSize(ts TeamSize) *requestBuilder {
	r.payload.TeamSize = ts.String()
	return r
}

func (r *requestBuilder) SetSearchPlayer(searchPlayer string) *requestBuilder {
	r.payload.SearchPlayer = searchPlayer
	return r
}

func (r *requestBuilder) SetPage(page int) *requestBuilder {
	r.payload.Page = page
	return r
}

func (r *requestBuilder) SetCount(count int) *requestBuilder {
	r.payload.Count = count
	return r
}

func (r *requestBuilder) Request() (Request, error) {
	if r.payload.Page < 1 {
		return nil, errors.New("cannot have a zero or negative page number")
	}
	if r.payload.Count < 1 {
		return nil, errors.New("cannot have a zero or negative result count")
	}
	if r.payload.Region < int(Europe) || r.payload.Region > int(Global) {
		return nil, errors.New("invalid region")
	}

	switch r.payload.MatchType {
	case "unranked", "custom":
		if r.payload.Versus == "ai" {
			return nil, fmt.Errorf("cannot have both match type as '%s' and versus as 'ai'", r.payload.MatchType)
		}
	case "aieasy", "aimedium", "aihard", "aiexpert":
		if r.payload.Versus == "players" {
			return nil, fmt.Errorf("cannot have both match type as '%s' and versus as 'players'", r.payload.MatchType)
		}
	}

	return &request{
		r.client,
		r.userAgent,
		r.payload,
	}, nil
}
