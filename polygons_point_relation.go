package manor

import (
	"errors"
)

// 点和多边形的位置关系计算, 使用射线和多边形的交点数判断点在多边形内部或者外部
// 计算都采用直线的一般式方程：Ax + By + C = 0

type Coordinate struct {
	X int
	Y int
}

// 方程一般式结构体
type GeneralParameter struct {
	A int
	B int
	C int
}

// 给定一点斜率为1的点斜式方程
// (3, 6) => x - y + 3 = 0
func pointSlopeForm(point Coordinate) (int, int, int) {
	return 1, -1, point.Y - point.X
}

// 给定两点计算一般式方程
// 返回0, 0, 0表示两点相同构不成方程
// (-2, -2), (-2, -2) => 0, 0, 0
// (-2, 3), (-2, 6) => x = -2 => x + 2 = 0
// (-2, 3), (-1, 3) => y = 3 => y - 3 = 0
// (-2, 3), (2, 6) => -3x + 4y - 18 = 0
func twoPointForm(p1, p2 Coordinate) (int, int, int) {
	if p1.X == p2.X && p1.Y == p2.Y {
		return 0, 0, 0
	}
	// x = a
	if p1.X == p2.X {
		return 1, 0, -p1.X
	}
	// y = a
	if p1.Y == p2.Y {
		return 0, 1, -p1.Y
	}
	// 一般式
	A := p1.Y - p2.Y
	B := p2.X - p1.X
	C := p1.X * p2.Y - p2.X * p1.Y
	return A, B, C
}

// 解二元一次方程组, 一般式, 根据克勒默法则求解
// 返回err表示方程无解或有无数多个解
func systemOfBinaryLinearQquationsSolving(gp1, gp2 GeneralParameter) (x, y float64, err error) {
	D := gp1.A * gp2.B - gp2.A * gp1.B
	Dx := gp1.B * gp2.C - gp1.C * gp2.B
	Dy := gp1.C * gp2.A - gp2.C * gp1.A
	if D != 0 {
		// 方程有唯一解
		x = float64(Dx) / float64(D)
		y = float64(Dy) / float64(D)
		err = nil
		return
	}
	if Dx != 0 || Dy != 0 {
		// 方程组无解, 两直线平行
		err = errors.New("Equations without solution!")
		return
	}
	// Dx == 0 且 Dy == 0, 两直线重合, 方程有无数个解
	err = errors.New("There are countless solutions to the equation!")
	return
}

// 判断点p是否在gp构成的直线上
// 在返回true, 否则返回false
func lineContain(gp GeneralParameter, p Coordinate) bool {
	return gp.A * p.X + gp.B * p.Y + gp.C == 0
}

// 判断点是否在多边形内
// pointSet: 按照切片的顺序依次连成闭合多边形
// 比如N = 4, 则0->1, 1->2, 2->3, 3->0, 依次连成4个边
// p: 要判断多边形是否包含该点
// p在多边形内返回1, p在多边形边上返回0, 在多边形外返回-1; 同时err为nil
// 当传入错误的数据时, 会返回err
func PolygonContain(pointSet []Coordinate, p Coordinate) (int, error) {
	N := len(pointSet)
	// N个直线方程
	equationSet := make([]GeneralParameter, N)
	for i := 0; i < N - 1; i++ {
		A, B, C := twoPointForm(pointSet[i], pointSet[i + 1])
		if A == 0 && B == 0 && C == 0 {
			return 0, errors.New("Coincident point!")
		}
		equationSet[i] = GeneralParameter{A: A, B: B, C: C}
	}
	A, B, C := twoPointForm(pointSet[N - 1], pointSet[0])
	if A == 0 && B == 0 && C == 0 {
		return 0, errors.New("Coincident point!")
	}
	equationSet[N - 1] = GeneralParameter{A: A, B: B, C: C}

	// 点斜式方程
	pA, pB, pC := pointSlopeForm(p)
	pointSlopeParam := GeneralParameter{A: pA, B: pB, C: pC}
	// 计算点和N个边的交点个数
	numberIntersections := 0
	for i := 0; i < N; i++ {
		// 线段
		lineSeg1 := pointSet[i]
		var lineSeg2 Coordinate
		if i < N - 1 {
			lineSeg2 = pointSet[i + 1]
		} else {
			lineSeg2 = pointSet[0]
		}

		// 边端点的坐标范围
		minX, minY := lineSeg1.X, lineSeg1.Y
		maxX, maxY := lineSeg2.X, lineSeg2.Y
		if minX > maxX {
			minX, maxX = maxX, minX
		}
		if minY > maxY {
			minY, maxY = maxY, minY
		}

		// 判断点是否在多边形边上
		if lineContain(equationSet[i], p) {
			if p.X >= minX && p.X <= maxX && p.Y >= minY && p.Y <= maxY {
				return 0, nil
			}
		}

		// 判断交点是否在射线上
		// 这样的交点算作1个
		if lineContain(pointSlopeParam, pointSet[i]) && pointSet[i].X >= p.X {
			numberIntersections++
			continue
		}

		ix, iy, err := systemOfBinaryLinearQquationsSolving(pointSlopeParam, equationSet[i])
		if err != nil {
			continue
		}
		if ix < float64(p.X) || iy < float64(p.Y) {
			continue
		}

		// 竖直边
		if minX == maxX {
			if iy > float64(minY) && iy < float64(maxY) {
				numberIntersections++
			}
			continue
		}
		// 水平边
		if minY == maxY {
			if ix > float64(minX) && ix < float64(maxX) {
				numberIntersections++
			}
			continue
		}
		// 普通边
		if ix > float64(minX) && ix < float64(maxX) && iy > float64(minY) && iy < float64(maxY) {
			numberIntersections++
		}
	}

	// 内部
	if numberIntersections % 2 != 0 {
		return 1, nil
	}
	// 外部
	return -1, nil
}

