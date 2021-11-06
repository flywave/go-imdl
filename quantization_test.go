package imdl

import "testing"

func TestQuantize(t *testing.T) {
	range_ := CreateRange3d([][3]float32{{0, -100, 200}, {50, 100, 10000}})
	qparams := &QParams3d{}
	qparams.SetFromRange(range_, rangeScale16)

	if qparams.Origin[0] != 0 {
		t.FailNow()
	}

	qpos := QuantizePoint3d([3]float32{20, 0, 800}, qparams)

	if !IsQuantized(qpos[0]) {
		t.FailNow()
	}
}
