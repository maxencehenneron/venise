package geo

import (
	"fmt"
	"math"
)

type Cell struct {
	X   float64 // Cell center X
	Y   float64 // Cell center Y
	H   float64 // Half cell size
	D   float64 // Distance from cell center to polygon
	Max float64 // max distance to polygon within a cell
}

func NewCell(x float64, y float64, h float64, p Polygon) *Cell {
	var d float64 = pointToPolygonDist(x, y, p)
	return &Cell{
		x,
		y,
		h,
		d,
		d + h*math.Sqrt2,
	}
}

func pointToPolygonDist(x float64, y float64, p Polygon) float64 {
	inside := false
	minDistSq := math.MaxFloat64

	for k := 0; k < len(p.Rings); k++ {
		ring := p.Rings[k]

		for i, length, j := 0, len(ring), len(ring)-1; i < length; i, j = i+1, i {
			a := ring[i]
			b := ring[j]

			if ((a.Y() > y) != (b.Y() > y)) &&
				(x < (b.X()-a.X())*(y-a.Y())/(b.Y()-a.Y())+a.X()) {
				inside = !inside
			}

			minDistSq = math.Min(minDistSq, getSegDistSq(x, y, a, b))
		}
	}

	if inside {
		return math.Sqrt(minDistSq)
	} else {
		return -1 * math.Sqrt(minDistSq)
	}
}

func getSegDistSq(px float64, py float64, a Point, b Point) float64 {
	var x = a.X()
	var y = a.Y()
	var dx = b.X() - x
	var dy = b.Y() - y

	if (dx != 0) || (dy != 0) {

		var t = ((px-x)*dx + (py-y)*dy) / (dx*dx + dy*dy)

		if t > 1 {
			x = b[0]
			y = b[1]

		} else if t > 0 {
			x += dx * t
			y += dy * t
		}
	}

	dx = px - x
	dy = py - y

	return dx*dx + dy*dy
}

func getCentroidCell(polygon Polygon) *Cell {
	var area = 0.0
	var x = 0.0
	var y = 0.0
	var points = polygon.Rings[0]

	for i, length, j := 0, len(points), len(points)-1; i < length; i, j = i+1, i {
		var a = points[i]
		var b = points[j]
		var f = a.X()*b.Y() - b.X()*a.Y()
		x += (a.X() + b.X()) * f
		y += (a.Y() + b.Y()) * f
		area += f * 3
	}

	if area == 0 {
		return NewCell(points[0].X(), points[0].Y(), 0, polygon)
	} else {
		return NewCell(x/area, y/area, 0, polygon)
	}
}

func Polylabel(polygon Polygon, precision float64, debug bool) Point {
	minX, minY, maxX, maxY := math.MaxFloat64, math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64

	for i := 0; i < len(polygon.Rings[0]); i++ {
		var p = polygon.Rings[0][i]

		minX = math.Min(minX, p.X())
		minY = math.Min(minY, p.Y())
		maxX = math.Max(maxX, p.X())
		maxY = math.Max(maxY, p.Y())
	}

	var width = maxX - minX
	var height = maxY - minY
	var cellSize = math.Min(width, height)
	var h = cellSize / 2

	cellQueue := NewPriorityQueue()

	if cellSize == 0 {
		return Point{minX, minY}
	}

	for x := minX; x < maxX; x += cellSize {
		for y := minY; y < maxY; y += cellSize {
			cell := NewCell(x+h, y+h, h, polygon)
			cellQueue.Insert(*cell, cell.Max)
		}
	}

	bestCell := getCentroidCell(polygon)

	var bboxCell = NewCell(minX+width/2, minY+height/2, 0, polygon)
	if bboxCell.D > bestCell.D {
		bestCell = bboxCell
	}

	numProbes := cellQueue.Len()

	for cellQueue.Len() != 0 {
		cellInterface, _ := cellQueue.Pop()
		cell := cellInterface.(Cell)

		// update the best cell if we found a better one
		if cell.D > bestCell.D {
			bestCell = &cell
		}

		// do not drill down further if there's no chance of a better solution
		if cell.Max-bestCell.D <= precision {
			continue
		}

		// split the cell into four cells
		h = cell.H / 2
		cells := []*Cell{
			NewCell(cell.X-h, cell.Y-h, h, polygon),
			NewCell(cell.X+h, cell.Y-h, h, polygon),
			NewCell(cell.X-h, cell.Y+h, h, polygon),
			NewCell(cell.X+h, cell.Y+h, h, polygon),
		}

		for _, ncell := range cells {
			cellQueue.Insert(*ncell, ncell.Max)
		}

		numProbes += 4
	}

	if debug {
		fmt.Sprintf("%d, num probes: ", numProbes)
	}

	return Point{bestCell.X, bestCell.Y}
}
