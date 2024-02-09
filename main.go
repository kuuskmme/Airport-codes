package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

func main() {
	// Command-line flag
	helpFlag := flag.Bool("h", false, "Display help")
	flag.Parse()

	if *helpFlag {
		fmt.Println("Itinerary usage:\n go run . ./input.txt ./output.txt ./airport-lookup.csv")
		return
	}

	// Validate arguments
	args := flag.Args()
	if len(args) != 3 {
		fmt.Println("Incorrect number of arguments")
		fmt.Println("Itinerary usage:\n go run . ./input.txt ./output.txt ./airport-lookup.csv")
		return
	}

	inputFile, outputFile, lookupFile := args[0], args[1], args[2]

	// Process itinerary
	err := processItinerary(inputFile, outputFile, lookupFile)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// Function to process the itinerary
func processItinerary(inputFile, outputFile, lookupFile string) error {
	// Read and parse airport lookup
	airportLookup, err := parseAirportLookup(lookupFile)
	if err != nil {
		return err
	}

	//Read input file
	input, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("Input not found")
	}

	//Process text
	processedText := processText(string(input), airportLookup)

	// Write to output file
	err = os.WriteFile(outputFile, []byte(processedText), 0644)
	if err != nil {
		return fmt.Errorf("Error writing to output file")
	}

	return nil
}

func parseAirportLookup(filepath string) (map[string]string, error) {
	// Open file
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("Airport lookup not found")
	}
	defer file.Close()

	// Read .csv content
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Airport lookup malformed")
	}

	// Process records
	lookup := make(map[string]string)
	for i, record := range records {
		if i == 0 { // Skip header row
			continue
		}
		if len(record) != 6 || record[0] == "" || record[3] == "" || record[4] == "" {
			return nil, fmt.Errorf("Airport lookup malformed")
		}

		// Map both IATA and ICAO codes to the airport name
		lookup["#"+record[4]] = record[0]  // IATA
		lookup["##"+record[3]] = record[0] // ICAO
	}

	return lookup, nil
}

func processText(text string, airportLookup map[string]string) string {
	// Replace airport codes
	for code, name := range airportLookup {
		text = strings.ReplaceAll(text, code, name)
	}

	// Replace D dates from the first code
	text = regexp.MustCompile(`D\(([^)]+)\)`).ReplaceAllStringFunc(text, func(match string) string {
		dateString := match[2 : len(match)-1]
		date, err := time.Parse("2006-01-02T15:04-07:00", dateString)
		if err != nil {
			date, err = time.Parse("2006-01-02T15:04Z", dateString)
			if err != nil {
				return match
			}
		}
		return date.Format("02 Jan 2006")
	})

	// Replace T12 times from the first code
	text = regexp.MustCompile(`T12\(([^)]+)\)`).ReplaceAllStringFunc(text, func(match string) string {
		timeString := match[4 : len(match)-1]
		t, err := time.Parse("2006-01-02T15:04-07:00", timeString)
		if err != nil {
			t, err = time.Parse("2006-01-02T15:04Z", timeString)
			if err != nil {
				return match
			}
		}
		return t.Format("03:04PM (-07:00)")
	})

	// Replace T24 times from the first code
	text = regexp.MustCompile(`T24\(([^)]+)\)`).ReplaceAllStringFunc(text, func(match string) string {
		timeString := match[4 : len(match)-1]
		t, err := time.Parse("2006-01-02T15:04-07:00", timeString)
		if err != nil {
			t, err = time.Parse("2006-01-02T15:04Z", timeString)
			if err != nil {
				return match
			}
		}
		return t.Format("15:04 (-07:00)")
	})

	// Replace line-break characters with \n and remove multiple consecutive blank lines
	text = strings.Replace(text, "\\v", "\n", -1)
	text = strings.Replace(text, "\\f", "\n", -1)
	text = strings.Replace(text, "\\r", "\n", -1)
	text = regexp.MustCompile(`\n{3,}`).ReplaceAllString(text, "\n\n")

	// // Remove multiple consecutive blank lines
	text = RemoveExtraNewLines(text)

	return text
}

func formatDate(input, layout string) string {
	// Extract the date from the matched string
	dateStr := strings.TrimPrefix(input, "D(")
	dateStr = strings.TrimSuffix(dateStr, ")")

	// Parse date
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return input // return original on error
	}

	// Return formatted date
	return t.Format(layout)
}

func formatTime(input string, is12HourFormat bool) string {
	// Extract the time part from the matched string
	timeStr := strings.TrimSuffix(strings.TrimPrefix(input, "T12("), ")")
	timeStr = strings.TrimSuffix(strings.TrimPrefix(timeStr, "T24("), ")")

	// Parse time
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return input // return original on error
	}

	var layout string
	if is12HourFormat {
		// 12-hour format
		layout = "03:04PM"
	} else {
		// 24-hour format
		layout = "15:04"
	}

	formattedTime := t.Format(layout)

	// Determine timezone offset
	_, offset := t.Zone()
	zone := ""
	if strings.HasSuffix(timeStr, "Z") {
		zone = "(+00:00)"
	} else {
		hours := offset / 3600
		minutes := (offset % 3600) / 60
		zone = fmt.Sprintf("(%+02d:%02d)", hours, minutes)
	}
	return formattedTime + " " + zone
}

func RemoveExtraNewLines(text string) string {
	// Regular expression to match two or more consecutive newlines
	re := regexp.MustCompile(`\n{2,}`)
	// Replace matches with a single newline
	return re.ReplaceAllString(text, "\n\n")
}
