package pkg

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

func ExtractEmbedFile(f embed.FS, rootDir string, targetDir string) error {
	return fs.WalkDir(f, rootDir, func(path string, d fs.DirEntry, err error) error {
		rootDir = strings.TrimSuffix(rootDir, "/")
		if err != nil {
			return err
		}
		if path != "." && path != rootDir {
			pathTarget := fmt.Sprintf("%s/%s", targetDir, strings.TrimPrefix(path, fmt.Sprintf("%s/", rootDir)))
			if d.IsDir() {
				_ = os.MkdirAll(pathTarget, 0700)
			} else {
				bs, err := f.ReadFile(path)
				if err != nil {
					fmt.Println("ERROR:", err.Error())
					return err
				}
				_ = os.WriteFile(pathTarget, bs, 0600)
			}
		}
		return nil
	})
}

func CommandExec(command, workDir string) (string, string, error) {
	var err error
	errInfo := fmt.Sprintf("exec %s error", command)
	var strOut, strErr string

	execCmd := exec.Command("sh", "-c", command)
	execCmd.Dir = workDir

	prOut, pwOut := io.Pipe()
	prErr, pwErr := io.Pipe()
	execCmd.Stdout = pwOut
	execCmd.Stderr = pwErr

	rOut := io.TeeReader(prOut, os.Stdout)
	rErr := io.TeeReader(prErr, os.Stderr)

	err = execCmd.Start()
	if err != nil {
		err = fmt.Errorf("%s: exec start error: %s", errInfo, err.Error())
		return strOut, strErr, err
	}

	var bOut, bErr bytes.Buffer

	go func() {
		_, _ = io.Copy(&bOut, rOut)
	}()

	go func() {
		_, _ = io.Copy(&bErr, rErr)
	}()

	err = execCmd.Wait()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				err = fmt.Errorf("%s: exit status: %d", errInfo, status.ExitStatus())
			}
		} else {
			err = fmt.Errorf("%s: exec run error: %s", errInfo, err.Error())
			return strOut, strErr, err
		}
	}

	strOut = bOut.String()
	strErr = bErr.String()

	return strOut, strErr, err
}

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
