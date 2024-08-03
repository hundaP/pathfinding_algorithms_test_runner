package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/tabwriter"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run csv_renderer.go <file_or_directory> [file_or_directory] ...")
		os.Exit(1)
	}

	for _, arg := range os.Args[1:] {
		info, err := os.Stat(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accessing %s: %v\n", arg, err)
			continue
		}

		if info.IsDir() {
			processDirectory(arg)
		} else {
			fmt.Println("File:", arg)
			renderCSV(arg)
		}
	}
}

func processDirectory(dir string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading directory %s: %v\n", dir, err)
		return
	}

	for _, fileInfo := range files {
		if !fileInfo.IsDir() && filepath.Ext(fileInfo.Name()) == ".csv" {
			filePath := filepath.Join(dir, fileInfo.Name())
			fmt.Println("File:", filePath)
			renderCSV(filePath)
		}
	}
}

func renderCSV(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','

	// Determine column widths for better formatting
	maxColumnWidths := make([]int, 0)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			return
		}

		if len(maxColumnWidths) == 0 {
			maxColumnWidths = make([]int, len(record))
		}

		for i, field := range record {
			if len(field) > maxColumnWidths[i] {
				maxColumnWidths[i] = len(field)
			}
		}
	}

	// Rewind the file to read it again
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error rewinding file: %v\n", err)
		return
	}

	reader = csv.NewReader(file)
	reader.Comma = ','

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 4, 2, ' ', 0)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			return
		}

		for i, field := range record {
			fmt.Fprintf(w, "%-*s\t", maxColumnWidths[i]+2, field)
		}
		fmt.Fprintln(w)
	}
	w.Flush()
	fmt.Println()
}
