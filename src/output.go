package src

import "os"

func Usage(u uint8, e bool) {
	Printerln("Usage:")
	switch u {
	case 0:
		Printerln("  bitmap <command> [arguments]")
		Printerln("")
		Printerln("The commands are:")
		Printerln("  header    prints bitmap files header information")
		Printerln("  apply     applies processing to the image and saves it to the file")
	case 1:
		Printerln("  bitmap header <source_file>...")
		Printerln("")
		Printerln("Description:")
		Printerln("  Prints bitmap files header information")
	case 2:
		Printerln("  bitmap apply [options] <source_file> <output_file>")
		Printerln("")
		Printerln("The options are:")
		Printerln("  -h, --help	prints program usage information")
		Printerln("--mirror= (horizontal, h, horizontally, hor, vertical, v, vertically, ver)	Mirror the image horizontally or vertically. ")
		Printerln("--filter")
		Printerln("--rotate")
		Printerln("--crop")
	}
	if e {
		os.Exit(1)
	}
	os.Exit(0)
}

func Err(e error, b bool) {
	os.Stdin.WriteString("ERROR: " + e.Error() + "\n")
	if b {
		os.Exit(1)
	}
}

func Printerln(s string) {
	os.Stdout.WriteString(s + "\n")
}

func WriterBmp(m *[][][3]byte, s *string) error {
	var pad uint8 = uint8((4 - (len((*m)[0])*3)%4) % 4)
	b := make([]byte, len(*m)*len((*m)[0])*3+54+(int(pad)*len(*m)))
	b[0], b[1], b[10], b[14], b[26], b[28] = 'B', 'M', 54, 40, 1, 24                  // headersize, headerinfosize, sloi, bit
	writeHeh4(uint32(len(b)), b[2:6])                                                 // filesize
	writeHeh4(uint32(len((*m)[0])), b[18:22])                                         // width
	writeHeh4(uint32(len((*m))), b[22:26])                                            // height
	writeHeh4(uint32(len(*m)*len((*m)[0])*3)+uint32((pad))*uint32(len(*m)), b[34:38]) // image size
	writeHeh4(uint32(CountColor(m)), b[46:50])                                        // colorused
	var k uint32 = 54
	for _, v := range *m {
		for _, g := range v {
			b[k], b[k+1], b[k+2], k = g[0], g[1], g[2], k+3
		}
		for i := uint8(0); i < pad; i++ {
			b[k] = 0
			k++
		}
	}
	return os.WriteFile(*s, b, 0o644)
}

func writeHeh4(n uint32, b []byte) {
	b[0], b[1], b[2], b[3] = byte((n % 256)), byte(n/256%256), byte(n/65536%256), byte(n/16777216%256)
}
