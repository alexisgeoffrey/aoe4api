package aoe4api

import (
	"errors"
	"fmt"
	"net/http"
)

type requestBuilder struct {
	client       *http.Client
	userAgent    string
	searchPlayer string
	versus       Versus
	matchType    MatchType
	teamSize     TeamSize
	region       Region
	page         int
	count        int
}

func NewRequestBuilder() *requestBuilder {
	return &requestBuilder{
		client:    http.DefaultClient,
		versus:    Players,
		matchType: Unranked,
		teamSize:  OneVOne,
		region:    Global,
		page:      1,
		count:     100,
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
	r.region = reg
	return r
}

func (r *requestBuilder) SetVersus(vs Versus) *requestBuilder {
	r.versus = vs
	return r
}

func (r *requestBuilder) SetMatchType(mt MatchType) *requestBuilder {
	r.matchType = mt
	return r
}

func (r *requestBuilder) SetTeamSize(ts TeamSize) *requestBuilder {
	r.teamSize = ts
	return r
}

func (r *requestBuilder) SetSearchPlayer(searchPlayer string) *requestBuilder {
	r.searchPlayer = searchPlayer
	return r
}

func (r *requestBuilder) SetPage(page int) *requestBuilder {
	r.page = page
	return r
}

func (r *requestBuilder) SetCount(count int) *requestBuilder {
	r.count = count
	return r
}

func (r *requestBuilder) Request() (Request, error) {
	if r.page < 1 {
		return nil, errors.New("cannot have a zero or negative page number")
	}
	if r.count < 1 {
		return nil, errors.New("cannot have a zero or negative result count")
	}
	if r.region < Europe || r.region > Global {
		return nil, errors.New("invalid region")
	}

	switch r.matchType {
	case Unranked, Custom:
		if r.versus == AI {
			return nil, fmt.Errorf("cannot have both match type as '%s' and versus as 'ai'", r.matchType)
		}
	case EasyAI, MediumAI, HardAI, ExpertAI:
		if r.versus == Players {
			return nil, fmt.Errorf("cannot have both match type as '%s' and versus as 'players'", r.matchType)
		}
	}

	return &request{
		r.client,
		r.userAgent,
		payload{
			string(r.versus),
			string(r.matchType),
			string(r.teamSize),
			r.searchPlayer,
			int(r.region),
			r.page,
			r.count,
		},
	}, nil
}
