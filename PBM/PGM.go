package Netpbm

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type PGM struct {
	data          [][]uint8
	width, height int
	magicNumber   string
	max           uint
}

// ReadPGM reads a PGM image from a file and returns a struct that represents the image.
func ReadPGM(filename string) (*PGM, error) {
	pgm := PGM{}

	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Scanner for text-based information
	scanner := bufio.NewScanner(file)

	// Get magic number
	scanner.Scan()
	pgm.magicNumber = scanner.Text()
	if pgm.magicNumber != "P5" && pgm.magicNumber != "P2" {
		return nil, fmt.Errorf("invalid magic number: %s", pgm.magicNumber)
	}

	// Get dimensions
	scanner.Scan()
	dimensions := strings.Fields(scanner.Text())
	if len(dimensions) != 2 {
		return nil, errors.New("bad input format")
	}
	pgm.width, _ = strconv.Atoi(dimensions[0])
	pgm.height, _ = strconv.Atoi(dimensions[1])
	// Check for valid values
	if pgm.width <= 0 || pgm.height <= 0 {
		return nil, fmt.Errorf("invalid size: %d x %d", pgm.width, pgm.height)
	}

	// Get max value
	scanner.Scan()
	maxValue, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return nil, fmt.Errorf("error parsing max value: %v", err)
	}
	pgm.max = uint(maxValue)

	// Move to the beginning of binary data
	for scanner.Scan() {
		if scanner.Text() == "" || strings.HasPrefix(scanner.Text(), "#") {
			continue
		} else {
			break
		}
	}

	if pgm.magicNumber == "P5" {
		// P5 format (raw binary)
		buffer := make([]byte, pgm.width*pgm.height)

		// Keep reading until we have read the expected number of bytes
		totalBytesRead := 0
		for totalBytesRead < pgm.width*pgm.height {
			n, err := file.Read(buffer[totalBytesRead:])
			if err != nil {
				if err == io.EOF {
					break // EOF reached
				} else {
					return nil, fmt.Errorf("error reading binary data: %v", err)
				}
			}
			if n == 0 {
				break // No more data to read
			}
			totalBytesRead += n
		}

		// Check if the total number of bytes read matches the expected size
		if totalBytesRead != pgm.width*pgm.height {
			return nil, fmt.Errorf("unexpected number of bytes read: got %d, expected %d", totalBytesRead, pgm.width*pgm.height)
		}

		// Debug information
		fmt.Printf("Total bytes read: %d\n", totalBytesRead)

		// Populate pgm.data with the binary data
		pgm.data = make([][]uint8, pgm.height)
		for i := 0; i < pgm.height; i++ {
			pgm.data[i] = make([]uint8, pgm.width)
			for j := 0; j < pgm.width; j++ {
				pgm.data[i][j] = buffer[i*pgm.width+j]
			}
		}

	} else if pgm.magicNumber == "P2" {
		// P2 format (ASCII)
		pgm.data = make([][]uint8, pgm.height)
		for i := 0; i < pgm.height; i++ {
			pgm.data[i] = make([]uint8, pgm.width)
			lineValues := strings.Fields(scanner.Text())
			if len(lineValues) != pgm.width {
				return nil, fmt.Errorf("bad row length: %d", len(lineValues))
			}
			for j := 0; j < pgm.width; j++ {
				value, err := strconv.Atoi(lineValues[j])
				if err != nil {
					return nil, fmt.Errorf("error reading pixel value: %v", err)
				}
				pgm.data[i][j] = uint8(value)
			}
			if !scanner.Scan() {
				break // Break if there are no more lines
			}
		}

	} else {
		// Handle other PGM formats here if needed
		return nil, fmt.Errorf("unsupported PGM format: %s", pgm.magicNumber)
	}

	return &pgm, nil
}
