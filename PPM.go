package Netpbm

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PPM struct {
	data          [][]Pixel
	width, height int
	magicNumber   string
	max           uint8
}

type Pixel struct {
	R, G, B uint8
}

// ReadPPM reads a PPM image from a file and returns a struct that represents the image.
func ReadPPM(filename string) (*PPM, error) {
	ppm := PPM{}
	pixel := Pixel{}

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
	ppm.magicNumber = scanner.Text()
	if ppm.magicNumber != "P3" {
		return nil, fmt.Errorf("unsupported PPM format: %s", ppm.magicNumber)
	}

	// Get dimensions
	scanner.Scan()
	sepa := strings.Fields(scanner.Text())
	if len(sepa) != 2 {
		return nil, errors.New("bad input format")
	}
	ppm.width, _ = strconv.Atoi(sepa[0])
	ppm.height, _ = strconv.Atoi(sepa[1])

	// Read and parse the maximum color value
	scanner.Scan()
	maxValue, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return nil, fmt.Errorf("error parsing max value: %v", err)
	}
	ppm.max = uint8(maxValue)

	// Create a 2D slice to store pixel data
	ppm.data = make([][]Pixel, ppm.height)
	for i := range ppm.data {
		ppm.data[i] = make([]Pixel, ppm.width)
	}

	// Read ASCII pixel data for P3 format
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			scanner.Scan() // Add this line to read newlines between lines
			_, err := fmt.Fscanf(file, "%d%d%d", &pixel.R, &pixel.G, &pixel.B)
			if err != nil {
				return nil, err
			}
			ppm.data[i][j] = pixel
		}
		scanner.Scan() // Add this line to read newlines between lines
	}

	return &ppm, nil
}
