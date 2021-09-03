package imageprocessing





type HSVLUT struct {
  hsvMap map[float64]map[float64]map[float64]*RGBType
}

func NewHSVLUT() *HSVLUT {
	hsvLUT := HSVLUT{}

	hsvLUT.hsvMap = make(map[float64]map[float64]map[float64]*RGBType)
	return &hsvLUT
}

func (l *HSVLUT) Add(h float64, s float64, v float64, rgb RGBType ) {
	if _, ok := l.hsvMap[h]; !ok {
		l.hsvMap[h] = make(map[float64]map[float64]*RGBType)
	}

	if _, ok := l.hsvMap[h][s]; !ok {
		l.hsvMap[h][s] = make(map[float64]*RGBType)
	}

	if _, ok := l.hsvMap[h][s][v]; !ok {
		l.hsvMap[h][s][v] = &rgb
	}
}


func (l *HSVLUT) Get(h float64, s float64, v float64) (*RGBType, bool) {
	if l.hsvMap[h] != nil && l.hsvMap[h][s] != nil && l.hsvMap[h][s][v] != nil  {
		return l.hsvMap[h][s][v], true
	}
	return nil, false
}




