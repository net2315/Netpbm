package Netpbm

import (
	"bufio"
	"errors"
	"fmt"
	// "io"
	"os"
	"strconv"
	"strings"
)

type PGM struct {
	data          [][]uint8
	width, height int
	magicNumber   string
	max           uint8
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
    sepa := strings.Fields(scanner.Text())
    if len(sepa) != 2 {
        return nil, errors.New("bad input format")
    }
    pgm.width, _ = strconv.Atoi(sepa[0])
    pgm.height, _ = strconv.Atoi(sepa[1])
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
    pgm.max = uint8(maxValue)

    // Move to the beginning of binary data
    for scanner.Scan() {
        if scanner.Text() == "" || strings.HasPrefix(scanner.Text(), "#") {
            continue
        } else {
            break
        }
    }

    // P5 format (raw binary)
    if pgm.magicNumber == "P5" {
    //     buffer := make([]byte, pgm.width*pgm.height)

    //     // Read the binary data directly into the buffer
    //     _, err := file.Read(buffer)
    //     if err != nil {
    //         return nil, fmt.Errorf("error reading binary data: %v", err)
    //     }

    //     // Populate pgm.data with the binary data
    //     pgm.data = make([][]uint8, pgm.height)
    //     for i := 0; i < pgm.height; i++ {
    //         pgm.data[i] = make([]uint8, pgm.width)
    //         for j := 0; j < pgm.width; j++ {
    //             pgm.data[i][j] = buffer[i*pgm.width+j]
    //         }
    //     }
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

// Size returns the width and height of the image.
func (pgm *PGM) Size() (int, int) {
	return pgm.height, pgm.width
}

// At returns the value of the pixel at (x, y).
func (pgm *PGM) At(x, y int) uint8 {
	return pgm.data[x][y]
}

// Set sets the value of the pixel at (x, y).
func (pgm *PGM) Set(x, y int, value uint8) {
	pgm.data[x][y] = value
}

// Save saves the PGM image to a file and returns an error if there was a problem.
func (pgm *PGM) Save(filename string) error {
	// Open the file for writing
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// Write magic number and sepa to the file
	fmt.Fprintf(file, "%s\n%d %d\n%d\n", pgm.magicNumber, pgm.width, pgm.height, pgm.max)

	// Write pixel values
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			fmt.Fprintf(file, "%d ", pgm.data[i][j])
		}
		fmt.Fprintln(file) // Newline after each row
	}

	return nil
}

// Invert inverts the colors of the PGM image.
func (pgm *PGM) Invert() {
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			pgm.data[i][j] = uint8(pgm.max) - pgm.data[i][j]
		}
	}
}

// Flip flips the PGM image horizontally.
func (pgm *PGM) Flip() {
	// Flip each row horizontally
	for i := 0; i < pgm.height; i++ {
		for j, k := 0, pgm.width-1; j < k; j, k = j+1, k-1 {
			// Swap pixel values
			pgm.data[i][j], pgm.data[i][k] = pgm.data[i][k], pgm.data[i][j]
		}
	}
}

// Flop flops the PGM image vertically.
func (pgm *PGM) Flop(){
     // Flop the rows vertically
	 for i, j := 0, pgm.height-1; i < j; i, j = i+1, j-1 {
        // Swap rows
        pgm.data[i], pgm.data[j] = pgm.data[j], pgm.data[i]
    }
}

// SetMagicNumber sets the magic number of the PGM image.
func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber
}

// SetMaxValue sets the max value of the PGM image.
func (pgm *PGM) SetMaxValue(maxValue uint8){
	for y := range pgm.data {
        for x := range pgm.data[y] {
            prevvalue := pgm.data[y][x];
            newvalue := prevvalue*uint8(5)/pgm.max
            pgm.data[y][x] = newvalue;
        }
    }
    pgm.max = maxValue;
}

// Rotate90CW rotates the PGM image 90° clockwise.
func (pgm *PGM) Rotate90CW(){
    // Create a new PGM image with swapped width and height
    rotatedPgm := &PGM{
        width:  pgm.height,
        height: pgm.width,
        max:    pgm.max,
        magicNumber: pgm.magicNumber,
        data:   make([][]uint8, pgm.width),
    }

    // Initialize the data array for the rotated image
    for i := 0; i < rotatedPgm.width; i++ {
        rotatedPgm.data[i] = make([]uint8, rotatedPgm.height)
    }

    // Populate the rotated image data
    for i := 0; i < pgm.height; i++ {
        for j := 0; j < pgm.width; j++ {
            rotatedPgm.data[j][pgm.height-1-i] = pgm.data[i][j]
        }
    }

    // Update the original image with the rotated values
    pgm.width, pgm.height = rotatedPgm.width, rotatedPgm.height
    pgm.data = rotatedPgm.data
}

// ToPBM converts the PGM image to PBM.
func (pgm *PGM) ToPBM() *PBM{
    pbmData := make([][]bool, pgm.height)
	for y := 0; y < pgm.height; y++ {
		pbmData[y] = make([]bool, pgm.width)
		for x := 0; x < pgm.width; x++ {
			pbmData[y][x] = pgm.data[y][x] < uint8(pgm.max/2)
		}
	}

	return &PBM{
		data:        pbmData,
		width:       pgm.width,
		height:      pgm.height,
		magicNumber: "P1",
	}
}
