package file

import (
	"fmt"
	"log"
	"os"
)

/*
Create is a simple helper function to create a file at path
Also checks if the file at path already exists
*/
func Create(path string) *os.File {
	// create file if not exists
	var file, err = os.Create("ROM")

	if isError(err) {
		os.Exit(65)
	}
	return file
}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}

/*
WriteBytes writes bytes to addr in the given file
*/
func WriteBytes(file *os.File, bytes []byte, particularOffset bool, addr int64) {
	var jmpFileLoc int64
	if particularOffset {
		originalOffset, _ := file.Seek(0, 1)
		jmpFileLoc = originalOffset
		file.Seek(addr, 0)
	}
	bytesWritten, err := file.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Wrote %d bytes.\n", bytesWritten)
	if particularOffset {
		file.Seek(jmpFileLoc, 0)
	}
}
