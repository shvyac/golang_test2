package subpack

import (
	"bufio"
	"fmt"
	"os"
)

type Node struct {
	hostname string
	ipaddr   string
	pstatus  string
	ppid     int
}

func Main22() {
	fileName := "2015-ALLJA-0.all" // replace with your file name
	file, err := os.Open(fileName)
	if err != nil {
		panic(fmt.Sprintf("error opening %s: %v", fileName, err))
	}
	defer file.Close()

	var nodes []*Node

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		nodes = append(nodes, &Node{ipaddr: scanner.Text()})
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading from file:", err)
		os.Exit(3)
	}

	for _, node := range nodes {
		fmt.Println(node.ipaddr)
	}
}
