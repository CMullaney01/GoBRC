package main

import "fmt"
import "os"
import "bufio"

type CityInfo struct {
	min int
	max int
	sum int
	count int
}


func main() {
    if len(os.Args) < 2 {
        fmt.Println("Please provide a file name as a command-line argument.")
        return
    }

	// cityList := make(map[string]CityInfo)

    filePath := "../Billion-Row-Challenge/" + os.Args[1]
    readFile, err := os.Open(filePath)
	if err != nil {
        fmt.Println(err)
    }

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
        fmt.Println(fileScanner.Text())
    }
}