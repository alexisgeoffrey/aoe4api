package aoe4api

type enum interface {
	String() string
}

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

func (r Region) String() string {
	switch r {
	case Europe:
		return "europe"
	case MiddleEast:
		return "middle east"
	case Asia:
		return "asia"
	case NorthAmerica:
		return "north america"
	case SouthAmerica:
		return "south america"
	case Oceania:
		return "oceania"
	case Africa:
		return "africa"
	case Global:
		return "global"
	}
	return "unknown"
}

type Versus string

const (
	Players Versus = "players"
	AI      Versus = "ai"
)

func (v Versus) String() string {
	return string(v)
}

type MatchType string

const (
	Unranked MatchType = "unranked"
	Custom   MatchType = "custom"
	EasyAI   MatchType = "aieasy"
	MediumAI MatchType = "aimedium"
	HardAI   MatchType = "aihard"
	ExpertAI MatchType = "aiexpert"
)

func (mt MatchType) String() string {
	return string(mt)
}

type TeamSize string

const (
	OneVOne     TeamSize = "1v1"
	TwoVTwo     TeamSize = "2v2"
	ThreeVThree TeamSize = "3v3"
	FourVFour   TeamSize = "4v4"
)

func (ts TeamSize) String() string {
	return string(ts)
}
