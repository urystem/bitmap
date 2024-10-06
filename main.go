package main

import (
	"fmt"
	"os"
	"time"

	"bitmap/src"
)

func main() {
	start := time.Now() // Запоминаем текущее время
	os.Args = os.Args[1:]
	if len(os.Args) == 0 || src.IsValid(&os.Args[0], true) {
		src.Usage(0, false)
	}
	switch os.Args[0] {
	case "header":
		os.Args = os.Args[1:]
		if len(os.Args) == 0 || src.IsValid(&os.Args[0], true) {
			src.Usage(1, false)
		}
		src.Header(&os.Args)
	case "apply":
		os.Args = os.Args[1:]
		if len(os.Args) == 0 || src.IsValid(&os.Args[0], true) {
			src.Usage(2, false)
		}
		if e := src.Apply(&os.Args); e != nil {
			src.Err(e, true)
		}
	default:
		src.Usage(0, true)
	}
	fmt.Printf("Время выполнения: %v\n", time.Since(start))
}
