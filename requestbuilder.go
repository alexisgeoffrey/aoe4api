package aoe4api

import "net/http"

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
		Build() (Request, error)
	}

	requestBuilder struct {
		client    *http.Client
		userAgent string
		pLoad     *payload
	}
)

func NewRequestBuilder() RequestBuilder {
	return &requestBuilder{
		client: http.DefaultClient,
		pLoad: &payload{
			Region:    int(Global),
			Versus:    "players",
			MatchType: "unranked",
			Page:      1,
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
	r.pLoad.Region = int(reg)
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

	r.pLoad.Versus = vsString
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

	r.pLoad.MatchType = mtString
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

	r.pLoad.Versus = teamSizeString
	return r
}

func (r *requestBuilder) SetSearchPlayer(searchPlayer string) RequestBuilder {
	r.pLoad.SearchPlayer = searchPlayer
	return r
}

func (r *requestBuilder) SetPage(page int) RequestBuilder {
	r.pLoad.Page = page
	return r
}

func (r *requestBuilder) SetCount(count int) RequestBuilder {
	r.pLoad.Count = count
	return r
}

func (r *requestBuilder) Build() (Request, error) {
	return &request{}, nil // TODO
}
