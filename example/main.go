package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/kumakichi/merge-images"
)

var (
	orientation string
	outfile     string
)

func init() {
	flag.StringVar(&orientation, "t", "", "Specify merge orientation: h(horizontal) or v(vertical)")
	flag.StringVar(&outfile, "o", "mergedImage", "Name of merged file")
}

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("Usage: %s [options] inFile1 [inFile2] ...\n\n", os.Args[0])
		flag.PrintDefaults()
		return
	}

	flag.Parse()

	if orientation == "" {
		fmt.Printf("Please specify merge orientation\n\n")
		fmt.Printf("Usage: %s [options] inFile1 [inFile2] ...\n\n", os.Args[0])
		flag.PrintDefaults()
		return
	}

	if orientation != "v" && orientation != "h" {
		fmt.Println("Invalid merge orientation, support only [h/v]")
		os.Exit(-1)
	}

	args := flag.Args()

	var err error

	merge_images.SetBackgroundColor(color.White)

	if orientation == "v" {
		err = merge_images.MergeImage(merge_images.VERTICLE, outfile, args...)
	} else {
		err = merge_images.MergeImage(merge_images.HORIZONTAL, outfile, args...)
	}

	if err != nil {
		log.Fatal(err)
	}
}
