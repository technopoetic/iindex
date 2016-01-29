package iindex

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

// inverted  index representation
var Index map[string][]int // ints index into indexed
var Indexed []Doc

type Doc struct {
	File  string
	Title string
}

func IndexDirectory(path string) (map[string][]int, error) {
	// initialize representation
	Index = make(map[string][]int)

	// build index
	if err := indexDir(path); err != nil {
		return nil, err
	}
	return Index, nil
}

func indexDir(dir string) error {
	df, err := os.Open(dir)
	if err != nil {
		return err
	}
	fis, err := df.Readdir(-1)
	if err != nil {
		return err
	}
	if len(fis) == 0 {
		return fmt.Errorf("no files in %s", dir)
	}
	Indexed := 0
	for _, fi := range fis {
		if !fi.IsDir() {
			if indexFile(dir + "/" + fi.Name()) {
				Indexed++
			}
		}
	}
	return nil
}

func indexFile(fn string) bool {
	f, err := os.Open(fn)
	if err != nil {
		fmt.Println(err)
		return false // only false return
	}

	// register new file
	x := len(Indexed)
	Indexed = append(Indexed, Doc{fn, fn})
	pdoc := &Indexed[x]

	// scan lines
	r := bufio.NewReader(f)
	lines := 0
	for {
		b, isPrefix, err := r.ReadLine()
		switch {
		case err == io.EOF:
			return true
		case err != nil:
			fmt.Println(err)
			return true
		case isPrefix:
			fmt.Printf("%s: unexpected long line\n", fn)
			return true
		case lines < 20 && bytes.HasPrefix(b, []byte("Title:")):
			// in a real program you would write code
			// to skip the Gutenberg document header
			// and not index it.
			pdoc.Title = string(b[7:])
		}
		// index line of text in b
		// TODO: Write a better word splitter
	wordLoop:
		for _, bword := range bytes.Fields(b) {
			bword := bytes.Trim(bword, ".,-~?!\"'`;:()<>[]{}\\|/=_+*&^%$#@")
			if len(bword) > 0 {
				word := string(bword)
				dl := Index[word]
				for _, d := range dl {
					if d == x {
						continue wordLoop
					}
				}
				Index[word] = append(dl, x)
			}
		}
	}
}
