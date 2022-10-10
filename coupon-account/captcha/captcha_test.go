package captcha

import (
	"fmt"
	"testing"
)

func TestGenCaptcha(t *testing.T) {
	_, res, err := GenCaptcha()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Sprintf(res)
}
