package main

import (
    "fmt"
    "os"
	"time"
	"io"
)

type EtagLine struct {
	Etag      string
	Timestamp time.Time
	Size      int64
	Offset    int64
	Comp      string
}



func ParseJpeg(localJpeg string, baseJpeg EtagLine, requestedJpeg EtagLine) error {
	file, err := os.Open(localJpeg)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()


	// offset := baseJpeg.Offset

	buffer := make([]byte, baseJpeg.Size) // Buffer the size of jpeg

	// Seek to the specified offset
	// We don't care about first 4kb
	_, err = file.Seek(baseJpeg.Offset + 4096, io.SeekStart)
	if err != nil {
		return fmt.Errorf("failed to seek in file: %v", err)
	}

	_, err = file.Read(buffer)
	if err != nil {
		fmt.Println("failed to reader buffer")
	}

	// for _, b := range buffer {
	// 	fmt.Printf("%02x ", b)
	// }

	// fmt.Println(buffer)

	// Open or create a new JPEG file to write the contents
	newFile, err := os.OpenFile("output.jpeg", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("failed to open/create jpeg file: %v\n", err)
	}
	defer newFile.Close()

	// Write the read contents to the new JPEG file
	_, err = newFile.Write(buffer)
	if err != nil {
		fmt.Printf("failed to write to jpeg file: %v\n", err)
	}

	fmt.Println("Content successfully written to output.jpeg")

	return nil
}


func main(){
	var err error 
	masterJpeg := "20240813210000.000_pre.jpeg"
	ts := "2024-08-13T21:59:48.744"
	timeFormat := "2006-01-02T15:04:05.000"

	requestedJpeg := EtagLine{ Etag: "PRFR",  Offset: 94058848, Size: 21994, Comp: "0x100"}
	requestedJpeg.Timestamp, err = time.Parse(timeFormat, ts)
	if err != nil {
		fmt.Println("Error Parsing timestamp")
	}

	baseJpeg := EtagLine{ Etag: "PRFR",  Offset: 93960505, Size: 98335, Comp: "0x200"}

	baseJpeg.Timestamp, err = time.Parse(timeFormat, "2024-08-13T21:59:46.016")
	if err != nil {
		fmt.Println("FAiled to parse timestamp")
	}

	// fmt.Printf("Args: %+v %+v %s", baseJpeg, requestedJpeg, masterJpeg)

	// Analyse the diffs between the jpegs
	err = ParseJpeg(masterJpeg, baseJpeg, requestedJpeg)
	if err != nil {
		fmt.Println(err)
	}
}