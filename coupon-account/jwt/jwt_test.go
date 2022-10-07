package jwt

import (
	"fmt"
	"testing"
)

var singedkey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiTmlja25hbWUiOiJ5b3VuZyJ9.4Sdx-vASxsuLXDYHenZBCMPyRKOazOp9B0dnpZ0jnIY"

func TestJWT_GenerateJWT(t *testing.T) {
	jwt := NewJWT()
	generateJWT, err := jwt.GenerateJWT(CustomClaims{
		ID:       1,
		Nickname: "young",
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(generateJWT)
}

func TestJWT_ParseJWT(t *testing.T) {
	jwt := NewJWT()
	claims, err := jwt.ParseJWT(singedkey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(claims)
}

func TestJWT_RefreshJWT(t *testing.T) {
	jwt := NewJWT()
	refreshJWT, err := jwt.RefreshJWT(singedkey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(refreshJWT)
}
