package geo

import "math"

// A Projector is a function that converts the given point to a different space.
type Projector func(p *Point)

// A Projection is a set of projectors to map forward and backwards to the projected space.
type Projection struct {
	Project Projector
	Inverse Projector
}

const mercatorPole = 20037508.34

// Mercator projection, performs EPSG:3857, sometimes also described as EPSG:900913.
var Mercator = Projection{
	Project: func(p *Point) {
		p.SetX(mercatorPole / 180.0 * p.Lng())

		y := math.Log(math.Tan((90.0+p.Lat())*math.Pi/360.0)) / math.Pi * mercatorPole
		p.SetY(math.Max(-mercatorPole, math.Min(y, mercatorPole)))
	},
	Inverse: func(p *Point) {
		p.SetLng(p.X() * 180.0 / mercatorPole)
		p.SetLat(180.0 / math.Pi * (2*math.Atan(math.Exp((p.Y()/mercatorPole)*math.Pi)) - math.Pi/2.0))
	},
}
