package v2

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

type CityInfo struct {
	Min   float64
	Max   float64
	Sum   float64
	Count int
}

func parseLine(cityList map[string]*CityInfo, ch <-chan string, wg *sync.WaitGroup) {
	defer wg.Done() // Decrement the wait group counter when done
	for line := range ch {
		parts := strings.Split(line, ";")
		if len(parts) != 2 {
			fmt.Println("Invalid line format:", line)
			continue
		}

		cityName := parts[0]
		temperatureStr := parts[1]

		temperature, err := strconv.ParseFloat(temperatureStr, 64)
		if err != nil {
			fmt.Println("Error parsing temperature:", err)
			continue
		}

		if cityInfo, ok := cityList[cityName]; ok {
			cityInfo.Count++
			cityInfo.Sum += temperature
			if temperature < cityInfo.Min || cityInfo.Count == 1 {
				cityInfo.Min = temperature
			}
			if temperature > cityInfo.Max || cityInfo.Count == 1 {
				cityInfo.Max = temperature
			}
		} else {
			cityList[cityName] = &CityInfo{
				Min:   temperature,
				Max:   temperature,
				Sum:   temperature,
				Count: 1,
			}
		}
	}
}

func v2() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a file name as a command-line argument.")
		return
	}

	cityList := make(map[string]*CityInfo)
	ch := make(chan string, 500000)

	var wg sync.WaitGroup
	wg.Add(1) // Increment the wait group counter

	go parseLine(cityList, ch, &wg)

	filePath := "../Billion-Row-Challenge/" + os.Args[1]
	readFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	for fileScanner.Scan() {
		ch <- fileScanner.Text()
	}
	close(ch)

	wg.Wait() // Wait for the parsing to finish

	for cityName, cityInfo := range cityList {
		fmt.Printf("City: %s\n", cityName)
		fmt.Printf("Min temperature: %.2f\n", cityInfo.Min)
		fmt.Printf("Max temperature: %.2f\n", cityInfo.Max)
		fmt.Printf("Average temperature: %.2f\n", cityInfo.Sum/float64(cityInfo.Count))
	}
}
