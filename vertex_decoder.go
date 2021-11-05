package imdl

import "image/color"

type VertexDecoder struct {
	data     []byte
	curIndex int
}

func (b *VertexDecoder) Decode8(index int) (int, byte) {
	return index + 1, b.data[index]
}

func (b *VertexDecoder) Decodeu16(index int) (int, uint16) {
	var bit0 byte
	var bit1 byte

	index, bit0 = b.Decode8(index)
	index, bit1 = b.Decode8(index)

	return index, uint16(bit0) | uint16(bit1)<<8
}

func (b *VertexDecoder) Decodeu32(index int) (int, uint32) {
	var bit0 uint16
	var bit1 uint16

	index, bit0 = b.Decodeu16(index)
	index, bit1 = b.Decodeu16(index)

	return index, uint32(bit0) | uint32(bit1)<<16
}

func (b *VertexDecoder) DecodeColor(index int) (int, color.RGBA) {
	var c color.RGBA

	index, c.R = b.Decode8(index)
	index, c.G = b.Decode8(index)
	index, c.B = b.Decode8(index)
	index, c.A = b.Decode8(index)

	return index, c
}

type SimpleDecoder struct {
	VertexDecoder
}

func (d *SimpleDecoder) DecodeQuantizedPosition(index int, pos *[3]uint16) int {
	index, pos[0] = d.Decodeu16(index)
	index, pos[1] = d.Decodeu16(index)
	index, pos[2] = d.Decodeu16(index)
	return index
}

func (d *SimpleDecoder) DecodeColorIndex(index int, colorIndex *uint16) int {
	if colorIndex != nil {
		index, *colorIndex = d.Decodeu16(index)
	} else {
		index += 2
	}
	return index
}

func (d *SimpleDecoder) DecodeFeatureIndex(index int, featureIndex *uint32) int {
	if featureIndex != nil {
		index, *featureIndex = d.Decodeu32(index)
	} else {
		index += 4
	}
	return index
}

func (d *SimpleDecoder) HasNext() bool {
	if d.curIndex < 0 {
		return false
	}
	return d.curIndex < len(d.data)
}

func (d *SimpleDecoder) VertexCount() int {
	return len(d.data) / 12
}

func (d *SimpleDecoder) Next() *SimpleVertex {
	if d.curIndex < 0 {
		return nil
	}
	curIndex := d.curIndex
	v := &SimpleVertex{ColorIndex: new(uint16), FeatureIndex: new(uint32)}
	curIndex = d.DecodeQuantizedPosition(curIndex, &v.QPos)
	curIndex = d.DecodeColorIndex(curIndex, v.ColorIndex)
	curIndex = d.DecodeFeatureIndex(curIndex, v.FeatureIndex)
	if curIndex < len(d.data) {
		d.curIndex = curIndex
	} else {
		d.curIndex = -1
	}
	return v
}

type SimplePolylineDecoder struct {
	SimpleDecoder
}

type MeshDecoder interface {
	Next() *MeshVertex
	HasNext() bool
	VertexCount() int
}

type SimpleMeshDecoder struct {
	SimpleDecoder
}

func (d *SimpleMeshDecoder) Next() *MeshVertex {
	vv := d.SimpleDecoder.Next()
	if vv != nil {
		return &MeshVertex{SimpleVertex: *vv}
	}
	return nil
}

type BaseMeshDecoder struct {
	SimpleMeshDecoder
}

func (b *BaseMeshDecoder) DecodeUV(index int, uv *[2]uint16) int {
	if uv != nil {
		index, uv[0] = b.Decodeu16(index)
		index, uv[1] = b.Decodeu16(index)
	} else {
		index += 4
	}
	return index
}

type TexturedMeshDecoder struct {
	BaseMeshDecoder
}

func (d *TexturedMeshDecoder) VertexCount() int {
	return len(d.data) / 16
}

func (d *TexturedMeshDecoder) Next() *MeshVertex {
	if d.curIndex < 0 {
		return nil
	}
	curIndex := d.curIndex
	v := &MeshVertex{SimpleVertex: SimpleVertex{FeatureIndex: new(uint32)}, QUV: &[2]uint16{}}
	curIndex = d.DecodeQuantizedPosition(curIndex, &v.QPos)
	curIndex += 2
	curIndex = d.DecodeFeatureIndex(curIndex, v.FeatureIndex)
	curIndex = d.DecodeUV(curIndex, v.QUV)
	if curIndex < len(d.data) {
		d.curIndex = curIndex
	} else {
		d.curIndex = -1
	}
	return v
}

type TexturedLitMeshDecoder struct {
	TexturedMeshDecoder
}

func (b *TexturedLitMeshDecoder) DecodeNormal(index int, normal *uint16) int {
	if normal != nil {
		index, *normal = b.Decodeu16(index)
	} else {
		index += 2
	}
	return index
}

func (d *TexturedLitMeshDecoder) VertexCount() int {
	return len(d.data) / 16
}

func (d *TexturedLitMeshDecoder) Next() *MeshVertex {
	if d.curIndex < 0 {
		return nil
	}
	curIndex := d.curIndex
	v := &MeshVertex{SimpleVertex: SimpleVertex{FeatureIndex: new(uint32)}, OctEncodedNormal: new(uint16), QUV: &[2]uint16{}}
	curIndex = d.DecodeQuantizedPosition(curIndex, &v.QPos)
	curIndex = d.DecodeNormal(curIndex, v.OctEncodedNormal)
	curIndex = d.DecodeFeatureIndex(curIndex, v.FeatureIndex)
	curIndex = d.DecodeUV(curIndex, v.QUV)
	if curIndex < len(d.data) {
		d.curIndex = curIndex
	} else {
		d.curIndex = -1
	}
	return v
}

type LitMeshDecoder struct {
	BaseMeshDecoder
}

func (d *LitMeshDecoder) VertexCount() int {
	return len(d.data) / 16
}

func (d *LitMeshDecoder) Next() *MeshVertex {
	if d.curIndex < 0 {
		return nil
	}
	curIndex := d.curIndex
	v := &MeshVertex{SimpleVertex: SimpleVertex{ColorIndex: new(uint16), FeatureIndex: new(uint32)}, OctEncodedNormal: new(uint16)}
	curIndex = d.DecodeQuantizedPosition(curIndex, &v.QPos)
	curIndex = d.DecodeColorIndex(curIndex, v.ColorIndex)
	curIndex = d.DecodeFeatureIndex(curIndex, v.FeatureIndex)
	curIndex, *v.OctEncodedNormal = d.Decodeu16(curIndex)
	curIndex += 2
	if curIndex < len(d.data) {
		d.curIndex = curIndex
	} else {
		d.curIndex = -1
	}
	return v
}

func CreateMeshDecoder(st SurfaceType, data []byte) MeshDecoder {
	switch st {
	case ST_Unlit:
		d := &BaseMeshDecoder{}
		d.curIndex = 0
		d.data = data
		return d
	case ST_Lit:
		d := &LitMeshDecoder{}
		d.curIndex = 0
		d.data = data
		return d
	case ST_Textured:
		d := &TexturedMeshDecoder{}
		d.curIndex = 0
		d.data = data
		return d
	case ST_TexturedLit:
		d := &TexturedLitMeshDecoder{}
		d.curIndex = 0
		d.data = data
		return d
	case ST_VolumeClassifier:
		d := &BaseMeshDecoder{}
		d.curIndex = 0
		d.data = data
		return d
	}
	return nil
}
