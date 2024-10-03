package src

import "errors"

func lastOut(arg *[]string) (bool, error) {
	if l := &(*arg)[len(*arg)-1]; !IsValid(l, false) { // if no file names
		return false, errors.New("GIVE bmp file")
	} else if len(*arg) == 1 {
		return false, errors.New("GIVE options there are only bmp file")
	} else if p := &(*arg)[len(*arg)-2]; IsValid(p, false) { // if has output and input files
		if len(*arg) == 2 { // if no options
			return false, errors.New("GIVE options there are only bmp files")
		}
		return true, nil
	}
	return false, nil
}

func Apply(arg *[]string) error {
	if f, e := lastOut(arg); e != nil {
		return e
	} else {
		var k uint8 = 1
		if f {
			k++
		}
		if er := applyFlags(arg, k); er != nil {
			return er
		}
		var m [][][3]byte
		if err := CheckAndHead(&(*arg)[len(*arg)-int(k)], &m); err != nil {
			return err
		} else {
			if errr := DoDone(&m, arg, k); errr != nil {
				return errr
			}
			WriterBmp(&m, &(*arg)[len(*arg)-1])
		}
		return nil
	}
}

func applyFlags(args *[]string, k uint8) error {
	for i, v := range (*args)[:len(*args)-int(k)] {
		if len(v) > 255 {
			return errors.New("this flag to long:" + v)
		} else if len(v) < 8 {
			return errors.New("invalid flag: " + v)
		} else if c, ca := v[:7], v[7:]; c == "--crop=" {
			if crop := splitCrop(&ca); crop == nil {
				return errors.New("invalid option for --crop=")
			} else if len(*crop) > 4 {
				return errors.New("too many cordinates in crop " + v)
			} else {
				(*args)[i] = ca
			}
		} else if len(v) < 10 {
			return errors.New("invalid flag: " + v)
		} else if c, ca = v[:9], v[9:]; c == "--mirror=" {
			if ca == "horizontal" || ca == "h" || ca == "horizontally" || ca == "hor" || ca == "vertical" || ca == "v" || ca == "vertically" || ca == "ver" {
				(*args)[i] = string(ca[0])
			} else {
				return errors.New("invalid option for --mirror " + v)
			}
		} else if c == "--filter=" {
			if ca == "blue" || ca == "red" || ca == "green" || ca == "grayscale" || ca == "negative" || ca == "pixelate" || ca == "blur" {
				(*args)[i] = ca
			} else {
				return errors.New("invalid option for --filter " + v)
			}
		} else if c == "--rotate=" {
			if ca == "right" || ca == "90" || ca == "-270" {
				(*args)[i] = "r"
			} else if ca == "left" || ca == "-90" || ca == "270" {
				(*args)[i] = "l"
			} else if ca == "180" || ca == "-180" {
				(*args)[i] = "c"
			} else {
				return errors.New("invalid option for --rotate= " + v)
			}
		} else {
			return errors.New("unknown flags " + (*args)[i])
		}
	}
	return nil
}

func splitCrop(s *string) *[]int {
	var crop []int
	var san string
	for i, v := range *s {
		if v != '-' {
			san += string(v)
		}
		if v == '-' || i == len(*s)-1 {
			if san == "auto" {
				crop = append(crop, -1)
			} else if n := itoa(&san); n != nil {
				crop = append(crop, int(*n))
			} else {
				return nil
			}
			san = ""
		}
	}
	return &crop
}

func itoa(s *string) *uint32 {
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

func DoDone(m *[][][3]byte, args *[]string, k uint8) error {
	var p int = 20
	for _, v := range (*args)[:len(*args)-int(k)] {
		switch v {
		case "h", "v":
			mirror(m, v == "h")
		case "blue", "green", "red", "grayscale", "negative":
			filterSimple(m, v[2]) // blue==117 green == 101 red=100 gray=97 negative=103
		case "pixelate":
			pixelate(m, &p)
			p += 10
		case "blur":
			blur(m)
		case "r", "l":
			rotate(m, v == "r")
		case "c":
			for i, j := 0, len(*m)-1; i < j; i, j = i+1, j-1 {
				(*m)[i], (*m)[j] = (*m)[j], (*m)[i]
			}
		default: // crop
			if e := crop(m, splitCrop(&v)); e != nil {
				return e
			}
		}
	}
	return nil
}

func filterSimple(m *[][][3]byte, n byte) {
	for i := range *m {
		for j := range (*m)[i] {
			switch n { // 0 blue //1 green //2 red
			case 97:
				g := byte(0.114*float64((*m)[i][j][0]) + 0.587*float64((*m)[i][j][1]) + 0.299*float64((*m)[i][j][2]))
				(*m)[i][j][0], (*m)[i][j][1], (*m)[i][j][2] = g, g, g
			case 100:
				(*m)[i][j][0], (*m)[i][j][1] = 0, 0
			case 101:
				(*m)[i][j][0], (*m)[i][j][2] = 0, 0
			case 103:
				(*m)[i][j][0], (*m)[i][j][1], (*m)[i][j][2] = 255-(*m)[i][j][0], 255-(*m)[i][j][1], 255-(*m)[i][j][2]
			case 117:
				(*m)[i][j][1], (*m)[i][j][2] = 0, 0
			}
		}
	}
}

func mirror(m *[][][3]byte, v bool) {
	if v {
		for i, j := 0, len(*m)-1; i < j; i, j = i+1, j-1 {
			(*m)[i], (*m)[j] = (*m)[j], (*m)[i]
		}
	} else {
		for i := range *m {
			for j, k := 0, len((*m)[i])-1; j < k; j, k = j+1, k-1 {
				(*m)[i][j], (*m)[i][k] = (*m)[i][k], (*m)[i][j]
			}
		}
	}
}

func rotate(m *[][][3]byte, r bool) {
	height, width := len(*m), len((*m)[0])
	c := make([][][3]byte, width) // Новая высота будет равна ширине оригинала
	for i := range c {
		c[i] = make([][3]byte, height) // Новая ширина будет равна высоте оригинала
		for j := range c[i] {
			if r {
				c[i][j] = (*m)[height-j-1][i]
			} else {
				c[i][j] = (*m)[j][width-1-i]
			}
		}
	}
	*m = c
}

func pixelate(m *[][][3]byte, p *int) {
	for i := 0; i < len(*m); i += *p {
		for j := 0; j < len((*m)[i]); j += *p {
			ei, ej := i+*p, j+*p
			if ei > len(*m) {
				ei = len(*m)
			}
			if ej > len((*m)[i]) {
				ej = len((*m)[i])
			}
			var c [3]uint32
			for ii := i; ii < ei; ii++ {
				for jj := j; jj < ej; jj++ {
					c[0], c[1], c[2] = c[0]+uint32((*m)[ii][jj][0]), c[1]+uint32((*m)[ii][jj][1]), c[2]+uint32((*m)[ii][jj][2])
				}
			}
			var k uint32 = uint32(((ei - i) * (ej - j)))
			c[0], c[1], c[2] = c[0]/k, c[1]/k, c[2]/k
			for ii := i; ii < ei; ii++ {
				for jj := j; jj < ej; jj++ {
					(*m)[ii][jj][0], (*m)[ii][jj][1], (*m)[ii][jj][2] = byte(c[0]), byte(c[1]), byte(c[2])
				}
			}
		}
	}
}

func blur(m *[][][3]byte) {
	c := make([][][3]byte, len(*m))
	for i := range c {
		c[i] = make([][3]byte, len((*m)[i]))
		for j := range c[i] {
			si, ei, sj, ej := i-10, i+10, j-10, j+10
			if si < 0 {
				si = 0
			}
			if ei > len(c) {
				ei = len(c)
			}
			if sj < 0 {
				sj = 0
			}
			if ej > len(c[i]) {
				ej = len(c[i])
			}
			var r [3]uint32
			for h := si; h < ei; h++ {
				for w := sj; w < ej; w++ {
					r[0], r[1], r[2] = r[0]+uint32((*m)[h][w][0]), r[1]+uint32((*m)[h][w][1]), r[2]+uint32((*m)[h][w][2])
				}
			}
			var k uint32 = uint32((ei - si) * (ej - sj))
			c[i][j][0], c[i][j][1], c[i][j][2] = byte(r[0]/k), byte(r[1]/k), byte(r[2]/k)
		}
	}
	*m = c
}

func crop(m *[][][3]byte, c *[]int) error {
	var oi, oj, ei, ej int = -1, -1, -1, -1
	if oi = (*c)[0]; oi >= len((*m)[0]) {
		return errors.New("over range offset x")
	}
	if len(*c) > 1 {
		if oj = (*c)[1]; oj >= len((*m)) {
			return errors.New("over range offset y")
		}
	}
	if len(*c) > 2 && (*c)[2] != -1 {
		if ei = oi + (*c)[2]; ei > len((*m)[0]) || oi >= ei {
			return errors.New("over range x")
		}
	}
	if len(*c) > 3 && (*c)[3] != -1 {
		if ej = oj + (*c)[3]; ej > len((*m)) || oj >= ej {
			return errors.New("over range y")
		}
	}
	if oi == -1 {
		oi = 0
	}
	if oj == -1 {
		oj = 0
	}
	if ei == -1 {
		ei = len((*m)[0])
	}
	if ej == -1 {
		ej = len(*m)
	}
	*m = (*m)[len(*m)-ej : len(*m)-oj]
	for i := range *m {
		(*m)[i] = (*m)[i][oi:ei]
	}
	return nil
}
