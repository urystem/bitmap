package src

import (
	"errors"
	"os"
)

func Header(arg *[]string) {
	var er bool
	for i, v := range *arg {
		if i != 0 {
			Printerln("")
		}
		if e := CheckAndHead(&v, nil); e != nil {
			Err(e, false)
			er = true
		}
	}
	if er {
		os.Exit(1)
	}
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
		} else if m != nil && b[0] != 66 && b[1] != 77 {//
			return errors.New("invalid signature of bmp file")
		} else if s, e := f.Stat(); e != nil {
			return e
		} else if b[5] == 255 {
			return errors.New("filefize seted too big")
		} else if q := cah4(b[2:6]); m != nil && q != uint32(s.Size()) {
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
				Printerln("File: " + *fn)
				Printerln("BMP Header:")
				Printerln("- FileType " + string(rune(b[0])) + string(rune(b[1])))
				Printerln("- FileSizeInBytes " + *atoiu(&q))
				Printerln("- HeaderSize " + *atoiu(&h))
				Printerln("DIB Header:")
				Printerln("- DibHeaderSize " + *atoiu(&hf))
				Printerln("- WidthInPixels " + *atoiu(&width))
				Printerln("- HeightInPixels " + *atoiu(&height))
				Printerln("- PixelSizeInBits " + *atoiu(&ib))
				Printerln("- ImageSizeInBytes " + if0(&l))
				Printerln("- Colors Used " + if0(&colors))
				return nil
			} else {
				b = make([]byte, l-54+h)
				if _, t := f.Read(b); t != nil {
					return t
				} else {
					b = b[h-54:]
				}
			}
		}
		var pad uint8 = uint8((4 - (width*3)%4) % 4)
		var k uint32
		*m = make([][][3]byte, height)
		for i := range *m {
			for j := uint32(0); j < width; j++ {
				(*m)[i] = append((*m)[i], [3]byte{})
				(*m)[i][j][0], (*m)[i][j][1], (*m)[i][j][2], k = b[k], b[k+1], b[k+2], k+3
			}
			k += uint32(pad)
		}
		if colors != 0 && colors != CountColor(m) {
			return errors.New("invalid colour counts")
		}
		return nil
	}
}

func cah4(b []byte) uint32 {
	return uint32(b[0]) + uint32(b[1])*256 + uint32(b[2])*256*256 + uint32(b[3])*256*256*256
}

func if0(n *uint32) string {
	if *n == 0 {
		return "0 (auto)"
	}
	return *atoiu(n)
}

func atoiu(n *uint32) *string {
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
