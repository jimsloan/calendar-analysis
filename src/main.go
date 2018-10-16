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

// cetralize error checks
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// the meeting data record
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

	// date formats
	csvdateformat := "1/2/2006 3:04:05 PM"
	dateformat := "2006/01/02 3:04:05 PM"
	timeformat := "3:04PM"

	// defaults
	defaultstart := "6:30:00 AM"
	defaultend := "3:30:00 PM"
	daylength := 540
	lunch := 60

	// counting
	lines := 0
	totalschedule := 0
	totalfree := 0

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
		showas, err := strconv.Atoi(record[6])
		check(err)

		// convert time to unixtime
		ustart, err := time.Parse(csvdateformat, startdate+" "+starttime)
		check(err)

		uend, err := time.Parse(csvdateformat, enddate+" "+endtime)
		check(err)

		// filter only meetings that show 'busy'
		if showas == 2 {
			dateindex := ustart.Format("2006/01/02")
			//println(dateindex)
			m[dateindex] = append(m[dateindex], meeting{
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
		totalschedule = totalschedule + daylength
		daystart, err := time.Parse(dateformat, days[d]+" "+defaultstart)
		check(err)
		dayend, err := time.Parse(dateformat, days[d]+" "+defaultend)
		check(err)

		// print the date header
		fmt.Printf("%v (%v - %v)\n------------------------\n", days[d], daystart.Format(timeformat), dayend.Format(timeformat))

		// sort meetings by start time
		sort.Slice(day, func(i, j int) bool {
			return day[i].start.Unix() < day[j].start.Unix()
		})

		// report on each day
		daytotal := 0
		freetotal := 0
		freeblock := 0
		freeblockminimum := 30
		meetingblockend := daystart

		// Beginning of the day
		freeblock = int((day[0].start.Unix() - daystart.Unix()) / 60)
		if freeblock > 0 {
			fmt.Printf("(%v - %v) Beginning Free Block: %d\n", daystart.Format(timeformat), day[0].start.Format(timeformat), freeblock)
			if freeblock > freeblockminimum {
				freetotal = freetotal + freeblock
			}
		}

		// loop through the day
		for i := range day {
			freeblock = 0
			// check for free time between meetings
			if i > 0 {

				// trying to account for double booking
				if meetingblockend.Unix() < day[i-1].end.Unix() {
					meetingblockend = day[i-1].end
				}

				freeblock = int((day[i].start.Unix() - meetingblockend.Unix()) / 60)
				if freeblock > 0 {
					fmt.Printf("(%v - %v) Free Block: %d\n", meetingblockend.Format(timeformat), day[i].start.Format(timeformat), freeblock)
					if freeblock > freeblockminimum {
						freetotal = freetotal + freeblock
					}
				}
			}

			fmt.Printf("\t(%v - %v) %v = %d\n", day[i].start.Format(timeformat), day[i].end.Format(timeformat), day[i].subject, (day[i].end.Unix()-day[i].start.Unix())/60)

			daytotal = daytotal + int((day[i].end.Unix() - day[i].start.Unix()))
		}
		if meetingblockend.Unix() < day[len(day)-1].end.Unix() {
			meetingblockend = day[len(day)-1].end
		}

		// End of the day
		freeblock = int((dayend.Unix() - meetingblockend.Unix()) / 60)
		if freeblock > 0 {
			fmt.Printf("(%v - %v) End of Day Free Block: %d\n", meetingblockend.Format(timeformat), dayend.Format(timeformat), freeblock)
			if freeblock > freeblockminimum {
				freetotal = freetotal + freeblock
			}
		}

		// daily totals
		freetotal = freetotal - lunch
		totalfree = totalfree + freetotal
		fmt.Printf("=================================\nfree = %d meetings = %d lunch = %d\n\n", freetotal, daytotal/60, lunch)
	}

	// report totals
	percentfree := (float64(totalfree) / float64(totalschedule) * float64(100))
	fmt.Printf("%d Days of %d for Total Scheduled = %d\nTotlal Free = %d (%0.2f %%)\n", len(days), daylength, totalschedule, totalfree, percentfree)
}
