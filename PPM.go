package Netpbm

import (
	"bufio"
	"fmt"
	"io"
	"os"
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
func ReadPPM(fileName string) (*PPM, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	// Get the magic number
	magicNumber, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading the magic number: %v", err)
	}
	magicNumber = strings.TrimSpace(magicNumber)

	// Get dimensions
	sepa, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error with dimensions: %v", err)
	}
	var width, height int
	_, err = fmt.Sscanf(strings.TrimSpace(sepa), "%d %d", &width, &height)
	if err != nil {
		return nil, fmt.Errorf("invalid dimensions: %v", err)
	}

	// Get the maximum value
	maxValue, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading the maximum value: %v", err)
	}
	maxValue = strings.TrimSpace(maxValue)
	var max uint8
	_, err = fmt.Sscanf(maxValue, "%d", &max)
	if err != nil {
		return nil, fmt.Errorf("invalid maximum value: %v", err)
	}

	// Get image data
	data := make([][]Pixel, height)
	expectedBytesPerPixel := 3

	if magicNumber == "P3" {
		//  The P3 format (ASCII)
		for y := 0; y < height; y++ {
			line, err := reader.ReadString('\n')
			if err != nil {
				return nil, fmt.Errorf("error reading data at line %d: %v", y, err)
			}
			fields := strings.Fields(line)
			rowdata := make([]Pixel, width)
			for x := 0; x < width; x++ {
				if x*3+2 >= len(fields) {
					return nil, fmt.Errorf("index out of range at line %d, column %d", y, x)
				}
				var pixel Pixel
				_, err := fmt.Sscanf(fields[x*3], "%d", &pixel.R)
				if err != nil {
					return nil, fmt.Errorf("error parsing the red value at line %d, column %d: %v", y, x, err)
				}
				_, err = fmt.Sscanf(fields[x*3+1], "%d", &pixel.G)
				if err != nil {
					return nil, fmt.Errorf("error parsing the green value at line %d, column %d: %v", y, x, err)
				}
				_, err = fmt.Sscanf(fields[x*3+2], "%d", &pixel.B)
				if err != nil {
					return nil, fmt.Errorf("error parsing the blue value at line %d, column %d: %v", y, x, err)
				}
				rowdata[x] = pixel
			}
			data[y] = rowdata
		}
	} else if magicNumber == "P6" {
		// The P6 format (binary)
		for y := 0; y < height; y++ {
			row := make([]byte, width*expectedBytesPerPixel)
			n, err := reader.Read(row)
			if err != nil {
				if err == io.EOF {
					return nil, fmt.Errorf("unexpected end of file at line %d", y)
				}
				return nil, fmt.Errorf("error reading pixel data at line %d: %v", y, err)
			}
			if n < width*expectedBytesPerPixel {
				return nil, fmt.Errorf("unexpected end of file at line %d, expected %d bytes, got %d", y, width*expectedBytesPerPixel, n)
			}
			rowdata := make([]Pixel, width)
			for x := 0; x < width; x++ {
				pixel := Pixel{R: row[x*expectedBytesPerPixel], G: row[x*expectedBytesPerPixel+1], B: row[x*expectedBytesPerPixel+2]}
				rowdata[x] = pixel
			}
			data[y] = rowdata
		}
	}

	return &PPM{data: data, width: width, height: height, magicNumber: magicNumber, max: max}, nil
}

// Size returns the width and height of the image.
func (ppm *PPM) Size() (int, int) {
	return ppm.height, ppm.width
}

// At returns the value of the pixel at (x, y).
func (ppm *PPM) At(x, y int) Pixel {
	// Check if the coordinates are within the valid range
	if x < 0 || x >= ppm.width || y < 0 || y >= ppm.height {
		// Return a default Pixel if the coordinates are out of range
		return Pixel{}
	}

	return ppm.data[y][x]
}

// Set sets the value of the pixel at (x, y).
func (ppm *PPM) Set(x, y int, value Pixel) {
	ppm.data[x][y] = value
}

// Save saves the PPM image to a file and returns an error if there was a problem.
func (ppm *PPM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Write the PPM header
	_, err = fmt.Fprintf(writer, "%s\n%d %d\n%d\n", ppm.magicNumber, ppm.width, ppm.height, ppm.max)
	if err != nil {
		return fmt.Errorf("error writing PPM header: %v", err)
	}

	// Write image data based on the magic number
	if ppm.magicNumber == "P3" {
		// P3 format (ASCII)
		for y := 0; y < ppm.height; y++ {
			for x := 0; x < ppm.width; x++ {
				_, err := fmt.Fprintf(writer, "%d %d %d ", ppm.data[y][x].R, ppm.data[y][x].G, ppm.data[y][x].B)
				if err != nil {
					return fmt.Errorf("error writing pixel data: %v", err)
				}
			}
			_, err := fmt.Fprint(writer, "\n")
			if err != nil {
				return fmt.Errorf("error writing newline character: %v", err)
			}
		}
	} else if ppm.magicNumber == "P6" {
		// P6 format (binary)
		for y := 0; y < ppm.height; y++ {
			for x := 0; x < ppm.width; x++ {
				_, err := writer.Write([]byte{ppm.data[y][x].R, ppm.data[y][x].G, ppm.data[y][x].B})
				if err != nil {
					return fmt.Errorf("error writing pixel data: %v", err)
				}
			}
		}
	}

	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("error flushing writer: %v", err)
	}

	return nil
}

// Invert inverts the colors of the PPM image.
func (ppm *PPM) Invert() {
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			// Invert the Red, Green, and Blue values of each pixel
			ppm.data[y][x].R = ppm.max - ppm.data[y][x].R
			ppm.data[y][x].G = ppm.max - ppm.data[y][x].G
			ppm.data[y][x].B = ppm.max - ppm.data[y][x].B
		}
	}
}

// Flip flips the PPM image horizontally.
func (ppm *PPM) Flip() {
	// Flip each row horizontally
	for i := 0; i < ppm.height; i++ {
		for j, k := 0, ppm.width-1; j < k; j, k = j+1, k-1 {
			// Swap pixel values
			ppm.data[i][j], ppm.data[i][k] = ppm.data[i][k], ppm.data[i][j]
		}
	}
}

// Flop flops the PPM image vertically.
func (ppm *PPM) Flop() {
	// Flop the rows vertically
	for i, j := 0, ppm.height-1; i < j; i, j = i+1, j-1 {
		// Swap rows
		ppm.data[i], ppm.data[j] = ppm.data[j], ppm.data[i]
	}
}

// SetMagicNumber sets the magic number of the PPM image.
func (ppm *PPM) SetMagicNumber(magicNumber string) {
	ppm.magicNumber = magicNumber
}

// SetMaxValue sets the max value of the PPM image.
func (ppm *PPM) SetMaxValue(maxValue uint8) {
	// Scale the color values of each pixel
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			// Scale each color component based on the new max value
			ppm.data[y][x].R = uint8(float64(ppm.data[y][x].R) * float64(maxValue) / float64(ppm.max))
			ppm.data[y][x].G = uint8(float64(ppm.data[y][x].G) * float64(maxValue) / float64(ppm.max))
			ppm.data[y][x].B = uint8(float64(ppm.data[y][x].B) * float64(maxValue) / float64(ppm.max))
		}
	}
	ppm.max = uint8(maxValue)
}

// Rotate90CW rotates the PPM image 90° clockwise.
func (ppm *PPM) Rotate90CW(){
    // Create a new PPM image with swapped dimensions
    newWidth, newHeight := ppm.height, ppm.width
    rotatedImage := &PPM{
        data:        make([][]Pixel, newHeight),
        width:       newWidth,
        height:      newHeight,
        magicNumber: ppm.magicNumber,
        max:         ppm.max,
    }

    // Initialize the data for the new PPM image
    for i := 0; i < newHeight; i++ {
        rotatedImage.data[i] = make([]Pixel, newWidth)
    }

    // Copy pixels in a rotated manner
    for y := 0; y < ppm.height; y++ {
        for x := 0; x < ppm.width; x++ {
            rotatedImage.data[x][ppm.height-y-1] = ppm.data[y][x]
        }
    }

    // Update the original PPM image with the rotated one
    ppm.data = rotatedImage.data
    ppm.width = rotatedImage.width
    ppm.height = rotatedImage.height
}

// ToPGM converts the PPM image to PGM.
func (ppm *PPM) ToPGM() *PGM{
	// Create a new PGM image with the same dimensions
    pgm := &PGM{
		width:       ppm.width,
		height:      ppm.height,
		magicNumber: "P2",
		max:         uint8(ppm.max),
	}

	    // Initialize the data for the new PGM image

	pgm.data = make([][]uint8, ppm.height)
	for i := range pgm.data {
		pgm.data[i] = make([]uint8, ppm.width)
	}

    // Convert RGB to grayscale and copy the pixel values
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			// Convert RGB to grayscale
			gray := uint8((int(ppm.data[y][x].R) + int(ppm.data[y][x].G) + int(ppm.data[y][x].B)) / 3)
			pgm.data[y][x] = gray
		}
	}

	return pgm
}

// ToPBM converts the PPM image to PBM.
func (ppm *PPM) ToPBM() *PBM {
    // Create a new PBM image with the same dimensions
    pbm := &PBM{
        width:       ppm.width,
        height:      ppm.height,
        magicNumber: "P1",
    }

    // Initialize the data for the new PBM image
    pbm.data = make([][]bool, ppm.height)
    for i := range pbm.data {
        pbm.data[i] = make([]bool, ppm.width)
    }

    // Define a threshold for converting color to monochrome
    threshold := uint8(ppm.max / 2)

    // Convert each pixel to monochrome based on average intensity
    for y := 0; y < ppm.height; y++ {
        for x := 0; x < ppm.width; x++ {
            // Calculate the average intensity using RGB values
            average := (uint16(ppm.data[y][x].R) + uint16(ppm.data[y][x].G) + uint16(ppm.data[y][x].B)) / 3
            pbm.data[y][x] = average < uint16(threshold)
        }
    }

    return pbm
}

type Point struct{
    X, Y int
}
