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

// tuple pair linking files and their mtime
type TLIFile struct {
	mtime int64
	fname string
}

// tuple triple representing a blend task alpha proportion between files a and b
type TLItask struct {
	afile string
	bfile string
	alpha float64
}

// return image file's creation date (from EXIF if possible, else fs mtime)
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
		log.Println(err)
	}
	tm, err := ex.DateTime()
	if err == nil {
		thing.mtime = tm.UnixMilli()
	}
	f.Close()
	return thing
}

// Convert list of files to list of tasks
func files2tasks(data []TLIFile, noframes int64) []TLItask {
	var ret []TLItask
	var i int64
	var task TLItask

	nopoints := int64(len(data))
	mint := float64(data[0].mtime)
	maxt := float64(data[nopoints-1].mtime)

	for i = 0; i <= noframes-1; i++ {
		tstmp := mint + (maxt-mint)*float64(i)/float64(noframes-1)
		var n int64
		var m int64

		alpha := float64(0)

		for n = 0; n < nopoints-1 && float64(data[n].mtime) < tstmp; n++ {
		}
		n--
		if n < 0 {
			n = 0
		}
		for m = nopoints - 1; m > n+1 && float64(data[m].mtime) > tstmp; m-- {
		}
		if i == noframes-1 {
			m = nopoints - 1
			n = m
		}
		if m == n {
			alpha = float64(0)
		} else {
			alpha = (tstmp - float64(data[n].mtime)) / float64(data[m].mtime-data[n].mtime)
		}
		fmt.Printf("i=%d, n=%d, m=%d, tstmp=%f, mtimes=[%d,%d], alpha=%f\n", i, n, m, tstmp, data[n].mtime, data[m].mtime, alpha)
		task.afile = data[n].fname
		task.bfile = data[m].fname
		task.alpha = alpha
		ret = append(ret, task)
	}
	return ret
}

// Generate list of image (jpg) files under given dir
func findfiles(dir string) []string {
	vals, _ := filepath.Glob(dir + "/*.jpg")
	return vals
}

func main() {
	indir := os.Args[1]
	framestr := os.Args[2]
	outdir := os.Args[3]
	frames, _ := strconv.Atoi(framestr)
	noframes := int64(frames)
	dirlist := findfiles(indir)
	fmt.Printf("Would generate %d files into %s\n", frames, outdir)
	fmt.Println("Found files: ", dirlist)
	moredata := fun.Map(dirlist, getCtime)
	fmt.Println("Sorting data...")
	sort.Slice(moredata, func(i int, j int) bool {
		return moredata[i].mtime < moredata[j].mtime
	})
	fmt.Println(moredata)
	tasks := files2tasks(moredata, noframes)
	fmt.Println("Task list")
	fmt.Println(tasks)
}
