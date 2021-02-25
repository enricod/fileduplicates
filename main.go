package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func fileWalker(hashes map[string][]string) filepath.WalkFunc {

	return func(path string, info os.FileInfo, err error) error {
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
			h := sha256.New()
			if _, err := io.Copy(h, f); err != nil {
				log.Fatal(err)
			}
			sha256Sum := fmt.Sprintf("%x", h.Sum(nil))
			hashes[sha256Sum] = append(hashes[sha256Sum], path) // md5s[md5value] = path
			// fmt.Printf("%s \t-> %s \n", path, md5value)
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

	totFiles := 0
	for key, value := range md5s {
		if len(value) > 1 {
			// fmt.Println("Key:", key, " => ", len(value))
			fmt.Println("Hash256 =", key, "| found", len(value), "files")
			for _, s := range value {
				fmt.Printf("\t %s\n", s)
			}
		}
		totFiles = totFiles + len(value)
	}

	fmt.Println("Files examined =", totFiles)
}
