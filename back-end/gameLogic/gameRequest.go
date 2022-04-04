package gamelogic

type GameRequest struct {
	RequestId         int64
	PlayerRank        int8
	PlayerId          int64
	RankAbove         uint8
	RankBelow         uint8
	BoardSize         uint8
	TimeoutInSeconds  uint32
	OverTimeInSeconds uint32
	OverTimeCount     uint8
}

func (gameRequest GameRequest) IsMatch(other GameRequest) bool {
	//TODO implement matching logic
	return true
}
