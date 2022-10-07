package password

import (
	"crypto/md5"
	"github.com/anaskhan96/go-password-encoder"
)

var encryptOptions = password.Options{
	SaltLen:      16,
	Iterations:   100,
	KeyLen:       32,
	HashFunction: md5.New,
}

// GenerateHashedPwd return salt,hashed password
func GenerateHashedPwd(pwd string) (string, string) {
	return password.Encode(pwd, &encryptOptions)
}

func CheckPwd(pwd, salt, hsdPwd string) bool {
	return password.Verify(pwd, salt, hsdPwd, &encryptOptions)
}
