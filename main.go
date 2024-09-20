package main

import (
	"errors"
	"fmt"
	"os"
)

func Usage(u uint8, e bool) {
	fmt.Println("Usage:")
	switch u {
	case 0:
		fmt.Println("  bitmap <command> [arguments]")
		fmt.Println()
		fmt.Println("The commands are:")
		fmt.Println("  header    prints bitmap files header information")
		fmt.Println("  apply     applies processing to the image and saves it to the file")
	case 1:
		fmt.Println("  bitmap header <source_file>...")
		fmt.Println()
		fmt.Println("Description:")
		fmt.Println("  Prints bitmap files header information")
	case 2:
		fmt.Println("  bitmap apply [options] <source_file> <output_file>")
		fmt.Println()
		fmt.Println("The options are:")
		fmt.Println("  -h, --help	prints program usage information")
		fmt.Println("--mirror= (horizontal, h, horizontally, hor, vertical, v, vertically, ver)	Mirror the image horizontally or vertically. ")
		fmt.Println("--filter")
		fmt.Println("--rotate")
		fmt.Println("--crop")
	}
	if e {
		os.Exit(1)
	}
	os.Exit(0)
}

func Err(e error) {
	os.Stdin.WriteString("ERROR: " + e.Error() + "\n")
	os.Exit(1)
}

func IsValid(s *string, h bool) bool {
	if h && (*s == "-h" || (len(*s) > 5 && (*s)[:5] == "--help")) {
		return true
	} else if !h && (len(*s) > 4 && (*s)[len(*s)-4:] == ".bmp") {
		return true
	}
	return false
}

func CheckAndHead(fn *string, m *[][][3]byte) error {
	if f, e := os.Open(*fn); e != nil {
		return e
	} else {
		defer f.Close()
		b := make([]byte, 54)
		var width, height, colors uint32
		if a, e := f.Read(b); e != nil {
			return e
		} else if a != 54 {
			return errors.New("invalid size of bmp file")
		} else if m != nil && b[0] != 66 && b[1] != 77 {
			return errors.New("invalid signature of bmp file")
		} else if s, e := f.Stat(); e != nil {
			return e
		} else if q := cah4(b[2:6]); m != nil && int64(q) != s.Size() {
			return errors.New("неправильео указан размер файла")
		} else if h, hf := cah4(b[10:14]), cah4(b[14:18]); m != nil && (h-hf != 14 || (hf != 40 && hf != 108 && hf != 124)) {
			return errors.New("invalid version. This program support 3 4 5 versions")
		} else if m != nil && (b[26] != 1 || b[27] != 0) {
			return errors.New("this program support only 1 bipanes")
		} else if ib := uint32(b[28]) + uint32(b[29])*256; m != nil && ib != 24 {
			return errors.New("this program support only 24 bit")
		} else if m != nil && cah4(b[30:34]) != 0 {
			return errors.New("this program not support compression")
		} else if l := cah4(b[34:38]); l > q-h {
			return errors.New("image size incorrect")
		} else {
			width, height, colors = cah4(b[18:22]), cah4(b[22:26]), cah4(b[46:50])
			if width*height > l {
				return errors.New("размер файл не правильно указан")
			} else if m == nil {
				os.Stdout.WriteString("File: " + *fn)
				os.Stdout.WriteString("\nBMP Header:")
				os.Stdout.WriteString("\n- FileType " + string(rune(b[0])) + string(rune(b[1])))
				os.Stdout.WriteString("\n- FileSizeInBytes " + *Atoiu(&q))
				os.Stdout.WriteString("\n- HeaderSize " + *Atoiu(&h))
				os.Stdout.WriteString("\nDIB Header:")
				os.Stdout.WriteString("\n- DibHeaderSize " + *Atoiu(&hf))
				os.Stdout.WriteString("\n- WidthInPixels " + *Atoiu(&width))
				os.Stdout.WriteString("\n- HeightInPixels " + *Atoiu(&height))
				os.Stdout.WriteString("\n- PixelSizeInBits " + *Atoiu(&ib))
				os.Stdout.WriteString("\n- ImageSizeInBytes " + if0(&l))
				os.Stdout.WriteString("\n- Colors Used " + if0(&colors) + "\n")
				return nil
			} else {
				b = make([]byte, l-54+uint32(h))
				if _, t := f.Read(b); t != nil {
					return t
				} else {
					b = b[h-54:]
				}
			}
		}
		var k uint64
		var cn uint32
		var c [][3]byte
		for i := 0; i < int(width); i++ {
			*m = append(*m, [][3]byte{})
			for j := 0; j < int(height); j++ {
				(*m)[i] = append((*m)[i], [3]byte{})
				(*m)[i][j][0] = b[k]
				(*m)[i][j][1] = b[k+1]
				(*m)[i][j][2] = b[k+2]
				k += 3
				if colors != 0 && checkColor(&c, &(*m)[i][j]) {
					cn++
					if cn > colors {
						return errors.New("invalid colour counts")
					}
				}
			}
		}
		if cn < colors {
			return errors.New("invalid colour counts")
		}
		return nil
	}
}

func if0(n *uint32) string {
	if *n == 0 {
		return "0 (auto)"
	}
	return *Atoiu(n)
}

func checkColor(c *[][3]byte, b *[3]byte) bool {
	for _, v := range *c {
		if v == *b {
			return false
		}
	}
	*c = append(*c, *b)
	return true
}

func cah4(b []byte) uint32 {
	return uint32(b[0]) + uint32(b[1])*256 + uint32(b[2])*256*256 + uint32(b[3])*256*256*256
}

func main() {
	var head bool
	if len(os.Args) == 1 || IsValid(&os.Args[1], true) {
		Usage(0, false)
	}
	switch os.Args[1] {
	case "header":
		head = true
		if len(os.Args) == 2 || IsValid(&os.Args[2], true) {
			Usage(1, false)
		}
	case "apply":
		if len(os.Args) == 2 || IsValid(&os.Args[2], true) {
			Usage(2, false)
		}
	default:
		Usage(0, true)
	}
	if head {
		var err bool
		for i, v := range os.Args[2:] {
			if i != 0 {
				os.Stdout.WriteString("\n")
			}
			if e := CheckAndHead(&v, nil); e != nil {
				fmt.Println(e)
				err = true
			}
		}
		if err {
			os.Exit(1)
		}
	} else {
		var f, o *string
		if l := &os.Args[len(os.Args)-1]; !IsValid(l, false) { // if no file names
			Err(errors.New("GIVE bmp file"))
		} else if p := &os.Args[len(os.Args)-2]; IsValid(p, false) { // if has output and input files
			if len(os.Args) < 5 { // if no options
				Err(errors.New("GIVE options"))
			}
			o, f, os.Args = l, p, os.Args[2:len(os.Args)-2]
		} else if len(os.Args) < 4 {
			Err(errors.New("GIVE options"))
		} else {
			f, o, os.Args = l, l, os.Args[2:len(os.Args)-1]
		}
		if err := ApplyOpCheck(&os.Args); err != nil {
			Err(err)
		} else {
			var m [][][3]byte
			if er := CheckAndHead(f, &m); er != nil {
				Err(er)
			} else {
				fmt.Println(o)
			}
		}
	}
}

func ApplyOpCheck(args *[]string) error {
	for i, v := range *args {
		if len(v) > 255 {
			return errors.New("too long flag")
		}
		if len(v) < 8 {
			return errors.New("too few flag")
		} else {
			if c, ca := v[:9], v[9:]; c == "--mirror=" {
				if ca == "horizontal" || ca == "h" || ca == "horizontally" || ca == "hor" || ca == "vertical" || ca == "v" || ca == "vertically" || ca == "ver" {
					(*args)[i] = string(ca[0])
				} else {
					return errors.New("invalid option for --mirror")
				}
			} else if c == "--filter=" {
				if ca == "blue" || ca == "red" || ca == "green" || ca == "grayscale" || ca == "negative" || ca == "pixelate" || ca == "blur" {
					(*args)[i] = ca
				} else {
					return errors.New("invalid option for --mirror")
				}
			} else if c == "--rotate=" {
				if ca == "right" || ca == "90" || ca == "-270" {
					(*args)[i] = "90"
				} else if ca == "left" || ca == "-90" || ca == "270" {
					(*args)[i] = "-90"
				} else if ca == "180" || ca == "-180" {
					(*args)[i] = "180"
				} else {
					return errors.New("invalid option for --rotate=")
				}
			} else if ca := v[7:]; v[:7] == "--crop=" {
				if crop := SplitCrop(&ca); crop == nil {
					return errors.New("invalid option for --crop=")
				} else if len(*crop) != 2 && len(*crop) != 4 {
					return errors.New("incorrect crop 2 4")
				} else {
					(*args)[i] = ca
				}
			} else {
				return errors.New("unknown flags")
			}
		}
	}
	return nil
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

func Atoiu(n *uint32) *string {
	var s string
	if *n == 0 {
		s = "0"
	}
	for *n != 0 {
		s = string(rune(*n%10+'0')) + s
		*n /= 10
	}
	return &s
}

func DoDone(m *[][]byte, args *[]string) error {
	return nil
}
