package particle

import "math"

func Center(points []Point) Point {
	if len(points) == 0 {
		return Point{}
	}
	x, y, alpha := 0, 0, 0
	for _, point := range points {
		x += point.X
		y += point.Y
		alpha += point.Alpha
	}
	return Point{X: x / len(points), Y: y / len(points), Alpha: alpha / len(points)}
}

func Spread(points []Point) float64 {
	if len(points) < 2 {
		return 0
	}
	center := Center(points)
	total := 0.0
	for _, point := range points {
		dx := float64(point.X - center.X)
		dy := float64(point.Y - center.Y)
		total += math.Sqrt(dx*dx + dy*dy)
	}
	return total / float64(len(points))
}

func Normalize(points []Point, width, height int) []Point {
	result := make([]Point, len(points))
	for i, point := range points {
		if width > 0 {
			point.X = point.X % width
		}
		if height > 0 {
			point.Y = point.Y % height
		}
		result[i] = point
	}
	return result
}
