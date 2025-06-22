package main

import (
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"
)

type ModbusReading struct {
	RegisterPair string  `json:"register_pair"`
	LSR          int16   `json:"LSR"`
	MSR          int16   `json:"MSR"`
	Value        float32 `json:"value"`
}

func parseModbusReadingsFromCSV(filepath string) ([]ModbusReading, error) {
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
			Value:        decodeFloat32(int16(lsr), int16(msr)),
		}

		readings = append(readings, reading)
	}

	return readings, nil
}

func decodeFloat32(lsr int16, msr int16) float32 {
	// Use uint16 to pack bytes
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint16(buf[0:], uint16(lsr))
	binary.LittleEndian.PutUint16(buf[2:], uint16(msr))
	bits := binary.LittleEndian.Uint32(buf)
	return math.Float32frombits(bits)
}

func main() {
	readings, err := parseModbusReadingsFromCSV("modbus_reading.csv")
	if err != nil {
		fmt.Printf("Error reading CSV: %v\n", err)
		return
	}

	fmt.Printf("Successfully parsed %d modbus readings\n", len(readings))

	//print the register pair and value
	for _, reading := range readings {
		fmt.Printf("%s: %.2f\n", reading.RegisterPair, reading.Value)
	}
}
