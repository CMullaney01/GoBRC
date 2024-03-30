package main

import (
	"fmt"
	"os"
	"strconv"
	"github.com/edsrzf/mmap-go"
	"sort"
	"sync"
)	

type CityInfo struct {
	Min   int
	Max   int
	Sum   int
	Count int
}

func updateMap(cityList map[string]*CityInfo, cityName string, temperature int) {
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

func parseChunk(chunkData *mmap.MMap, cityList map[string]*CityInfo, start, end, fileSize int) {
	fmt.Println(start, end)
	// Forward to start of next line
	for (start != 0) && ((*chunkData)[start-1] != '\n') {
		start++
	}

	for start < end && start < fileSize {
		cityName := ""
		// Find the position of ';' character or end of chunk
		endIndex := start
		for (*chunkData)[endIndex] != ';' {
			endIndex++
		}

		// Extract city name from chunkData slice
		cityName = string((*chunkData)[start:endIndex])
		start = endIndex + 1 // Move past the ';' character
		fmt.Println(cityName)
		tmp := []byte{}
		// Extract temperature until '\n' character
		// asssumes that the last line ends with \n
		for (*chunkData)[start] != '\n' && start < fileSize{
			if (*chunkData)[start] != '.' {
				tmp = append(tmp, (*chunkData)[start])
			}
			start++
		}

		// Convert temperature bytes to integer
		temperature, err := strconv.Atoi(string(tmp))
		fmt.Println(temperature)
		if err != nil {
			fmt.Println("Error converting temperature to integer:", err)
			continue
		}

		updateMap(cityList, cityName, temperature)
	}
}


func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func min(a, b int) int {
	if a < b {
        return a
    }
    return b
}

func mergeCityLists(cityLists []map[string]*CityInfo) map[string]*CityInfo {
	merged := make(map[string]*CityInfo)
	for _, cityList := range cityLists {
		for city, info := range cityList {
			if _, ok := merged[city]; ok {
				merged[city].Min = min(merged[city].Min, info.Min)
				merged[city].Max = max(merged[city].Max, info.Max)
				merged[city].Sum += info.Sum
				merged[city].Count += info.Count
			} else {
				merged[city] = info
			}
		}
	}
	return merged
}

func printCityList(cityList map[string]*CityInfo) {
	var keys []string
	for city := range cityList {
		keys = append(keys, city)
	}
	sort.Strings(keys)

	for _, city := range keys {
		info := cityList[city]
		min_temp :=  float64(info.Min) / 10
		max_temp := float64(info.Max) / 10

		fmt.Printf("City: %s\n", city)
		fmt.Printf("Min temperature: %d\n", min_temp)
		fmt.Printf("Max temperature: %d\n", max_temp)
		fmt.Printf("Average temperature: %.2f\n", (float64(info.Sum)/float64(info.Count)) / 10)
		fmt.Println()
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a file name as a command-line argument.")
		return
	}

	filePath := "../Billion-Row-Challenge/" + os.Args[1]
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
        panic(err)
    }
    defer file.Close()

	fileInfo, err := file.Stat()
    if err != nil {
        panic(err)
    }
    fileSize := fileInfo.Size()

	// Memory map the file
    data, err := mmap.Map(file, mmap.RDWR, 0)
    if err != nil {
        panic(err)
    }
    defer data.Unmap()

	num_chunks := 1
	chunkSize := int(fileSize) / num_chunks
	dataPtr := &data
    var wg sync.WaitGroup

	// Create a slice of maps to store city lists for each chunk
	cityLists := make([]map[string]*CityInfo, num_chunks)
	for i := range cityLists {
		cityLists[i] = make(map[string]*CityInfo)
	}

	// Iterate over the data and split into chunks
	for i := 0; i < num_chunks; i++ {
		start := i * chunkSize
		end := (i + 1) * chunkSize
		if i == num_chunks-1 {
			// Last chunk might be smaller if the file size is not evenly divisible
			end = int(fileSize)
		}

		wg.Add(1)

		go func(i, start, end int) {
            defer wg.Done() // Decrement the wait group counter when the goroutine completes
            parseChunk(dataPtr, cityLists[i], start, end, int(fileSize))
        }(i, start, end)
	}
	wg.Wait()
	fmt.Println("finished parsing chunks")

	// merge and sort city lists
	cityList := mergeCityLists(cityLists)

	// Print the parsed information
	printCityList(cityList)
}

