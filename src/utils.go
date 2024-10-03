package src

func IsValid(s *string, h bool) bool {
	if h && (*s == "-h" || (len(*s) > 5 && (*s)[:5] == "--help")) {
		return true
	} else if !h && (len(*s) > 4 && (*s)[len(*s)-4:] == ".bmp") {
		return true
	}
	return false
}

func CountColor(m *[][][3]byte) uint32 {
	u := make(map[[3]byte]struct{})
	for _, v := range *m {
		for _, g := range v {
			u[g] = struct{}{}
		}
	}
	return uint32(len(u))
}
