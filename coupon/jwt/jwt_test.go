package jwt

import (
	"fmt"
	"testing"
	"time"
)

var singedkey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjU5MzAzNDIsIklEIjoxLCJNb2JpbGUiOiJ5b3VuZyJ9.kcwl5AcUzRR3pQliSSVIWK8VawsjssOOAjTRX6stjHg"

func TestJWT_GenerateJWT(t *testing.T) {
	jwt := NewJWT()
	claims := CustomClaims{
		ID:     1,
		Mobile: "young",
	}
	generateJWT, err := jwt.GenerateJWT(claims)
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
	//bytes, _ := json.Marshal(&claims)
	//timestamp := fmt.Sprintf("%d", claims.ExpiresAt)
	//格式化为字符串,tm为Time类型

	tm := time.Unix(claims.ExpiresAt, 0)

	fmt.Println(time.Now().Before(tm))
	//expiredAt, _ := time.Parse("2006-04-02 15:04", timestamp)
	//fmt.Println(expiredAt)
}

func TestJWT_RefreshJWT(t *testing.T) {
	jwt := NewJWT()
	refreshJWT, err := jwt.RefreshJWT(singedkey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(refreshJWT)
}
