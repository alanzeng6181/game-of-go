package gameLogic_test

import (
	"encoding/json"
	"testing"

	gamelogic "github.com/alanzeng6181/game-of-go/gameLogic"
)

func TestGameRequestJson(t *testing.T) {
	gameRequestStr := `
	{
		"PlayerRank":3,
		"PlayerId":1,
		"RankAbove":1,
		"RankBelow":1,
		"BoardSize":19,
		"TimeoutInSeconds" : 300,
		"OverTimeInSeconds" :  30,
		"OverTimeCount": 3
	}
	`
	var gameRequest gamelogic.GameRequest
	err := json.Unmarshal([]byte(gameRequestStr), &gameRequest)
	if err != nil {
		t.Error(err)
	}

	assertEqual(19, gameRequest.BoardSize, t)
	assertEqual(1, gameRequest.PlayerId, t)
	assertEqual(3, gameRequest.PlayerRank, t)
	assertEqual(1, gameRequest.RankAbove, t)
	assertEqual(1, gameRequest.RankBelow, t)
	assertEqual(300, gameRequest.TimeoutInSeconds, t)
	assertEqual(30, gameRequest.OverTimeInSeconds, t)
	assertEqual(3, gameRequest.OverTimeCount, t)
}
