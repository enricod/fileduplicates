package main

import (
	"crypto/sha256"
	"encoding/json"
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

func filtraHashesConPiuDiUnFile(hashes map[string][]string) map[string][]string {
	result := make(map[string][]string)
	for key, value := range hashes {
		if len(value) > 1 {
			result[key] = value
		}
	}
	return result
}
func main() {
	log.SetFlags(log.Lshortfile)

	hashes := make(map[string][]string)
	dir := os.Args[1]
	err := filepath.Walk(dir, fileWalker(hashes))
	if err != nil {
		log.Fatal(err)
	}

	/*
		totFiles := 0
		for key, value := range hashes {
			if len(value) > 1 {
				// fmt.Println("Key:", key, " => ", len(value))
				fmt.Println("Hash256 =", key, "| found", len(value), "files")
				for _, s := range value {
					fmt.Printf("\t %s\n", s)
				}
			}
			totFiles = totFiles + len(value)
		}
	*/

	b, err := json.MarshalIndent(filtraHashesConPiuDiUnFile(hashes), "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
	os.Stdout.Write([]byte("\n"))

	//fmt.Println("Files examined =", totFiles)
}
