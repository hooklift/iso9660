package main

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hooklift/iso9660"
)

func main() {
	file, err := os.Open("test.iso")
	if err != nil {
		panic(err)
	}

	r, err := iso9660.NewReader(file)
	if err != nil {
		panic(err)
	}

	destPath := "tmp"
	for {
		fi, err := r.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		fp := filepath.Join(destPath, fi.Name())
		if fi.IsDir() {
			if err := os.MkdirAll(fp, fi.Mode()); err != nil {
				panic(err)
			}
			continue
		}

		parentDir, _ := filepath.Split(fp)
		if err := os.MkdirAll(parentDir, fi.Mode()); err != nil {
			panic(err)
		}

		freader := fi.Sys().(io.Reader)
		f, err := os.Create(fp)
		if err != nil {
			panic(err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				panic(err)
			}
		}()

		if err := f.Chmod(fi.Mode()); err != nil {
			panic(err)
		}

		if _, err := io.Copy(f, freader); err != nil {
			panic(err)
		}
	}
}

// Sanitizes name to avoid overwriting sensitive system files when unarchiving
func sanitize(name string) string {
	// Gets rid of volume drive label in Windows
	if len(name) > 1 && name[1] == ':' && runtime.GOOS == "windows" {
		name = name[2:]
	}

	name = filepath.Clean(name)
	name = filepath.ToSlash(name)
	for strings.HasPrefix(name, "../") {
		name = name[3:]
	}
	return name
}
