package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kumakichi/merge-images"
)

var (
	direction string
	outfile   string
)

func init() {
	flag.StringVar(&direction, "t", "", "Specify merge direction: h(horizontal) or v(vertical)")
	flag.StringVar(&outfile, "o", "mergedImage", "Name of merged file")
}

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("Usage: %s [options] inFile1 [inFile2] ...\n\n", os.Args[0])
		flag.PrintDefaults()
		return
	}

	flag.Parse()

	if direction != "v" && direction != "h" {
		fmt.Println("Invalid merge direction, support only [h/v]")
		os.Exit(-1)
	}

	args := flag.Args()

	var err error
	if direction == "v" {
		err = merge_images.MergeImage(merge_images.VERTICLE, outfile, args...)
	} else {
		err = merge_images.MergeImage(merge_images.HORIZONTAL, outfile, args...)
	}

	if err != nil {
		log.Fatal(err)
	}
}
