package all_scans

import (
	"fmt"
	"os"
	"path/filepath"
	"scanner/don_patt"
	"scanner/ec_patt"
	"scanner/schemas"
)

func scanList() int {
	fmt.Println("Available Scans:")
	fmt.Println("1. Recent Donchian 5 Breakouts on Weekly")

	var choice int
	fmt.Print("Enter the number of the scan you want to perform (or 0 to exit): ")
	_, err := fmt.Scan(&choice)
	if err != nil {
		fmt.Println("Invalid input. Please enter a number.")
		return scanList()
	}

	return choice
}

func ControlScans(dailyData, weeklyData schemas.AllSymData, dailyIndicators, weeklyIndicators schemas.AllIndicatorsData) {

	extraSyminKeys := func(myDict map[string]string) []string {
		var keys []string
		for key := range myDict {
			keys = append(keys, key)
		}
		return keys
	}

	ec_filter := func(given_symbols []string, ec_symbols []string) []string {
		var filtered []string

		for _, sym := range given_symbols {
			for _, ec_sym := range ec_symbols {
				if sym == ec_sym {
					filtered = append(filtered, sym)
					break
				}
			}
		}
		return filtered
	}

	// get ec symbols
	long_ec, short_ec := ec_patt.GetEc(dailyData, dailyIndicators)

	//scan
	long_ec_keys := extraSyminKeys(long_ec)
	short_ec_keys := extraSyminKeys(short_ec)

	//scans based on user choice
	for {
		choice := scanList()

		switch choice {
		case 0:
			fmt.Println("Exiting scans.")
			return
		case 1:
			long_don5_bo_sym, short_don5_bo_sym := don_patt.Recent_don5_bo(weeklyData, weeklyIndicators)

			//apply ec filter
			fil_long_don5_bo_sym := ec_filter(long_don5_bo_sym, long_ec_keys)
			fil_short_don5_bo_sym := ec_filter(short_don5_bo_sym, short_ec_keys)

			fmt.Println("Recent Weekly Donchian 5 Breakouts")
			fmt.Println("Long Symbols:", len(fil_long_don5_bo_sym), fil_long_don5_bo_sym)
			fmt.Println("Short Symbols:", len(fil_short_don5_bo_sym), fil_short_don5_bo_sym)

			//save it to TLS file
			writeSymbolsToTLS(fil_long_don5_bo_sym, "don5_bo_l.tls", "ow")
			writeSymbolsToTLS(fil_short_don5_bo_sym, "don5_bo_s.tls", "ow")
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}

}

// writeSymbolsToTLS writes only the symbol values to a TLS file in the specified directory
func writeSymbolsToTLS(symbols []string, filename string, where string) error {

	// Define the target directory
	var scanDir string
	if where == "ow" {
		scanDir = `C:\Users\madar\OneDrive\Desktop\scan_results\weekly\once_a_week`
	}

	// Create the full file path
	fullPath := filepath.Join(scanDir, filename)

	// Create the file
	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// Write symbols to file
	for _, value := range symbols {
		_, err := file.WriteString(value + "\n")
		if err != nil {
			return fmt.Errorf("error writing to file: %v", err)
		}
	}

	return nil
}
