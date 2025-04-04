package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Word struct {
	Gendered   string
	Substitute string
}

var dictionary map[string]Word = make(map[string]Word)
var punctuationBuffer byte = ' '

func loadDictionary() {

	data, err := os.ReadFile("gender.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(data, &dictionary)
	if err != nil {
		log.Println(err)
		return
	}
}

func checkFile(sourceFileName string) {
	lineNumber := 0
	file, err := os.Open(sourceFileName)
	if err != nil {
		log.Printf("%s is not a filename.", sourceFileName)
		word, ok := dictionary[sourceFileName]
		if ok {
			log.Println("Word " + sourceFileName + " found in dictionary suggest replacing with " + word.Substitute)
		}
		return
	}

	Scanner := bufio.NewScanner(file)
	Scanner.Split(bufio.ScanWords)
	for Scanner.Scan() {
		punctuationBuffer = Scanner.Text()[len(Scanner.Text())-1]
		var searchWord string
		if punctuationBuffer == '.' || punctuationBuffer == ',' || punctuationBuffer == '?' {
			searchWord = Scanner.Text()[:len(Scanner.Text())-1]
		} else {
			searchWord = Scanner.Text()
		}

		word, ok := dictionary[searchWord]
		if ok {
			log.Printf("%s found in %s on line %d replace with %s", searchWord, sourceFileName, lineNumber, word.Substitute)
		}
		lineNumber++
	}

	if err := Scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	var arguments = os.Args
	if len(arguments) == 1 {
		log.Fatal("You need to at least pass a filename as parameter.")
	}

	if len(arguments) == 2 {
		var fileName = arguments[1]
		loadDictionary()
		checkFile(fileName)
	}
}
