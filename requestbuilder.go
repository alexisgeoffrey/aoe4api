package aoe4api

import (
	"errors"
	"fmt"
	"net/http"
)

type (
	RequestBuilder interface {
		SetHttpClient(*http.Client) RequestBuilder
		SetUserAgent(string) RequestBuilder
		SetRegion(Region) RequestBuilder
		SetVersus(Versus) RequestBuilder
		SetMatchType(MatchType) RequestBuilder
		SetTeamSize(TeamSize) RequestBuilder
		SetSearchPlayer(string) RequestBuilder
		SetPage(int) RequestBuilder
		SetCount(int) RequestBuilder
		Request() (Request, error)
	}

	requestBuilder struct {
		client    *http.Client
		userAgent string
		payload   *payload
	}
)

func NewRequestBuilder() RequestBuilder {
	return &requestBuilder{
		client: http.DefaultClient,
		payload: &payload{
			Region:    int(Global),
			Versus:    "players",
			MatchType: "unranked",
			TeamSize:  "1v1",
			Page:      1,
			Count:     100,
		},
	}
}

func (r *requestBuilder) SetHttpClient(client *http.Client) RequestBuilder {
	r.client = client
	return r
}

func (r *requestBuilder) SetUserAgent(userAgent string) RequestBuilder {
	r.userAgent = userAgent
	return r
}

func (r *requestBuilder) SetRegion(reg Region) RequestBuilder {
	r.payload.Region = int(reg)
	return r
}

func (r *requestBuilder) SetVersus(vs Versus) RequestBuilder {
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

func (r *requestBuilder) SetMatchType(mt MatchType) RequestBuilder {
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

func (r *requestBuilder) SetTeamSize(ts TeamSize) RequestBuilder {
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

	r.payload.Versus = teamSizeString
	return r
}

func (r *requestBuilder) SetSearchPlayer(searchPlayer string) RequestBuilder {
	r.payload.SearchPlayer = searchPlayer
	return r
}

func (r *requestBuilder) SetPage(page int) RequestBuilder {
	r.payload.Page = page
	return r
}

func (r *requestBuilder) SetCount(count int) RequestBuilder {
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

	return &request{
		r.client,
		r.userAgent,
		r.payload,
	}, nil
}
