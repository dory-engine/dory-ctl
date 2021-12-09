package pkg

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

func CheckRandomStringStrength(password string, length int, enableSpecialChar bool) error {
	var err error

	if len(password) < length {
		err = fmt.Errorf("password must at least %d charactors", length)
		return err
	}
	lowerChars := "abcdefghijklmnopqrstuvwxyz"
	upperChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberChars := "0123456789"
	specialChars := `~!@#$%^&*()_+-={}[]\|:";'<>?,./`
	lowerOK := strings.ContainsAny(password, lowerChars)
	upperOK := strings.ContainsAny(password, upperChars)
	numberOK := strings.ContainsAny(password, numberChars)
	specialOK := strings.ContainsAny(password, specialChars)
	if enableSpecialChar && !(lowerOK && upperOK && numberOK && specialOK) {
		err = fmt.Errorf("password must include lower upper case charactors and number and special charactors")
		return err
	} else if !enableSpecialChar && !(lowerOK && upperOK && numberOK) {
		err = fmt.Errorf("password must include lower upper case charactors and number")
		return err
	}

	return err
}

func RandomString(n int, enableSpecialChar bool, suffix string) string {
	var letter []rune
	lowerChars := "abcdefghijklmnopqrstuvwxyz"
	upperChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberChars := "0123456789"
	specialChars := `~!@#$%^&*()_+-={}[]\|:";'<>?,./`
	if enableSpecialChar {
		chars := fmt.Sprintf("%s%s%s%s", lowerChars, upperChars, numberChars, specialChars)
		letter = []rune(chars)
	} else {
		chars := fmt.Sprintf("%s%s%s", lowerChars, upperChars, numberChars)
		letter = []rune(chars)
	}
	var pwd string
	for {
		b := make([]rune, n)
		seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
		for i := range b {
			b[i] = letter[seededRand.Intn(len(letter))]
		}
		pwd = string(b)
		err := CheckRandomStringStrength(pwd, n, enableSpecialChar)
		if err == nil {
			break
		}
	}
	return fmt.Sprintf("%s%s", pwd, suffix)
}

func ValidateIpAddress(s string) error {
	var err error
	if net.ParseIP(s).To4() == nil {
		err = fmt.Errorf(`not ipv4 address`)
		return err
	}
	return err
}
