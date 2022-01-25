package fetcher

import (
	"fmt"
	"phoenix/lambda/function"
	//"phoenix/ligo"
	"math"
)

type Area struct {
	start  function.Vector
	offset function.Vector
}

const Unit = 100.0

func GetOblong(begin, end function.Vector) []Area {
	// Element[0]: Start
	// Element[1]: End
	var area []Area
	minX := math.Min(begin[0], end[0])
	minY := math.Min(begin[1], end[1])
	minZ := math.Min(begin[2], end[2])

	offsetX := math.Abs(begin[0] - end[0])
	offsetY := math.Abs(begin[1] - end[1])
	offsetZ := math.Abs(begin[2] - end[2])

	splitX := SplitLen(offsetX)
	splitZ := SplitLen(offsetZ)

	for x := .0; x < splitX; x++ {
		for z := .0; z < splitZ; z++ {
			start := function.Vector{
				minX + x*Unit,
				minY,
				minZ + z*Unit,
			}
			oz := 0.0
			if z == splitZ-1 {
				oz = offsetZ - (splitZ-1)*Unit
				fmt.Println(oz)
			} else {
				oz = Unit
			}
			ox := 0.0
			if x == splitX-1 {
				ox = offsetX - (splitX-1)*Unit
			} else {
				ox = Unit
			}
			offset := function.Vector{
				ox,
				offsetY,
				oz,
			}
			area = append(
				area, Area{
					start:  start,
					offset: offset,
				},
			)
		}

	}
	return area
}

func SplitLen(length float64) float64 {
	curr := 1.0
	for length/curr > Unit {
		curr += 1.0
	}
	return curr
}
