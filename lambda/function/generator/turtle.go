package generator

import (
	"math"
	"phoenix/lambda/function"
)

type QuickSave struct {
	Pos   function.Vector
	Block struct {
		name string
		data uint8
	}
	Matrix [3][3]float64
	Pen    bool
}

type Turtle struct {
	Pos   function.Vector
	Space *function.Space
	Pen   bool
	Block struct {
		name string
		data uint8
	}
	Stack       []QuickSave
	Orientation float64
	PitchValue  float64
	Matrix      [3][3]float64
}

var TO_RADIANS = math.Pi / 180.
var TO_DEGREES = 180. / math.Pi

func NewTurtle(space *function.Space, blockName string, blockData uint8) *Turtle {
	return &Turtle{
		Pos:   space.GetPointer(),
		Space: space,
		Pen:   false,
		Block: struct {
			name string
			data uint8
		}{},
		Orientation: 0,
	}
}

func (t *Turtle) Push() {
	t.Stack = append(t.Stack, t.Save())
}

func (t *Turtle) Save() QuickSave {
	return QuickSave{
		Pos:    t.Pos,
		Block:  t.Block,
		Matrix: t.Matrix,
		Pen:    t.Pen,
	}
}

func (t *Turtle) Restore(save QuickSave) {
	t.Block = save.Block
	t.Matrix = save.Matrix
	t.Pen = save.Pen
	t.Pos = save.Pos
}

func (t *Turtle) Pop() {
	save := t.Stack[len(t.Stack)-1]
	t.Stack = t.Stack[0 : len(t.Stack)-1]
	t.Restore(save)
}

func (t *Turtle) Pitch(angle float64) {
	t.Matrix = matrixMultiply(t.Matrix, pitchMatrix(angle))
	t.DirectionOut()
}

func (t *Turtle) Forward(dist float64) {
	d := t.GetHeading()
	newX := t.Pos[0] + d[0]*dist
	newY := t.Pos[1] + d[1]*dist
	newZ := t.Pos[2] + d[2]*dist
	line := getLine(function.Vector3f{t.Pos[0], t.Pos[1], t.Pos[2]}, function.Vector3f{newX, newY, newZ})
	if t.Pen {
		for _, b := range line {
			t.Space.Plot(b)
		}
	}
	t.Pos = function.Vector{newX, newY, newZ}
}

func (t *Turtle) RollAngle(angle float64) {
	angles := t.GetAngle()
	m0 := matrixMultiply(yawMatrix(angles[0]), pitchMatrix(angles[1]))
	t.Matrix = matrixMultiply(m0, rollMatrix(angle))
}

func (t *Turtle) Yaw(angle float64) {
	t.Matrix = matrixMultiply(t.Matrix, yawMatrix(angle))
	t.DirectionOut()
}

func (t *Turtle) DirectionOut() {
	heading := t.GetHeading()
	xz := math.Sqrt(heading[0]*heading[0] + heading[1]*heading[1])
	pitch := math.Atan2(-heading[1], xz) * TO_DEGREES
	t.SetPitch(pitch)
	if xz >= 1e-9 {
		rot := math.Atan2(-heading[0], heading[2]) * TO_DEGREES
		t.SetRotation(rot)
	}
}

func (t *Turtle) GetHeading() [3]float64 {
	return [3]float64{t.Matrix[0][2], t.Matrix[1][2], t.Matrix[2][2]}
}

func (t *Turtle) GetAngle() [2]float64 {
	heading := t.GetHeading()
	rot := .0
	pitch := .0
	if heading[0] == math.Trunc(heading[0]) && heading[1] == math.Trunc(heading[1]) && heading[2] == math.Trunc(heading[2]) {
		xz := math.Abs(heading[0]) + math.Abs(heading[2])
		if xz != 0 {
			rot = iAtan2(-int64(heading[0]), int64(heading[2]))
		}
		pitch = iAtan2(-int64(heading[1]), int64(xz))
	} else {
		xz := math.Sqrt(heading[0]*heading[0] + heading[1]*heading[1])
		if xz >= 1e-9 {
			rot = math.Atan2(-heading[0], heading[2]) * TO_RADIANS
		}
		pitch = math.Atan2(-heading[1], xz) * TO_RADIANS
	}
	return [2]float64{rot, pitch}
}

func (t *Turtle) Backward(dist float64) {
	t.Forward(-dist)
}

func (t *Turtle) Right(angle float64) {
	t.Matrix = matrixMultiply(yawMatrix(angle), t.Matrix)
	t.DirectionOut()
}

func (t *Turtle) Left(angle float64) {
	t.Right(-angle)
}

func (t *Turtle) Up(angle float64) {
	t.Pitch(angle)
}

func (t *Turtle) Down(angle float64) {
	t.Up(-angle)
}

// SetAngle Set roll angle of turtle (compass, vertical, roll) in degrees
func (t *Turtle) SetAngle(compass, vertical, roll float64) {
	t.Matrix = makeMatrix(compass, vertical, roll)
}

func (t *Turtle) PenUp() {
	t.Pen = false
}

func (t *Turtle) PenDoown() {
	t.Pen = true
}

func (t *Turtle) Goto(v function.Vector) {
	t.Pos = v
	t.DirectionOut()
}

func (t *Turtle) VerticalAngle(angle float64) {
	angles := t.GetAngle()
	t.Matrix = matrixMultiply(yawMatrix(angles[0]), pitchMatrix(angle))
	t.DirectionOut()
}

// GridAlign Align positions to grid
func (t *Turtle) GridAlign() {
	bestDist := 2 * 9.
	bestMatrix := makeMatrix(0, 0, 0)
	for _, compass := range []float64{0, 90, 180, 270} {
		for _, pitch := range []float64{0, 90, 180, 270} {
			for _, roll := range []float64{0, 90, 180, 270} {
				m := makeMatrix(compass, pitch, roll)
				dist := matrixDistanceSquared(t.Matrix, m)
				if dist < bestDist {
					bestMatrix = m
					bestDist = dist
				}
			}
		}
	}
	t.Matrix = bestMatrix
	t.DirectionOut()
}

func (t *Turtle) Roll(angle float64) {
	t.Matrix = matrixMultiply(t.Matrix, rollMatrix(angle))
	t.DirectionOut()
}

func (t *Turtle) SetPitch(p float64) {
	t.PitchValue = p
}

func (t *Turtle) SetRotation(r float64) {
	t.Orientation = r
}

func yawMatrix(angle float64) [3][3]float64 {
	if angle == math.Trunc(angle) && int64(angle)%90 == 0 {
		angleI := int64(math.Trunc(angle))
		return [3][3]float64{
			{iCos(angleI), 0, -iSin(angleI)},
			{0, 1, 0},
			{iSin(angleI), 0, iCos(angleI)},
		}
	} else {
		theta := angle * TO_RADIANS
		return [3][3]float64{
			{math.Cos(theta), 0., -math.Sin(theta)},
			{0, 1, 0},
			{math.Sin(theta), 0, math.Cos(theta)},
		}
	}
}

func pitchMatrix(angle float64) [3][3]float64 {
	if angle == math.Trunc(angle) && int64(angle)%90 == 0 {
		angleI := int64(math.Trunc(angle))
		return [3][3]float64{
			{1, 0, 0},
			{0, iCos(angleI), iSin(angleI)},
			{0, -iSin(angleI), iCos(angleI)},
		}
	} else {
		theta := angle * TO_RADIANS
		return [3][3]float64{
			{1, 0, 0},
			{0, math.Cos(theta), math.Sin(theta)},
			{0, -math.Sin(theta), math.Cos(theta)},
		}
	}
}

func rollMatrix(angle float64) [3][3]float64 {
	if angle == math.Trunc(angle) && int64(angle)%90 == 0 {
		angleI := int64(math.Trunc(angle))
		return [3][3]float64{
			{iCos(angleI), -iSin(angleI), 0},
			{iSin(angleI), iCos(angleI), 0},
			{0, 0, 1},
		}
	} else {
		theta := angle * TO_RADIANS
		return [3][3]float64{
			{math.Cos(theta), -math.Sin(theta), 0},
			{math.Sin(theta), math.Cos(theta), 0},
			{0, 0, 1},
		}
	}
}

func matrixMultiply(a [3][3]float64, b [3][3]float64) [3][3]float64 {
	var c [3][3]float64
	for i := 0; i < 3; i++ {
		for j := 0; i < 3; j++ {
			c[i][j] = a[i][0]*b[0][j] + a[i][1]*b[1][j] + a[i][2]*b[2][j]
		}
	}
	return c
}

var ICOS = [4]float64{1, 0, -1, 0}
var ISIN = [4]float64{0, 1, 0, -1}

func iCos(angle int64) float64 {
	return ICOS[(angle%360)/90]
}

func iSin(angle int64) float64 {
	return ISIN[(angle%360)/90]
}

func iAtan2(y int64, x int64) float64 {
	if x == 0 {
		if y > 0 {
			return 90
		} else {
			return -90
		}
	} else {
		if x > 0 {
			return 0
		} else {
			return 180
		}
	}
}

func getLine(begin, end function.Vector3f) []function.Vector {
	var BlockSet []function.Vector
	sx, sy, sz := begin[0], begin[1], begin[2]
	ex, ey, ez := end[0], end[1], end[2]
	i, j, k := sx, sy, sz
	t := 0.0
	s := 1 / math.Sqrt(math.Pow(ex-sx, 2)+math.Pow(ey-sy, 2)+math.Pow(ez-sz, 2))
	for t >= 0 && t <= 1 {
		i = t*(ex-sx) + sx
		j = t*(ey-sy) + sy
		k = t*(ez-sz) + sz
		t += s
		BlockSet = append(BlockSet, function.Vector{i, j, k})
	}
	return BlockSet
}

func makeMatrix(compass, vertical, roll float64) [3][3]float64 {
	m0 := matrixMultiply(yawMatrix(compass), pitchMatrix(vertical))
	return matrixMultiply(m0, rollMatrix(roll))
}

func matrixDistanceSquared(m1, m2 [3][3]float64) float64 {
	d2 := 0.
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			d2 += math.Pow(m1[i][j]-m2[i][j], 2)
		}
	}
	return d2
}
