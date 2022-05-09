package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// shuffle the lines
	for i := len(lines) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		lines[i], lines[j] = lines[j], lines[i]
	}

	for _, line := range lines {
		fmt.Println(line)
	}
}
