package gamelogic

type Player struct {
	UserId        string
	Level         int8
	PlayerId      int64
	CurrentGameId *int64
	Wins          int32
	Losses        int32
	//TODO temp code
	Password string
}
