package captcha

import (
	"fmt"
	"testing"
)

func TestGenCaptcha(t *testing.T) {
	res, err := GenCaptcha("123")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Sprintf(res)
}
