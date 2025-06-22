package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type ModbusReading struct {
	RegisterPair string `json:"register_pair"`
	LSR          int16  `json:"LSR"`
	MSR          int16  `json:"MSR"`
}

func readModbusReadingsFromCSV(filepath string) ([]ModbusReading, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV: %w", err)
	}

	// Skip header row
	var readings []ModbusReading
	for i := 1; i < len(records); i++ {
		record := records[i]

		lsr, err := strconv.ParseInt(record[1], 10, 16)
		if err != nil {
			return nil, fmt.Errorf("error parsing LSR at row %d: %w", i, err)
		}

		msr, err := strconv.ParseInt(record[2], 10, 16)
		if err != nil {
			return nil, fmt.Errorf("error parsing MSR at row %d: %w", i, err)
		}

		reading := ModbusReading{
			RegisterPair: record[0],
			LSR:          int16(lsr),
			MSR:          int16(msr),
		}

		readings = append(readings, reading)
	}

	return readings, nil
}

func main() {
	readings, err := readModbusReadingsFromCSV("modbus_reading.csv")
	if err != nil {
		fmt.Printf("Error reading CSV: %v\n", err)
		return
	}

	fmt.Printf("Successfully read %d modbus readings\n", len(readings))

	// Print first few readings as example
	for i := 0; i < 3 && i < len(readings); i++ {
		fmt.Printf("%+v\n", readings[i])
	}
}
