package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	fun "github.com/luraim/fun"
	exif "github.com/rwcarlsen/goexif/exif"
)

type TLIFile struct {
	mtime int64
	fname string
}

func getCtime(fname string) TLIFile {
	fi, _ := os.Stat(fname)
	var thing TLIFile
	thing.fname = fname
	thing.mtime = fi.ModTime().UnixMilli()

	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}

	ex, err := exif.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	tm, err := ex.DateTime()
	if err == nil {
		thing.mtime = tm.UnixMilli()
	}
	f.Close()
	return thing
}

func findfiles(dir string) []string {
	vals, _ := filepath.Glob(dir + "/*.jpg")
	return vals
}

func main() {
	indir := os.Args[1]
	framestr := os.Args[2]
	outdir := os.Args[3]
	frames, _ := strconv.Atoi(framestr)
	dirlist := findfiles(indir)
	fmt.Printf("Would generate %d files into %s\n", frames, outdir)
	fmt.Println("Found files: ", dirlist)
	moredata := fun.Map(dirlist, getCtime)
	fmt.Println("Sorting data...")
	sort.Slice(moredata, func(i int, j int) bool {
		return moredata[i].mtime < moredata[j].mtime
	})
	fmt.Println(moredata)
}
