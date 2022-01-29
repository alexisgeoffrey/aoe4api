package aoe4api

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
