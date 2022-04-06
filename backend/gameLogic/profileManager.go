package gamelogic

import (
	"errors"
)

type ProfileManager struct {
	players []Player
}

func NewProfileManager() *ProfileManager {
	pm := ProfileManager{
		players: make([]Player, 0),
	}

	pm.players = append(pm.players, Player{
		UserId:   "player1",
		Level:    1,
		PlayerId: 1,
		Password: "Password",
	}, Player{
		UserId:   "player1",
		Level:    1,
		PlayerId: 1,
		Password: "Password",
	})
	return &pm
}

func (pm *ProfileManager) UpdateStat(winnerPlayerId int64, loserPlayerId int64) {

}

func (pm *ProfileManager) GetPlayer(playerId int64) *Player {
	for _, p := range pm.players {
		if p.PlayerId == playerId {
			return &p
		}
	}
	return nil
}

func (pm *ProfileManager) GetPlayerByUserId(userId string) *Player {
	for _, p := range pm.players {
		if p.UserId == userId {
			return &p
		}
	}
	return nil
}

func (pm *ProfileManager) GetPlayerId(user string, password string) (int64, error) {
	for _, p := range pm.players {
		if p.UserId == user && p.Password == password {
			return p.PlayerId, nil
		}
	}
	return -1, errors.New("invalid credentail")
}

func (pm *ProfileManager) GetUserId(user string, password string) (string, error) {
	for _, p := range pm.players {
		if p.UserId == user && p.Password == password {
			return p.UserId, nil
		}
	}
	return -1, errors.New("invalid credentail")
}
