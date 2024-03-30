package v1

import "fmt"
import "os"
import "bufio"
import "strings"
import "strconv"

type CityInfo struct {
    Min   float64
    Max   float64
    Sum   float64
    Count int
}

func parseLine(cityList map[string]*CityInfo, line string) {
    // Split the line by ';' to separate city name and temperature
    parts := strings.Split(line, ";")
    if len(parts) != 2 {
        fmt.Println("Invalid line format:", line)
        return
    }

    // Extract city name and temperature
    cityName := parts[0]
    temperatureStr := parts[1]

    // Convert temperature string to float64
    temperature, err := strconv.ParseFloat(temperatureStr, 64)
    if err != nil {
        fmt.Println("Error parsing temperature:", err)
        return
    }

    // Update cityList with the parsed information
    if cityInfo, ok := cityList[cityName]; ok {
        // Update existing city's statistics
        cityInfo.Count++
        cityInfo.Sum += temperature
        if temperature < cityInfo.Min || cityInfo.Count == 1 {
            cityInfo.Min = temperature
        }
        if temperature > cityInfo.Max || cityInfo.Count == 1 {
            cityInfo.Max = temperature
        }
    } else {
        // Create a new CityInfo entry if the city doesn't exist
        cityList[cityName] = &CityInfo{
            Min:   temperature,
            Max:   temperature,
            Sum:   temperature,
            Count: 1,
        }
    }
}


func main() {
    if len(os.Args) < 2 {
        fmt.Println("Please provide a file name as a command-line argument.")
        return
    }

	cityList := make(map[string]*CityInfo)

    filePath := "../Billion-Row-Challenge/" + os.Args[1]
    readFile, err := os.Open(filePath)
	if err != nil {
        fmt.Println(err)
    }

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
        parseLine(cityList, fileScanner.Text())
    }
	// Print the parsed information
    for cityName, cityInfo := range cityList {
        fmt.Printf("City: %s\n", cityName)
        fmt.Printf("Min temperature: %.2f\n", cityInfo.Min)
        fmt.Printf("Max temperature: %.2f\n", cityInfo.Max)
        fmt.Printf("Average temperature: %.2f\n", cityInfo.Sum/float64(cityInfo.Count))
    }
}