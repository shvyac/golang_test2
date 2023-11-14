package subpack

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type ZlogQso struct {
	//2015/04/25 21:00 AAAAAA       599 15H     599 21M     21    -     3.5  CW   1  %%CCCCCC%%         TX#1
	//2015/04/25 21:01 BBBBBB       59  15H     59  14M     14    -     7    SSB  1  %%CCCCCC%%         TX#2
	DateTime time.Time
	Callsign string
	RSTsent  string
	NRsent   string
	RSTrcvd  string
	NRrcvd   string
	Mult     string
	Mult2    string
	Band     string
	Mode     string
	Point    string
	Oper     string
	TxNo     string
	//series 1 and 2
	SeriesNo    int
	ElapsedTime int64
	TxCount     int
	QsoCount    int
}

type Record struct {
	Date     string
	Time     string
	CallSign string
	RST      string
	Band     string
	Mode     string
	TX       string
}

func Readfile() []*ZlogQso {
	fileName := "2015-ALLJA-0.all" // replace with your file name
	file, err := os.Open(fileName)
	if err != nil {
		panic(fmt.Sprintf("error opening %s: %v", fileName, err))
	}
	defer file.Close()
	var zlogqso []*ZlogQso
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 14 {
			continue //skip header
		}
		qsodatetime := ToJstTimeFromString(fields[0] + fields[1])
		if fields[2] == "JH4WBY" {
			fmt.Println("JH4WBY: ", fields[2], qsodatetime, fields[9])
		}
		zlogqso = append(zlogqso, &ZlogQso{
			DateTime: qsodatetime,
			Callsign: fields[2],
			RSTsent:  fields[3],
			NRsent:   fields[4],
			RSTrcvd:  fields[5],
			NRrcvd:   fields[6],
			Mode:     fields[7],
			Mult2:    fields[8],
			Band:     fields[9],
			Mult:     fields[10],
			Point:    fields[11],
			Oper:     strings.ReplaceAll(fields[12], "%%", ""),
			TxNo:     fields[13],
		})
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading from file:", err)
		os.Exit(3)
	}

	fmt.Println("ZlogQso: ", len(zlogqso), "records")
	// for _, qso := range zlogqso {
	// 	fmt.Println(qso)
	// }
	return zlogqso
}

func ToJstTimeFromString(timeString string) time.Time {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	layout := "2006/01/0215:04"
	ts, err := time.ParseInLocation(layout, timeString, loc)
	if err != nil {
		fmt.Println(err)
	}
	return ts
}
