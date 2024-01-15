package Netpbm

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

func main()  {
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
	input := scanner.Text()
	confirmOne := false
	confirmTwo := false

	for scanner.Scan() {
		if strings.HasPrefix(input, "#") {
			continue
		} else if !confirmOne {
			PBMstock.magicNumber = input
			confirmOne = true
		} else if confirmTwo {
			sepa := strings.Fields(input)
			if len(sepa)>2 {
			width, errW := strconv.Atoi(sepa[0])
			height, errH := strconv.Atoi(sepa[1])
			if errW != nil || errH != nil {
				return nil, fmt.Errorf("Erreur de conversion de largeur/hauteur: %v, %v", errW, errH)
			}
			PBMstock.width = width
			PBMstock.height = height
			confirmTwo = true
			}
		} 
		if PBMstock.magicNumber == "P1" {
			for i := 0; i < PBMstock.width; i++ {
				for j := 1; j < PBMstock.height; j++ {
					if PBMstock.data[i][j] == true {
						fmt.Println("true")
					} else {
						fmt.Println("false")
					}
				}
			}
		}
	}
	fmt.Printf("%v\n", PBMstock)
	return &PBM{PBMstock.data,PBMstock.height,PBMstock.width,PBMstock.magicNumber}, err
}
