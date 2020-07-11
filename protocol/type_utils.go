package protocol

import "math"

func NewPosition(x, y, z int) Position {
    return Position(((x & 0x3ffffff) << 38) | ((y & 0xFFF) << 26) | (z & 0x3ffffff))
}

func (p Position) X() int {
    return int((int64(p) >> 38) & 0x3ffffff)
}

func (p Position) Y() int {
    return int((int64(p) >> 26) & 0xfff)
}

func (p Position) Z() int {
    return int(int64(p) & 0x3ffffff)
}

func AngleFromDegrees(deg float64) Angle {
    return Angle(deg * 256 / 360)
}

func AngleFromRadians(rad float64) Angle {
    return Angle(rad * 256 / (2 * math.Pi))
}

func (a Angle) Degrees() float64 {
    return float64(a) * 360 / 256
}

func (a Angle) Radians() float64 {
    return float64(a) * 2 * math.Pi / 256
}