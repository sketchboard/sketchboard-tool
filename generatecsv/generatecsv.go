package main

import (
	"encoding/csv"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/tjarratt/babble"
)

func main() {
	babbler := babble.NewBabbler()

	// set the number of words you want
	babbler.Count = 10
	babbler.Separator = " "
	println(babbler.Babble()) // antibiomicrobrial (or some other word)

	rows := make([][]string, 0)
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	for i := 0; i < 50; i++ {
		babbler.Count = r.Intn(20)
		words := babbler.Babble()
		rows = append(rows, []string{words})
	}

	output := os.Args[1]
	csvfile, err := os.Create(output)
	if err != nil {
		log.Fatalln(err)
	}

	writer := csv.NewWriter(csvfile)

	for _, row := range rows {
		err := writer.Write(row)
		if err != nil {
			log.Fatalln(err)
		}
	}

	writer.Flush()
}
