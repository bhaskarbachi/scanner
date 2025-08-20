package data_process

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"scanner/schemas"
	"strconv"
	"strings"
	"time"
)

// Column indices - define these constants
const (
	symbolIndex = 0
	dateIndex   = 1
	openIndex   = 2
	highIndex   = 3
	lowIndex    = 4
	closeIndex  = 5
	volumeIndex = 6
)

// Configuration
const minRequiredBars = 30

// dateParser handles date parsing with format detection
type dateParser struct {
	detectedLayout string
	layoutDetected bool
}

// newDateParser creates a new date parser
func newDateParser() *dateParser {
	return &dateParser{}
}

// Parse attempts to parse a date string, using cached format when possible
func (dp *dateParser) Parse(dateStr string, row int) (time.Time, error) {
	cleanDateStr := cleanDateString(dateStr)

	// Try cached layout first
	if dp.layoutDetected {
		if date, err := time.Parse(dp.detectedLayout, cleanDateStr); err == nil {
			return date, nil
		}
		// Reset if cached layout fails
		dp.layoutDetected = false
	}

	// Try different layouts
	layouts := []string{
		"2006-01-02",          // ISO date format (most common)
		"01/02/2006",          // US date format
		"02/01/2006",          // UK date format
		"2006-01-02 15:04:05", // Date with time
		"01-02-2006",          // US date with dashes
		"02-01-2006",          // UK date with dashes
		"2006/01/02",          // Alternative ISO format
	}

	for _, layout := range layouts {
		if date, err := time.Parse(layout, cleanDateStr); err == nil {
			dp.detectedLayout = layout
			dp.layoutDetected = true
			return date, nil
		}
	}

	return time.Time{}, fmt.Errorf("error parsing date in row %d: unable to parse date '%s'", row, dateStr)
}

func getDailyData(filePath string) (schemas.AllSymData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening csv file: %w", err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)

	// Read and validate header
	header, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading header: %w", err)
	}

	if len(header) < 7 {
		return nil, fmt.Errorf("invalid header: expected at least 7 columns, got %d", len(header))
	}

	totalData := make(schemas.AllSymData)
	dateParser := newDateParser()

	// Process records
	for rowNum := 2; ; rowNum++ { // Start from 2 (header is row 1)
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading record at row %d: %w", rowNum, err)
		}

		if len(record) < 7 {
			fmt.Printf("Skipping incomplete record at row %d (has %d columns)\n", rowNum, len(record))
			continue
		}

		// Parse all fields
		symbol := strings.TrimSpace(record[symbolIndex])
		if symbol == "" {
			fmt.Printf("Skipping row %d: empty symbol\n", rowNum)
			continue
		}

		date, err := dateParser.Parse(record[dateIndex], rowNum)
		if err != nil {
			return nil, err
		}

		open, err := parseFloat(record[openIndex], rowNum, "open")
		if err != nil {
			return nil, err
		}

		high, err := parseFloat(record[highIndex], rowNum, "high")
		if err != nil {
			return nil, err
		}

		low, err := parseFloat(record[lowIndex], rowNum, "low")
		if err != nil {
			return nil, err
		}

		close, err := parseFloat(record[closeIndex], rowNum, "close")
		if err != nil {
			return nil, err
		}

		volume, err := parseInt(record[volumeIndex], rowNum, "volume")
		if err != nil {
			return nil, err
		}

		// Validate OHLC data
		if high < low || open < 0 || high < 0 || low < 0 || close < 0 || volume < 0 {
			fmt.Printf("Warning: invalid OHLC data at row %d for symbol %s\n", rowNum, symbol)
		}

		data := schemas.TOHLCV{
			Date:   date,
			Open:   open,
			High:   high,
			Low:    low,
			Close:  close,
			Volume: volume,
		}

		totalData[symbol] = append(totalData[symbol], data)
	}

	// Filter symbols with insufficient data
	filteredCount := 0
	for symbol, data := range totalData {
		if len(data) < minRequiredBars {
			delete(totalData, symbol)
			filteredCount++
		}
	}

	if filteredCount > 0 {
		fmt.Printf("Filtered out %d symbols with less than %d bars\n", filteredCount, minRequiredBars)
	}

	if len(totalData) == 0 {
		return nil, fmt.Errorf("no valid symbols found with sufficient data")
	}

	return totalData, nil
}

// Optimized cleanDateString with fewer string operations
func cleanDateString(dateStr string) string {
	dateStr = strings.TrimSpace(dateStr)

	if len(dateStr) == 0 {
		return dateStr
	}

	// Find first occurrence of 'T' or space
	for i, char := range dateStr {
		if char == 'T' || char == ' ' {
			return dateStr[:i]
		}
	}

	return dateStr
}

func parseFloat(valueStr string, row int, fieldName string) (float64, error) {
	valueStr = strings.TrimSpace(valueStr)
	if valueStr == "" {
		return 0, fmt.Errorf("empty %s value in row %d", fieldName, row)
	}

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing %s in row %d: %w", fieldName, row, err)
	}
	return value, nil
}

func parseInt(valueStr string, row int, fieldName string) (int, error) {
	valueStr = strings.TrimSpace(valueStr)
	if valueStr == "" {
		return 0, fmt.Errorf("empty %s value in row %d", fieldName, row)
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("error parsing %s in row %d: %w", fieldName, row, err)
	}
	return value, nil
}
