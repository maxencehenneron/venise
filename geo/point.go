package geo

// A Point is a simple Lng/Lat 2d point.
type Point [2]float64

// Transform applies a given projection or inverse projection to the current point.
func (p *Point) Transform(projector Projector) *Point {
	projector(p)
	return p
}

// Creates a point from a latitude and longitude
func NewPointFromLatLng(lat, lng float64) *Point {
	return &Point{lng, lat}
}

// X returns the x/horizontal component of the point.
func (p *Point) X() float64 {
	return p[0]
}

// SetX sets the x/horizontal component of the point.
func (p *Point) SetX(x float64) *Point {
	p[0] = x
	return p
}

// Y returns the y/vertical component of the point.
func (p *Point) Y() float64 {
	return p[1]
}

// SetY sets the y/vertical component of the point.
func (p *Point) SetY(y float64) *Point {
	p[1] = y
	return p
}

// Lat returns the latitude/vertical component of the point.
func (p *Point) Lat() float64 {
	return p[1]
}

// SetLat sets the latitude/vertical component of the point.
func (p *Point) SetLat(lat float64) *Point {
	p[1] = lat
	return p
}

// Lng returns the longitude/horizontal component of the point.
func (p *Point) Lng() float64 {
	return p[0]
}

// SetLng sets the longitude/horizontal component of the point.
func (p *Point) SetLng(lng float64) *Point {
	p[0] = lng
	return p
}
