package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"unicode"
)

func checkWord(word string, letters []rune) bool {
	checked := make([]bool, len(letters))
	for _, word_letter := range word {
		found := false
		for pos, letter := range letters {
			if unicode.ToLower(word_letter) == unicode.ToLower(letter) && !checked[pos] {
				checked[pos] = true
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func worker(jobs <-chan string, done chan<- bool, letters []rune) {
	for {
		j, more := <-jobs
		if more {
			if checkWord(j, letters) {
				fmt.Printf("%s\n", j)
			}
		} else {
			done <- true
			return
		}
	}
}

func main() {
	letters := []rune{'a', 'g', 'e', 'n', 't', 'u', 'r'}
	// initialize max number of concurrent jobs
	numGoroutines := runtime.GOMAXPROCS(0)
	// create channels
	jobs := make(chan string, numGoroutines)
	done := make(chan bool, numGoroutines)

	// start workers
	for w := 0; w < numGoroutines; w++ {
		go worker(jobs, done, letters)
	}

	// open file
	file, err := os.Open("dictionary.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// read file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// and put the lines into the jobs channel
		// fmt.Printf("> %s\n", scanner.Text())
		jobs <- scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	close(jobs)

	for i := 0; i < numGoroutines; i++ {
		<-done
	}
}
