package main

import (
	"fmt"
	"time"

	manor "github.com/monchickey/manor-go/v2"
)

type Coordinate manor.Coordinate

func PolygonContain(pointSet []Coordinate, p Coordinate) (int, error) {
	newPointSet := make([]manor.Coordinate, len(pointSet))
	for i, c := range pointSet {
		newPointSet[i] = manor.Coordinate(c)
	}
	return manor.PolygonContain(newPointSet, manor.Coordinate(p))
}

func main() {
	nowTimestamp := time.Now().Unix()
	nowTimeStr := manor.TimestampToString(nowTimestamp, "2006-01-02 15:04:05")
	fmt.Println(nowTimeStr)

	numSeq := []uint8{72, 101, 108, 108, 111, 32, 109, 111, 110, 99, 104, 105, 99, 107, 101, 121, 33}
	raw := manor.Uint8ToBytes(numSeq)
	fmt.Println(string(raw))
	fmt.Println(manor.Base64Encode(raw))

	geoHash, err := manor.GeohashEncode(113.56291, 36.9271, 12)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("(113.56291, 36.9271) geohash Encode:", geoHash)
	}
	longitude, latitude, err := manor.GeohashDecode(geoHash)
	if err == nil {
		fmt.Println(geoHash, "Decode:(", longitude, latitude, ")")
	}

	pointSet := []Coordinate{
		Coordinate{1, 1},
		Coordinate{1, 4},
		Coordinate{4, 4},
		Coordinate{4, 1},
	}

	fmt.Println("Polygon: (1,1)-(1,4)-(4,4)-(4,1): ")
	v, _ := PolygonContain(pointSet, Coordinate{1, 1})
	fmt.Println("  (1, 1) in", v) // 边上
	v, _ = PolygonContain(pointSet, Coordinate{2, 2})
	fmt.Println("  (2, 2) in", v) // 内部
	v, _ = PolygonContain(pointSet, Coordinate{5, 1})
	fmt.Println("  (5, 1) in", v) // 外部
}
