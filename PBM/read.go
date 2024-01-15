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

func main() {
	ReadPBM("test1.pbm")
}

// ReadPBM lie l'image PBM du fichier et retourne dans la struct les infos de l'image.
func ReadPBM(filename string) (*PBM, error) {
	PBMstock := PBM{}
	//ouverture du fichier mis en parametre
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	//Scanner le texte + conditions
	scanner := bufio.NewScanner(file)
	confirmOne := false
	confirmTwo := false
	Line := 0

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "#") {
			continue
		} else if !confirmOne {
			PBMstock.magicNumber = scanner.Text()
			confirmOne = true
		} else if !confirmTwo {
			sepa := strings.Fields(scanner.Text())
			if len(sepa) > 0 {
				PBMstock.width, _ = strconv.Atoi(sepa[0])
				PBMstock.height, _ = strconv.Atoi(sepa[1])
			}
			confirmTwo = true

			PBMstock.data = make(([][] bool), PBMstock.height)
			for i := range PBMstock.data {
				PBMstock.data[i] = make(([]bool), PBMstock.width)
			}

		} else {

			if PBMstock.magicNumber == "P1" {
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
	return &PBM{PBMstock.data, PBMstock.height, PBMstock.width, PBMstock.magicNumber}, err
}
