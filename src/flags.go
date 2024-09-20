package src

import (
	"flag"
	"fmt"
	"strings"
)

type flagvalues []string

func (ss *flagvalues) String() string {
	return strings.Join(*ss, ", ")
}

func (ss *flagvalues) Set(value string) error {
	*ss = append(*ss, value)
	return nil
}

func Flag() (bool, []string, []string, []string, []string) {
	var (
		filters flagvalues
		rotates flagvalues
		mirror  flagvalues
		crop    flagvalues
	)
	flag.Var(&filters, "filter", "")
	flag.Var(&rotates, "rotate", "")
	flag.Var(&mirror, "mirror", "")
	flag.Parse()
	fmt.Println("Filters:", filters)
	fmt.Println("Rotates:", rotates)
	fmt.Println("Mirrors:", mirror)
	fmt.Println("crops:", crop)
	return true, filters, rotates, mirror, crop
}
