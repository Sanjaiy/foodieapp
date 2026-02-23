package main

import (
	"bufio"
	"fmt"
	"log"
	"math/bits"
	"os"
	"sort"
	"strings"
)

func process(filename string, bit uint8, store map[string]uint8) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("opening %s: %w", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		code := scanner.Text()
		if len(code) >= 8 && len(code) <= 10 {
			store[code] |= bit
		}
	}

	return scanner.Err()
}

func main() {
	couponFiles := os.Getenv("COUPON_FILES")
	if couponFiles == "" {
		couponFiles = "data/coupon1.txt,data/coupon2.txt,data/coupon3.txt"
	}
	files := strings.Split(couponFiles, ",")

	outputPath := os.Getenv("OUTPUT_PATH")
	if outputPath == "" {
		outputPath = "data/valid_codes.txt"
	}

	store := make(map[string]uint8, 1000)

	for i, f := range files {
		log.Printf("Processing %s...", f)
		if err := process(f, 1<<uint(i), store); err != nil {
			log.Fatalf("Error processing %s: %v", f, err)
		}
	}

	var valid []string
	for code, mask := range store {
		if bits.OnesCount8(mask) >= 2 {
			valid = append(valid, code)
		}
	}

	sort.Strings(valid)

	out, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer out.Close()

	writer := bufio.NewWriter(out)
	for i, code := range valid {
		if i > 0 {
			writer.WriteByte('\n')
		}
		writer.WriteString(code)
	}
	writer.WriteByte('\n')
	writer.Flush()

	log.Printf("Done. %d valid codes written to %s", len(valid), outputPath)
}
