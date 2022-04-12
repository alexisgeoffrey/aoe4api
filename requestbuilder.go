package aoe4api

import (
	"errors"
	"fmt"
	"net/http"
)

type requestBuilder struct{ request }

type Region int

const (
	Europe Region = iota
	MiddleEast
	Asia
	NorthAmerica
	SouthAmerica
	Oceania
	Africa
	Global
)

type Versus int

const (
	Players Versus = iota
	AI
)

type MatchType int

const (
	Unranked MatchType = iota
	Custom
	EasyAI
	MediumAI
	HardAI
	ExpertAI
)

type TeamSize int

const (
	OneVOne TeamSize = iota + 1
	TwoVTwo
	ThreeVThree
	FourVFour
)

func NewRequestBuilder() *requestBuilder {
	return &requestBuilder{
		request{
			client: http.DefaultClient,
			payload: &payload{
				Region:    int(Global),
				Versus:    "players",
				MatchType: "unranked",
				TeamSize:  "1v1",
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
	var vsString string

	switch vs {
	case Players:
		vsString = "players"
	case AI:
		vsString = "ai"
	}

	r.payload.Versus = vsString
	return r
}

func (r *requestBuilder) SetMatchType(mt MatchType) *requestBuilder {
	var mtString string

	switch mt {
	case Unranked:
		mtString = "unranked"
	case Custom:
		mtString = "custom"
	case EasyAI:
		mtString = "aieasy"
	case MediumAI:
		mtString = "aimedium"
	case HardAI:
		mtString = "aihard"
	case ExpertAI:
		mtString = "aiexpert"
	}

	r.payload.MatchType = mtString
	return r
}

func (r *requestBuilder) SetTeamSize(ts TeamSize) *requestBuilder {
	var teamSizeString string

	switch ts {
	case OneVOne:
		teamSizeString = "1v1"
	case TwoVTwo:
		teamSizeString = "2v2"
	case ThreeVThree:
		teamSizeString = "3v3"
	case FourVFour:
		teamSizeString = "4v4"
	}

	r.payload.TeamSize = teamSizeString
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
		return nil, errors.New("cannot have a negative page number")
	}
	if r.payload.Count < 1 {
		return nil, errors.New("cannot have a negative result count")
	}
	if r.payload.Region < int(Europe) || r.payload.Region > int(Global) {
		return nil, errors.New("invalid region")
	}

	switch r.payload.MatchType {
	case "unranked", "custom":
		if r.payload.Versus == "ai" {
			return nil, fmt.Errorf("cannot have both match type as '%s' and versus as 'AI'", r.payload.MatchType)
		}
	case "aieasy", "aimedium", "aihard", "aiexpert":
		if r.payload.Versus == "players" {
			return nil, fmt.Errorf("cannot have both match type as '%s' and versus as 'Players'", r.payload.MatchType)
		}
	}

	payloadCopy := *r.payload

	return &request{
		r.client,
		r.userAgent,
		&payloadCopy,
	}, nil
}
