package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type meeting struct {
	subject    string
	start, end time.Time
	show       int
	allday     bool
}

func main() {

	// the csv file
	arg := os.Args
	inputfile := arg[1]
	f, _ := os.Open(inputfile)

	// date format
	dateformat := "1/2/2006 3:04:05 PM"
	timeformat := "3:04PM"
	defaultstart := "6:00:00 AM"
	defaultend := "3:00:00 PM"

	// counting lines
	lines := 0

	// create a new reader.
	r := csv.NewReader(bufio.NewReader(f))

	m := make(map[string][]meeting)

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
		ustart, err := time.Parse(dateformat, startdate+" "+starttime)
		check(err)

		uend, err := time.Parse(dateformat, enddate+" "+endtime)
		check(err)

		if showas == 2 {
			m[startdate] = append(m[startdate], meeting{
				subject: subject,
				start:   ustart,
				end:     uend,
				allday:  allday,
				show:    showas,
			})
		}
	}

	// sort the days
	days := make([]string, 0, len(m))
	for k := range m {
		days = append(days, k)
	}
	sort.Strings(days)

	// loop over the sorted days
	for d := range days {
		day := m[days[d]]

		daystart, err := time.Parse(dateformat, days[d]+" "+defaultstart)
		check(err)
		dayend, err := time.Parse(dateformat, days[d]+" "+defaultend)
		check(err)

		fmt.Printf("%v (%v - %v)\n", days[d], daystart.Format(timeformat), dayend.Format(timeformat))

		// sort meetings by start time
		sort.Slice(day, func(i, j int) bool {
			return day[i].start.Unix() < day[j].start.Unix()
		})

		// report on each day
		daytotal := 0
		freeblock := 0
		startfree := daystart
		endfree := dayend
		for i := range day {
			if i > 0 {
				startfree = day[i-1].end
			}
			if i < len(day) {
				endfree = day[i].start
			}

			freeblock = int((endfree.Unix() - startfree.Unix()) / 60)
			daytotal = daytotal + int((day[i].end.Unix() - day[i].start.Unix()))
			if freeblock > 0 {
				fmt.Printf("\tBlock Free: %d (%v - %v)\n", (endfree.Unix()-startfree.Unix())/60, startfree.Format(timeformat), endfree.Format(timeformat))
			}
			fmt.Printf("\t%v\n\t\t%v - %v = %d\n", day[i].subject, day[i].start.Format(timeformat), day[i].end.Format(timeformat), (day[i].end.Unix()-day[i].start.Unix())/60)
		}
		freeblock = int((endfree.Unix() - startfree.Unix()) / 60)
		if freeblock > 0 {
			fmt.Printf("Last Block Free: %d (%v - %v)\n", (endfree.Unix()-startfree.Unix())/60, startfree.Format(timeformat), endfree.Format(timeformat))
		}

		fmt.Printf("====================================\n\ttotal = %d\n\n", daytotal/60)
	}
}
