package helpers

import (
	"fmt"
	"os"
	"sort"
	"syscall"
)

type CouponLookup struct {
	data    []byte
	offsets []int
}

func NewCouponLookup(path string) (*CouponLookup, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("promo lookup open: %w", err)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("promo lookup stat: %w", err)
	}

	size := int(info.Size())
	if size == 0 {
		return &CouponLookup{}, nil
	}

	data, err := syscall.Mmap(int(f.Fd()), 0, size, syscall.PROT_READ, syscall.MAP_PRIVATE)
	if err != nil {
		return nil, fmt.Errorf("promo lookup mmap: %w", err)
	}

	offsets := []int{0}
	for i := 0; i < size; i++ {
		if data[i] == '\n' && i+1 < size {
			offsets = append(offsets, i+1)
		}
	}

	return &CouponLookup{data: data, offsets: offsets}, nil
}

func (p *CouponLookup) lineAt(idx int) string {
	start := p.offsets[idx]
	end := start
	for end < len(p.data) && p.data[end] != '\n' {
		end++
	}
	return string(p.data[start:end])
}

func (p *CouponLookup) IsValid(code string) bool {
	if len(code) < 8 || len(code) > 10 {
		return false
	}
	if len(p.offsets) == 0 {
		return false
	}

	i := sort.Search(len(p.offsets), func(i int) bool {
		return p.lineAt(i) >= code
	})

	return i < len(p.offsets) && p.lineAt(i) == code
}

func (p *CouponLookup) Close() error {
	if p.data != nil {
		return syscall.Munmap(p.data)
	}
	return nil
}
