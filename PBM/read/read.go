package Netpbm

import (
	"bufio"
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

	//ouverture du fichier mis en parametre
	path := "./" + filename
	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}

	//lire la premiere ligne
	fileScanner := bufio.NewScanner(file)
	fileScanner.Text()

}
