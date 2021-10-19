package pkg

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Tar struct {
}

func (t Tar) TarPackage(tarName string, paths map[string]string) error {
	errInfo := fmt.Sprintf("tar file %s create error", tarName)
	var err error

	tarFile, err := os.Create(tarName)
	if err != nil {
		return err
	}
	defer tarFile.Close()
	gw := gzip.NewWriter(tarFile)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	for path, dst := range paths {
		pi, err := os.Stat(path)
		if err != nil {
			return err
		}
		if pi.IsDir() {
			// path为目录情况下如下处理
			walker := func(f string, fi os.FileInfo, err error) error {
				hdr, err := tar.FileInfoHeader(fi, fi.Name())

				relFilePath, err := filepath.Rel(path, f)
				if err != nil {
					return err
				}
				if strings.HasSuffix(dst, "/") {
					hdr.Name = fmt.Sprintf("%s%s", dst, relFilePath)
				} else {
					hdr.Name = fmt.Sprintf("%s/%s", dst, relFilePath)
				}
				if err := tw.WriteHeader(hdr); err != nil {
					return err
				}
				if fi.Mode().IsDir() {
					return nil
				}
				srcFile, err := os.Open(f)
				defer srcFile.Close()
				_, err = io.Copy(tw, srcFile)
				if err != nil {
					return err
				}
				return nil
			}
			if err := filepath.Walk(path, walker); err != nil {
				err = fmt.Errorf("%s: failed to add %s to tar error: %s", errInfo, path, err.Error())
				return err
			}
		} else {
			hdr, err := tar.FileInfoHeader(pi, pi.Name())
			if err != nil {
				return err
			}

			if strings.HasSuffix(dst, "/") {
				hdr.Name = fmt.Sprintf("%s%s", dst, pi.Name())
			} else {
				hdr.Name = fmt.Sprintf("%s", dst)
			}
			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}
			srcFile, err := os.Open(path)
			if err != nil {
				return err
			}
			_, err = io.Copy(tw, srcFile)
			if err != nil {
				return err
			}
			srcFile.Close()
		}
	}
	return err
}
