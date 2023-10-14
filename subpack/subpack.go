package subpack

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Node struct {
	hostname string
	ipaddr   string
	pstatus  string
	ppid     int
}
type ZlogQso struct {
	//2015/04/25 21:00 JR2NMJ       599 15H     599 21M     21    -     3.5  CW   1  %%JG1AVR%%         TX#1
	//2015/04/25 21:01 JH1XDW       59  15H     59  14M     14    -     7    SSB  1  %%JR1TCY%%         TX#2
	DateQSO  string
	TimeQSO  string
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

func main333() {
	data := `2015/04/25 21:00 JR2NMJ       599 15H     599 21M     21    -     3.5  CW   1  %%JG1AVR%%         TX#1`
	lines := strings.Split(data, "\n")

	for _, line := range lines {
		fields := strings.Fields(line)
		record := Record{
			Date:     fields[0],
			Time:     fields[1],
			CallSign: fields[2],
			RST:      fields[3] + " " + fields[4],
			Band:     fields[8],
			Mode:     fields[9],
			TX:       fields[11],
		}
		fmt.Println(record)
	}
}

func Readfile() []*ZlogQso {
	fileName := "2015-ALLJA-0.all" // replace with your file name
	file, err := os.Open(fileName)
	if err != nil {
		panic(fmt.Sprintf("error opening %s: %v", fileName, err))
	}
	defer file.Close()

	//var nodes []*Node
	var zlogqso []*ZlogQso

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		//nodes = append(nodes, &Node{ipaddr: scanner.Text()})
		fields := strings.Fields(scanner.Text())
		if len(fields) < 14 {
			continue
		}
		zlogqso = append(zlogqso, &ZlogQso{
			DateQSO:  fields[0],
			TimeQSO:  fields[1],
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
	for _, qso := range zlogqso {
		fmt.Println(qso)
	}
	return zlogqso
}
