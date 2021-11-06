package imdl

import "math"

const rangeScale16 = 0xffff
const rangeScale8 = 0xff

func computeScale(extent float32, rangeScale uint16) float32 {
	if 0.0 == extent {
		return extent
	}
	return float32(rangeScale) / extent
}

func isInRange(qpos uint16, rangeScale uint16) bool {
	return qpos >= 0 && qpos < rangeScale+1
}

func Quantize(pos float32, origin float32, scale float32, rangeScale uint16) uint16 {
	return uint16(math.Floor(math.Max(0.0, math.Min(float64(rangeScale), float64(0.5+(pos-origin)*scale)))))
}

func IsQuantizable(pos float32, origin float32, scale float32, rangeScale uint16) bool {
	return isInRange(Quantize(pos, origin, scale, rangeScale), rangeScale16)
}

func UnQuantize(qpos uint16, origin float32, scale float32) float32 {
	if 0.0 == scale {
		return origin
	}
	return origin + float32(qpos)/scale
}

func IsQuantized(qpos uint16) bool {
	return isInRange(qpos, rangeScale16) && qpos == uint16(math.Floor(float64(qpos)))
}

type QParams2d struct {
	Origin [2]float32
	Scale  [2]float32
}

func (p *QParams2d) SetFromRange(range_ *Range2d, rangeScale uint16) {
	p.Origin[0] = range_.Low[0]
	p.Origin[1] = range_.Low[1]
	p.Scale[0] = computeScale(range_.High[0]-range_.Low[0], rangeScale)
	p.Scale[1] = computeScale(range_.High[1]-range_.Low[1], rangeScale)
}

func (p *QParams2d) rangeDiagonal() [2]float32 {
	var x float32
	var y float32

	if p.Scale[0] == 0 {
		x = 0
	} else {
		x = rangeScale16 / p.Scale[0]
	}

	if p.Scale[1] == 0 {
		y = 0
	} else {
		y = rangeScale16 / p.Scale[1]
	}

	return [2]float32{x, y}
}

func (p *QParams2d) GetRange() *Range2d {
	r := &Range2d{}
	r.Extend(p.Origin)
	d := p.rangeDiagonal()
	r.ExtendXY(p.Origin[0]+d[0], p.Origin[1]+d[1])
	return r
}

func QuantizePoint2d(pos [2]float32, params *QParams2d) [2]uint16 {
	var out [2]uint16
	out[0] = Quantize(pos[0], params.Origin[0], params.Scale[0], rangeScale16)
	out[1] = Quantize(pos[1], params.Origin[1], params.Scale[1], rangeScale16)
	return out
}

func UnQuantizePoint2d(qpos [2]uint16, params *QParams2d) [2]float32 {
	var out [2]float32
	out[0] = UnQuantize(qpos[0], params.Origin[0], params.Scale[0])
	out[1] = UnQuantize(qpos[1], params.Origin[1], params.Scale[1])
	return out
}

type QParams3d struct {
	Origin [3]float32
	Scale  [3]float32
}

func (p *QParams3d) SetFromRange(range_ *Range3d, rangeScale uint16) {
	p.Origin[0] = range_.Low[0]
	p.Origin[1] = range_.Low[1]
	p.Origin[2] = range_.Low[2]

	p.Scale[0] = computeScale(range_.High[0]-range_.Low[0], rangeScale)
	p.Scale[1] = computeScale(range_.High[1]-range_.Low[1], rangeScale)
	p.Scale[2] = computeScale(range_.High[2]-range_.Low[2], rangeScale)
}

func (p *QParams3d) rangeDiagonal() [3]float32 {
	var x float32
	var y float32
	var z float32

	if p.Scale[0] == 0 {
		x = 0
	} else {
		x = rangeScale16 / p.Scale[0]
	}

	if p.Scale[1] == 0 {
		y = 0
	} else {
		y = rangeScale16 / p.Scale[1]
	}

	if p.Scale[2] == 0 {
		z = 0
	} else {
		z = rangeScale16 / p.Scale[2]
	}

	return [3]float32{x, y, z}
}

func (p *QParams3d) GetRange() *Range3d {
	r := &Range3d{}
	r.Extend(p.Origin)
	d := p.rangeDiagonal()
	r.ExtendXYZ(p.Origin[0]+d[0], p.Origin[1]+d[1], p.Origin[2]+d[2])
	return r
}

func QuantizePoint3d(pos [3]float32, params *QParams3d) [3]uint16 {
	var out [3]uint16
	out[0] = Quantize(pos[0], params.Origin[0], params.Scale[0], rangeScale16)
	out[1] = Quantize(pos[1], params.Origin[1], params.Scale[1], rangeScale16)
	out[2] = Quantize(pos[2], params.Origin[2], params.Scale[2], rangeScale16)
	return out
}

func UnQuantizePoint3d(qpos [3]uint16, params *QParams3d) [3]float32 {
	var out [3]float32
	out[0] = UnQuantize(qpos[0], params.Origin[0], params.Scale[0])
	out[1] = UnQuantize(qpos[1], params.Origin[1], params.Scale[1])
	out[2] = UnQuantize(qpos[2], params.Origin[2], params.Scale[2])
	return out
}
