package main

import (
	"math"
)

func Direct(lat1 float64, lon1 float64, dist float64, azi float64) (lat2, lon2 float64) {

	radLat1 := lat1 * math.Pi / 180.0
	radLon1 := lon1 * math.Pi / 180.0
	radDist := dist / 6372795.0

	x, y, z := spherToCart(math.Pi/2.0-radDist, math.Pi-azi)
	z, x = rotate(z, x, radLat1-math.Pi/2)
	x, y = rotate(x, y, 0.0-radLon1)
	radLat2, radLon2 := cartToSpher(x, y, z)
	lat2, lon2 = radLat2*180/math.Pi, radLon2*180/math.Pi
	return
}

func spherToCart(lat float64, lon float64) (x, y, z float64) {
	x = math.Cos(lat) * math.Cos(lon)
	y = math.Cos(lat) * math.Sin(lon)
	z = math.Sin(lat)
	return x, y, z
}

func cartToSpher(x float64, y float64, z float64) (lat, lon float64) {
	lat = math.Atan2(z, math.Sqrt(x*x+y*y))
	lon = math.Atan2(y, x)
	return lat, lon
}

func rotate(x float64, y float64, a float64) (u, v float64) {
	c := math.Cos(a)
	s := math.Sin(a)
	u = x*c + y*s
	v = -x*s + y*c
	return u, v
}
