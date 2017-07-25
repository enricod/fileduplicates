package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"crypto/md5"
	"io"
	"strings"
)



func fileWalker( md5s map[string][]string) filepath.WalkFunc {

	return func (path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Print(err)
			return nil
		}

		if !info.IsDir() && !strings.HasSuffix(path, "sock") {
			f, err := os.Open(path)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			defer f.Close()
			h := md5.New()
			if _, err := io.Copy(h, f); err != nil {
				fmt.Println(err)
				return nil
			}


			md5value := fmt.Sprintf("%x",  h.Sum(nil))
			md5s[md5value] = append(md5s[md5value], path) // md5s[md5value] = path

			fmt.Printf("%s \t-> %s \n", path, md5value)
		}
		return nil
	}
}

func main() {
	log.SetFlags(log.Lshortfile)

	md5s := make(map[string][]string)
	dir := os.Args[1]
	err := filepath.Walk(dir, fileWalker(md5s))
	if err != nil {
		log.Fatal(err)
	}

	for key, value := range md5s {
		if len(value) > 1 {
			// fmt.Println("Key:", key, " => ", len(value))
			fmt.Println("Key:", key, " => ", len(value))
			for _, s := range value {
				fmt.Printf("\t %s\n", s)
			}
		}

	}
}

