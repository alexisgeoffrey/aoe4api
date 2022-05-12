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

type Versus int

const (
	Players Versus = iota
	AI
)

func (v Versus) String() string {
	switch v {
	case Players:
		return "players"
	case AI:
		return "ai"
	}
	return "unknown"
}

type MatchType int

const (
	Unranked MatchType = iota
	Custom
	EasyAI
	MediumAI
	HardAI
	ExpertAI
)

func (mt MatchType) String() string {
	switch mt {
	case Unranked:
		return "unranked"
	case Custom:
		return "custom"
	case EasyAI:
		return "aieasy"
	case MediumAI:
		return "aimedium"
	case HardAI:
		return "aihard"
	case ExpertAI:
		return "aiexpert"
	}
	return "unknown"
}

type TeamSize int

const (
	OneVOne TeamSize = iota + 1
	TwoVTwo
	ThreeVThree
	FourVFour
)

func (ts TeamSize) String() string {
	switch ts {
	case OneVOne:
		return "1v1"
	case TwoVTwo:
		return "2v2"
	case ThreeVThree:
		return "3v3"
	case FourVFour:
		return "4v4"
	}
	return "unknown"
}
