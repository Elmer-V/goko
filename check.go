package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func check() {

	id := "."
	nameid := "."
	if len(os.Args) >= 3 {
		id = os.Args[2] + probe(os.Args[2])
		nameid = id
	} else {
		wd, _ := os.Getwd()
		id = wd
		nameid = filepath.Base(wd)
	}
	// count, err := getCount(id)
	// if err != nil {
	// 	panic(err)
	// }
	stems, err := sampleStems(id)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	build(id, nameid)
	getResults(id, nameid, stems)
	CompareResults(id, stems)

}

func sampleStems(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read dir %q: %w", dir, err)
	}

	var stems []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".inp") {
			continue
		}
		stems = append(stems, strings.TrimSuffix(name, ".inp"))
	}

	if len(stems) == 0 {
		return nil, fmt.Errorf("no sample inputs in %q", dir)
	}
	return stems, nil
}

func getCount(dir string) (int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, fmt.Errorf("read dir %q: %w", dir, err)
	}

	const prefix = "sample-"
	const suffix = ".inp"

	maxIdx := 0
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasPrefix(name, prefix) || !strings.HasSuffix(name, suffix) {
			continue
		}
		// Extract the number between prefix and suffix.
		numPart := name[len(prefix) : len(name)-len(suffix)]
		n, err := strconv.Atoi(numPart)
		if err != nil {
			continue // not a numeric index; skip
		}
		if n > maxIdx {
			maxIdx = n
		}
	}
	return maxIdx, nil
}

func build(id, nameid string) {
	p1build := exec.Command("clang++", filepath.Join(id, nameid+".cpp"), "-o", filepath.Join(id, nameid+".x"))
	if err := p1build.Run(); err != nil {
		fmt.Println("build error ", err)
	}
}

// func getResults(id string, count int) {
// 	for i := 1; i <= count; i++ {
// 		runTests := exec.Command(filepath.Join(id, id+".x"), "<", filepath.Join(id, "sample-"+strconv.Itoa(i)+".inp"), ">", filepath.Join(id, "sample-"+strconv.Itoa(i)+".out"))
// 		runTests.Run()
// 	}
// }

func getResults(id, nameid string, stems []string) {
	for i, stem := range stems {
		inPath := filepath.Join(id, stem+".inp")
		outPath := filepath.Join(id, stem+".out")

		in, err := os.Open(inPath) // *os.File, error  — read-only
		if err != nil {
			fmt.Println("open input:", err)
			return
		}

		out, err := os.Create(outPath) // *os.File, error  — creates/truncates
		if err != nil {
			in.Close()
			fmt.Println("create output:", err)
			return
		}

		runTests := exec.Command(filepath.Join(id, nameid+".x"))
		runTests.Stdin = in   // io.Reader
		runTests.Stdout = out // io.Writer
		// runTests.Stderr = os.Stderr       // uncomment to see your program's own errors

		err = runTests.Run()

		in.Close() // close AFTER the run
		out.Close()

		if err != nil {
			fmt.Println("test", i, "run error:", err)
		}
	}
}

func CompareResults(id string, stems []string) {
	for i, stem := range stems {
		cor := filepath.Join(id, stem+".cor")
		out := filepath.Join(id, stem+".out")
		// var out []byte
		// var err error
		output, err := exec.Command("diff", "-q", cor, out).CombinedOutput()
		if err == nil {
			fmt.Println(i, "OK") // exit 0
		} else if ee, ok := err.(*exec.ExitError); ok && ee.ExitCode() == 1 {
			fmt.Println(i, "DIFFER") // exit 1
		} else {
			fmt.Println(i, "error:", err, string(output)) // exit 2+ or couldn't run
		}
	}
}
