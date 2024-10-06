package src

import "sync"

func Blurgo(m *[][][3]byte) {
	c := make([][][3]byte, len(*m))
	for k := 0; k < len(*m); k += 100 {
		var ch [100]chan bool
		for i := range ch {
			if k+i < len(*m) {
				ch[i] = make(chan bool)
				go blurgoH(m, &c, k+i, ch[i])
			} else {
				break
			}
		}
		for i := range ch {
			if k+i < len(*m) {
				<-ch[i]
				close(ch[i])
			} else {
				break
			}
		}
	}

	*m = c
}

func BlurgoSy(m *[][][3]byte) {
	c := make([][][3]byte, len(*m))
	var wg sync.WaitGroup
	for i := range *m {
		wg.Add(1)
		go blurgoSyH(m, &c, i, &wg)
	}
	wg.Wait()
	*m = c
}

func blurgoSyH(m, c *[][][3]byte, i int, wg *sync.WaitGroup) {
	(*c)[i] = make([][3]byte, len((*m)[i]))
	for j := range (*m)[i] {
		si, ei, sj, ej := i-10, i+10, j-10, j+10
		if si < 0 {
			si = 0
		}
		if ei > len(*m) {
			ei = len(*m)
		}
		if sj < 0 {
			sj = 0
		}
		if ej > len((*m)[i]) {
			ej = len((*m)[i])
		}
		var r [3]uint32
		for h := si; h < ei; h++ {
			for w := sj; w < ej; w++ {
				r[0], r[1], r[2] = r[0]+uint32((*m)[h][w][0]), r[1]+uint32((*m)[h][w][1]), r[2]+uint32((*m)[h][w][2])
			}
		}
		var k uint32 = uint32((ei - si) * (ej - sj))
		(*c)[i][j][0], (*c)[i][j][1], (*c)[i][j][2] = byte(r[0]/k), byte(r[1]/k), byte(r[2]/k)
	}
	(*wg).Done()
}

func blurgoH(m, c *[][][3]byte, i int, ch chan bool) {
	(*c)[i] = make([][3]byte, len((*m)[i]))
	for j := range (*m)[i] {
		si, ei, sj, ej := i-10, i+10, j-10, j+10
		if si < 0 {
			si = 0
		}
		if ei > len(*m) {
			ei = len(*m)
		}
		if sj < 0 {
			sj = 0
		}
		if ej > len((*m)[i]) {
			ej = len((*m)[i])
		}
		var r [3]uint32
		for h := si; h < ei; h++ {
			for w := sj; w < ej; w++ {
				r[0], r[1], r[2] = r[0]+uint32((*m)[h][w][0]), r[1]+uint32((*m)[h][w][1]), r[2]+uint32((*m)[h][w][2])
			}
		}
		var k uint32 = uint32((ei - si) * (ej - sj))
		(*c)[i][j][0], (*c)[i][j][1], (*c)[i][j][2] = byte(r[0]/k), byte(r[1]/k), byte(r[2]/k)
	}
	ch <- true
}
