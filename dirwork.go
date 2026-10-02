package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func getdir() string {
	out, err := os.Getwd()
	if err != nil {
		fmt.Printf("error")
	}
	return out
}

func unzip(src, dst string) error {
	// 1. Open the zip
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	// 2. Loop over every entry inside
	for _, f := range r.File {
		target := filepath.Join(dst, f.Name)

		// Directory entry? Just create it.
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
			continue
		}

		// Make sure the parent directory exists.
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}

		// Open the entry for reading
		rc, err := f.Open()
		if err != nil {
			return err
		}

		// Create the destination file
		out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			rc.Close()
			return err
		}

		// Copy the data and close both handles
		if _, err := io.Copy(out, rc); err != nil {
			out.Close()
			rc.Close()
			return err
		}
		out.Close()
		rc.Close()
	}
	return nil
}
