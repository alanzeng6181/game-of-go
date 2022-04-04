package gamelogic

type GameState struct {
	Positions     []Stone
	WhiteCaptures int8
	BlackCaptures int8
	Status        Status
	ClockState    ClockState
	Comments      []Comment
}
