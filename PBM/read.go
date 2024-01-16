package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

// ReadPBM lie l'image PBM du fichier et retourne dans la struct avec les infos de l'image.
func ReadPBM(filename string) (*PBM, error) {
	PBMstock := PBM{}

	// Ouverture du fichier
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	// Scanner le texte + variables
	scanner := bufio.NewScanner(file)
	reader := bufio.NewReader(file)
	confirmOne := false
	confirmTwo := false
	Line := 0

	// Ignorer le #
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "#") {
			continue
		} else if !confirmOne { // Prendre en compte le P1 ou P4
			PBMstock.magicNumber = scanner.Text()
			confirmOne = true
		} else if !confirmTwo { // Prendre en compte la taille de l'image
			sepa := strings.Fields(scanner.Text())
			if len(sepa) > 0 {
				PBMstock.width, _ = strconv.Atoi(sepa[0])
				PBMstock.height, _ = strconv.Atoi(sepa[1])
			}
			confirmTwo = true

			PBMstock.data = make(([][]bool), PBMstock.height) // Création de la matrice pour le data en prenant en compte la taille
			for i := range PBMstock.data {
				PBMstock.data[i] = make(([]bool), PBMstock.width)
			}

		} else {

			if PBMstock.magicNumber == "P1" { // Pour le P1 rendre les 0 en false et les 1 en true
				test := strings.Fields(scanner.Text())
				for i := 0; i < PBMstock.width; i++ {
					if test[i] == "1" {
						PBMstock.data[Line][i] = true
					} else {
						PBMstock.data[Line][i] = false
					}
				}
				Line++
			}else if PBMstock.magicNumber == "P4" {
				// Read P4 format (binary)
				expectedBytesPerRow := (PBMstock.width + 7) / 8
				for y := 0; y < PBMstock.height; y++ {
					row := make([]byte, expectedBytesPerRow)
					n, err := reader.Read(row)
					if err != nil {
						if err == io.EOF {
							return nil, fmt.Errorf("unexpected end of file at row %d", y)
						}
						return nil, fmt.Errorf("error reading pixel data at row %d: %v", y, err)
					}
					if n < expectedBytesPerRow {
						return nil, fmt.Errorf("unexpected end of file at row %d, expected %d bytes, got %d", y, expectedBytesPerRow, n)
					}
		
					for x := 0; x < PBMstock.width; x++ {
						byteIndex := x / 8
						bitIndex := 7 - (x % 8)
		
						// Convert ASCII to decimal and extract the bit
						decimalValue := int(row[byteIndex])
						bitValue := (decimalValue >> bitIndex) & 1
		
						PBMstock.data[y][x] = bitValue != 0
					}
				}
			}
		}
	}

	fmt.Printf("%v\n", PBMstock)
	return &PBM{PBMstock.data, PBMstock.height, PBMstock.width, PBMstock.magicNumber}, err // Retourner la struct en entiere
}

// Size retourne la hauteur et la largeur de l'image.
func (pbm *PBM) Size() (int, int) {
	return pbm.height, pbm.width
}

// At retourne la valeur de chaque pixel en (x, y).
func (pbm *PBM) At(x, y int) bool {
	return pbm.data[x][y]
}

// Set défini la valeur de chaque pixel à (x, y).
func (pbm *PBM) Set(x, y int, value bool) {
	pbm.data[x][y] = value
}

// Save saves the PBM image to a file and returns an error if there was a problem.
func (pbm *PBM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// Ecrire le magique number et la taille de l'image dans le nouveau fichier
	_, err = fmt.Fprintf(file, "%s\n%d %d\n", pbm.magicNumber, pbm.width, pbm.height)
	if err != nil {
		return fmt.Errorf("error writing magic number and dimensions: %v", err)
	}

	// Entrer les donnees de l'image
	if pbm.magicNumber == "P1" { 												//Pour le P1 
		for _, row := range pbm.data {
			for _, pixel := range row {
				if pixel {
					_, err = file.WriteString("1 ")
				} else {
					_, err = file.WriteString("0 ")
				}
				if err != nil {
					return fmt.Errorf("error writing pixel data: %v", err)
				}
			}
			_, err = file.WriteString("\n")
			if err != nil {
				return fmt.Errorf("error writing pixel data: %v", err)
			}
		}
	} else if pbm.magicNumber == "P4" { 										// Pour le P4
		for _, row := range pbm.data {
			for x := 0; x < pbm.width; x += 8 {
				var byteValue byte
				for i := 0; i < 8 && x+i < pbm.width; i++ {
					bitIndex := 7 - i
					if row[x+i] {
						byteValue |= 1 << bitIndex
					}
				}
				_, err = file.Write([]byte{byteValue})
				if err != nil {
					return fmt.Errorf("error writing pixel data: %v", err)
				}
			}
		}
	}

	return nil
}
