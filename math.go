package main

import (
	"math"

	"github.com/westphae/quaternion"
)

func add3f(x, y [3]float64) [3]float64 {
	return [3]float64{x[0] + y[0], x[1] + y[1], x[2] + y[2]}
}

func sub3f(x, y [3]float64) [3]float64 {
	return [3]float64{x[0] - y[0], x[1] - y[1], x[2] - y[2]}
}

func mul3f(x [3]float64, y float64) [3]float64 {
	return [3]float64{x[0] * y, x[1] * y, x[2] * y}
}

func combineEuler(x, y [3]float64) [3]float64 {
	xq := quaternion.FromEuler(x[2], x[1], x[0])
	yq := quaternion.FromEuler(y[2], y[1], y[0])
	qx, qy, qz := quaternion.Prod(yq, xq).Euler()
	return [3]float64{qz, qy, qx}
}

func quatAdd(q, p [4]float64) [4]float64 {
	return [4]float64{
		q[3]*p[0] - q[2]*p[1] + q[1]*p[2] + q[0]*p[3],
		q[2]*p[0] + q[3]*p[1] - q[0]*p[2] + q[1]*p[3],
		-q[1]*p[0] + q[0]*p[1] + q[3]*p[2] + q[2]*p[3],
		-q[0]*p[0] - q[1]*p[1] - q[2]*p[2] + q[3]*p[3],
	}
}

func eulerToQuat(e [3]float64) [4]float64 {
	c := [3]float64{cos(e[0]), cos(e[1]), cos(e[2])}
	s := [3]float64{sin(e[0]), sin(e[1]), sin(e[2])}
	return [4]float64{
		c[0]*s[1]*s[2] - s[0]*c[1]*c[2],
		-c[0]*s[1]*c[2] - s[0]*c[1]*s[2],
		c[0]*c[1]*s[2] - s[0]*s[1]*c[2],
		c[0]*c[1]*c[2] + s[0]*s[1]*s[2],
	}
}

func quatToEuler(q [4]float64) [3]float64 {
	return [3]float64{
		atan2(2.0*(q[2]*q[3]+q[0]*q[1]), q[0]*q[0]-q[1]*q[1]-q[2]*q[2]+q[3]*q[3]),
		asin(2.0 * (q[0]*q[2] - q[1]*q[3])),
		atan2(2.0*(q[1]*q[2]+q[0]*q[3]), q[0]*q[0]+q[1]*q[1]-q[2]*q[2]-q[3]*q[3]),
	}
}

func cos(x float64) float64 {
	return math.Cos(x)
}

func sin(x float64) float64 {
	return math.Sin(x)
}

func atan2(y, x float64) float64 {
	return math.Atan2(y, x)
}

func asin(x float64) float64 {
	return math.Asin(x)
}
