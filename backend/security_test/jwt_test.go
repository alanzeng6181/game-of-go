package security_test

import (
	"fmt"
	"testing"

	"github.com/alanzeng6181/game-of-go/security"
)

func TestJWT(t *testing.T) {
	userId := "player1"
	token, err := security.GetToken(userId)
	if err != nil {
		t.Error(err)
		return
	}
	userIdActual, err := security.GetUserId(token)

	if err != nil {
		t.Error(err)
	}

	if userIdActual != userId {
		t.Errorf(fmt.Sprintf("expected %s, but got %s", userId, userIdActual))
	}
}
