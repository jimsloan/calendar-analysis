package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	// the csv file
	arg := os.Args
	inputfile := arg[1]
	f, _ := os.Open(inputfile)

	// date format
	layout := "1/2/2006 3:04:05 PM"

	// counting lines
	lines := 0

	// create a new reader.
	r := csv.NewReader(bufio.NewReader(f))

	// read and process
	for {
		record, err := r.Read()
		lines++
		if lines == 1 {
			// skip header line
			continue
		}
		// stop at EOF.
		if err == io.EOF {
			break
		}

		// map csv columns
		subject := record[0]
		startdate := record[1]
		starttime := record[2]
		enddate := record[3]
		endtime := record[4]
		allday := false
		if strings.EqualFold("true", record[5]) {
			allday = true
		}
		showas, err := strconv.Atoi(record[21])
		check(err)

		// convert time to unixtime
		ustart, err := time.Parse(layout, startdate+" "+starttime)
		check(err)

		uend, err := time.Parse(layout, enddate+" "+endtime)
		check(err)

		// display results
		if showas == 2 {
			fmt.Printf("%v\n", subject)
			fmt.Printf("%v ALLDAY:%t SHOWAS:%v\n", startdate, allday, showas)
			fmt.Printf("%d - %d = %d\n\n", ustart.Unix(), uend.Unix(), (uend.Unix()-ustart.Unix())/60)
		}

	}
}
