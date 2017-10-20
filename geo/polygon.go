package geo

type Polygon struct {
	Rings [][]Point
}

func (p *Polygon) Transform(projector Projector) {
	for ridx, _ := range p.Rings {
		for pidx, _ := range p.Rings[ridx] {
			p.Rings[ridx][pidx].Transform(projector)
		}
	}
}
