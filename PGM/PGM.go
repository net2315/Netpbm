package Netpbm

type PGM struct{
    data [][]uint8
    width, height int
    magicNumber string
    max uint
}

// ReadPGM reads a PGM image from a file and returns a struct that represents the image.
func ReadPGM(filename string) (*PGM, error){
    // ...
}
