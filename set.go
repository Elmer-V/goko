package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func set() {
	//base_dir := getdir()
	//id := "P97969"
	if len(os.Args) < 3 {
		fmt.Println("usage: goko set <id> [app]")
		os.Exit(1)
	}
	id := os.Args[2]
	appid := "xdg-open"
	if len(os.Args) >= 4 {
		appid = os.Args[3]
	}
	//os.Mkdir(id, 0775)
	// Only extract when the download actually produced a file; otherwise the
	// unzip error below would just be "no such file", hiding the real cause.
	if !downloader(id) {
		os.Exit(1)
	}
	//os.Rename(id+".zip", filepath.Join(id, id+".zip"))
	// Check the error: a failed extraction used to be discarded, so a bad ZIP
	// exited 0 with no output and no message.
	if err := unzip(filepath.Join(id+".zip"), "."); err != nil {
		fmt.Println("unzip failed:", err)
		os.Exit(1)
	}
	// Only remove the archive once it has been extracted successfully.
	os.Remove(id + ".zip")
	os.Create(filepath.Join(id+probe(id), id+probe(id)+".cpp"))
	openApp := exec.Command(appid, filepath.Join(id+probe(id), id+probe(id)+".cpp"))
	if err := openApp.Start(); err != nil {
		fmt.Println("XDG error ", err)
		return
	}
}
