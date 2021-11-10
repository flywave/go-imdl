package imdl

import "math"

type MeshVertex struct {
	SimpleVertex
	UV               *[2]float32
	QUV              *[2]uint16
	Normal           *[3]float32
	OctEncodedNormal *uint16
}

type MeshData struct {
	Type    SurfaceType
	Indices []uint32
	Vertexs []MeshVertex
}

func (d *MeshData) GetPosRange() *Range3d {
	r := &Range3d{Low: [3]float32{math.MaxFloat32, math.MaxFloat32, math.MaxFloat32}, High: [3]float32{-math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32}}
	for i := range d.Vertexs {
		r.Extend(d.Vertexs[i].Pos)
	}
	return r
}

func (d *MeshData) GetPosQParams3d() *QParams3d {
	qparams := &QParams3d{}
	qparams.SetFromRange(d.GetPosRange(), rangeScale16)
	return qparams
}

func (d *MeshData) GetUvRange() *Range2d {
	r := &Range2d{Low: [2]float32{math.MaxFloat32, math.MaxFloat32}, High: [2]float32{-math.MaxFloat32, -math.MaxFloat32}}
	hasRange := false
	for i := range d.Vertexs {
		if d.Vertexs[i].UV != nil {
			r.Extend(*d.Vertexs[i].UV)
			hasRange = true
		}
	}
	if hasRange {
		return r
	}
	return nil
}

func (d *MeshData) GetUvQParams2d() *QParams2d {
	r := d.GetUvRange()
	if r != nil {
		qparams := &QParams2d{}
		qparams.SetFromRange(d.GetUvRange(), rangeScale16)
		return qparams
	}
	return nil
}

func (d *MeshData) Quantize() (*Range3d, *Range2d) {
	qParams3d := d.GetPosQParams3d()
	qParams2d := d.GetUvQParams2d()

	for i := range d.Vertexs {
		d.Vertexs[i].QPos = QuantizePoint3d(d.Vertexs[i].Pos, qParams3d)
		if qParams2d != nil {
			uv := QuantizePoint2d(*d.Vertexs[i].UV, qParams2d)
			d.Vertexs[i].QUV = &uv
		}
		if d.Vertexs[i].Normal != nil {
			qn := encodeXYZ(d.Vertexs[i].Normal[0], d.Vertexs[i].Normal[1], d.Vertexs[i].Normal[2])
			d.Vertexs[i].OctEncodedNormal = &qn
		}
	}

	if qParams2d == nil {
		return qParams3d.GetRange(), nil
	}

	return qParams3d.GetRange(), qParams2d.GetRange()
}

func (d *MeshData) UnQuantize(qParams3d *QParams3d, qParams2d *QParams2d) {
	for i := range d.Vertexs {
		d.Vertexs[i].Pos = UnQuantizePoint3d(d.Vertexs[i].QPos, qParams3d)
		if d.Vertexs[i].QUV != nil && qParams2d != nil {
			uv := UnQuantizePoint2d(*d.Vertexs[i].QUV, qParams2d)
			d.Vertexs[i].UV = &uv
		}
		if d.Vertexs[i].OctEncodedNormal != nil {
			norms := decodeValue(*d.Vertexs[i].OctEncodedNormal)
			d.Vertexs[i].Normal = &norms
		}
	}
}

func (d *MeshData) EncodeIndices() []byte {
	return EncodeVertexIndices(d.Indices)
}

func (d *MeshData) DecodeIndices(bytes []byte) {
	d.Indices = DecodeVertexIndices(bytes)
}

func (d *MeshData) EncodeVertexs() []byte {
	builder := CreateMeshBuilder(d.Type, len(d.Vertexs))
	builder.Process(d.Vertexs)
	return builder.GetData()
}

func (d *MeshData) DecodeVertexs(data []byte, count uint32) {
	decoder := CreateMeshDecoder(d.Type, data)
	d.Vertexs = make([]MeshVertex, count)

	for i := range d.Vertexs {
		d.Vertexs[i] = *decoder.Next()
	}
}

type SimpleVertex struct {
	Pos          [3]float32
	QPos         [3]uint16
	ColorIndex   *uint16
	FeatureIndex *uint32
}

type PolylineData struct {
	Indices []uint32
	Vertexs []SimpleVertex
}

func (d *PolylineData) GetRange() *Range3d {
	r := &Range3d{Low: [3]float32{math.MaxFloat32, math.MaxFloat32, math.MaxFloat32}, High: [3]float32{-math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32}}
	for i := range d.Vertexs {
		r.Extend(d.Vertexs[i].Pos)
	}
	return r
}

func (d *PolylineData) GetQParams3d() *QParams3d {
	qparams := &QParams3d{}
	qparams.SetFromRange(d.GetRange(), rangeScale16)
	return qparams
}

func (d *PolylineData) Quantize() *Range3d {
	qParams3d := d.GetQParams3d()
	for i := range d.Vertexs {
		d.Vertexs[i].QPos = QuantizePoint3d(d.Vertexs[i].Pos, qParams3d)
	}
	return qParams3d.GetRange()
}

func (d *PolylineData) UnQuantize(qParams3d *QParams3d) {
	for i := range d.Vertexs {
		d.Vertexs[i].Pos = UnQuantizePoint3d(d.Vertexs[i].QPos, qParams3d)
	}
}

func (d *PolylineData) EncodeIndices() []byte {
	return EncodeVertexIndices(d.Indices)
}

func (d *PolylineData) DecodeIndices(bytes []byte) {
	d.Indices = DecodeVertexIndices(bytes)
}

func (d *PolylineData) EncodeVertexs() []byte {
	builder := newSimplePolylineBuilder(len(d.Vertexs))
	builder.Process(d.Vertexs)
	return builder.data
}

func (d *PolylineData) DecodeVertexs(data []byte) {
	decoder := &SimplePolylineDecoder{}
	decoder.data = data
	decoder.curIndex = 0

	d.Vertexs = make([]SimpleVertex, decoder.VertexCount())
	j := 0
	for decoder.HasNext() {
		d.Vertexs[j] = *decoder.Next()
		j++
	}
}

type PointStringData struct {
	Indices []uint32
	Vertexs []SimpleVertex
}

func (d *PointStringData) GetRange() *Range3d {
	r := &Range3d{Low: [3]float32{math.MaxFloat32, math.MaxFloat32, math.MaxFloat32}, High: [3]float32{-math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32}}
	for i := range d.Vertexs {
		r.Extend(d.Vertexs[i].Pos)
	}
	return r
}

func (d *PointStringData) GetQParams3d() *QParams3d {
	qparams := &QParams3d{}
	qparams.SetFromRange(d.GetRange(), rangeScale16)
	return qparams
}

func (d *PointStringData) Quantize() *Range3d {
	qParams3d := d.GetQParams3d()
	for i := range d.Vertexs {
		d.Vertexs[i].QPos = QuantizePoint3d(d.Vertexs[i].Pos, qParams3d)
	}
	return qParams3d.GetRange()
}

func (d *PointStringData) UnQuantize(qParams3d *QParams3d) {
	for i := range d.Vertexs {
		d.Vertexs[i].Pos = UnQuantizePoint3d(d.Vertexs[i].QPos, qParams3d)
	}
}

func (d *PointStringData) EncodeIndices() []byte {
	return EncodeVertexIndices(d.Indices)
}

func (d *PointStringData) DecodeIndices(bytes []byte) {
	d.Indices = DecodeVertexIndices(bytes)
}

func (d *PointStringData) EncodeVertexs() []byte {
	builder := newSimplePolylineBuilder(len(d.Vertexs))
	builder.Process(d.Vertexs)
	return builder.data
}

func (d *PointStringData) DecodeVertexs(data []byte) {
	decoder := &SimplePolylineDecoder{}
	decoder.data = data
	decoder.curIndex = 0

	d.Vertexs = make([]SimpleVertex, decoder.VertexCount())
	j := 0
	for decoder.HasNext() {
		d.Vertexs[j] = *decoder.Next()
		j++
	}
}

type chunkData struct {
	name string
	data []byte
}
