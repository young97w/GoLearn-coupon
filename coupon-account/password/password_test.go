package password

import (
	"fmt"
	"testing"
)

var pwd = "65535"
var hsdPwd = "c179d9618ce1565222f02c6b2494022d0059dfca446616304994d402711fe632"
var salt = "8xmrwb5ouK9cVfoM"

func TestGenerateHashedPwd(t *testing.T) {
	fmt.Println(GenerateHashedPwd("65535"))
}

func TestCheckPwd(t *testing.T) {
	if !CheckPwd(pwd, salt, hsdPwd) {
		t.Fatal("")
	}
	fmt.Println("check true")
}
