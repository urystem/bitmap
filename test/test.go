package main

import "fmt"

func cah4(b []byte) uint32 {
	if len(b) == 0 {
		return uint32(b[0])
	}
	if len(b) == 4 {
	}
	return uint32(b[0]) + cah4(b[1:])
}

func main() {
	s := "2-32-24-522"
	fmt.Println(SplitCrop(&s))
}

func SplitCrop(s *string) *[]uint32 {
	var crop []uint32
	var san string
	for i, v := range *s {
		if v != '-' {
			san += string(v)
		}
		if v == '-' || i == len(*s)-1 {
			fmt.Println(string(v), san)
			if n := Itoa(&san); n != nil {
				crop = append(crop, *n)
			} else {
				return nil
			}
			san = ""
		}
	}
	return &crop
}

func Itoa(s *string) *uint32 {
	if len(*s) == 0 {
		fmt.Println("dd")
		return nil
	}
	var n uint32
	for _, v := range *s {
		if v < '0' || v > '9' {
			return nil
		}
		n = n*10 + uint32(v-'0')
	}
	return &n
}
