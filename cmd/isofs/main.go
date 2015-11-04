package main

import (
	"io"
	"os"

	isofs "github.com/hooklift/iso"
)

func main() {
	file, err := os.Open("myfile.iso")
	if err != nil {
		panic(err)
	}

	r, err := isofs.NewReader(file)
	if err != nil {
		panic(err)
	}

	for r.Next() {
		fi, err := r.Value()
		if err != nil {
			panic(err)
		}

		if fi.IsDir() {
			os.MkdirAll(fi.Name(), fi.Mode().Perm())
		} else {
			freader := fi.Sys().(io.Reader)
			f, err := os.Create(fi.Name())
			if err != nil {
				panic(err)
			}

			if _, err := io.Copy(f, freader); err != nil {
				panic(err)
			}
			f.Close()
		}
	}
}
