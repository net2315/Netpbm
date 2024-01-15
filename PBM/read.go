package main

import (
	"bufio"
	"fmt"
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


