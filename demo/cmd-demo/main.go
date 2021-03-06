package main

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"syscall"
)

func commandExec(command, workDir string) (string, string, error) {
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

func extractEmbedFile(rootDir string) error {
	return fs.WalkDir(f, rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path != "." {
			if d.IsDir() {
				_ = os.MkdirAll(path, 0700)
			} else {
				bs, err := f.ReadFile(path)
				if err != nil {
					fmt.Println("ERROR:", err.Error())
					return err
				}
				_ = os.WriteFile(path, bs, 0600)
			}
		}
		fmt.Println("#", path)
		return nil
	})
}

//go:embed shell/* shell/.test_dir/.test.txt shell/.test_dir/_helper.tpl shell2
var f embed.FS

func main() {
	err := extractEmbedFile("shell")
	if err != nil {
		err = fmt.Errorf("walk error: %s", err.Error())
		fmt.Println("[ERROR]", err.Error())
		return
	}

	strOut, strErr, err := commandExec("sh shell/certs.sh", "")
	if err != nil {
		fmt.Println("[ERROR]", err.Error())
	} else {
		fmt.Println("[SUCCESS]")
	}
	fmt.Println("##############")
	fmt.Println("stdout:")
	fmt.Println(strOut)
	fmt.Println("stderr:")
	fmt.Println(strErr)
}
