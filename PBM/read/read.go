package Netpbm

import (
	"fmt"
	"os"
)

type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

// ReadPBM reads a PBM image from a file and returns a struct that represents the image.
func ReadPBM(filename string) (*PBM, error) {
	path := "./" + filename
	_, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	
}
