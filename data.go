package imdl

type MeshVertex struct {
	SimpleVertex
	UV      *[2]uint16
	Normals *uint16
}

type MeshData struct {
	Type    SurfaceType
	Indices []uint32
	Vertexs []MeshVertex
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
	Pos          [3]uint16
	ColorIndex   *uint16
	FeatureIndex *uint32
}

type PolylineData struct {
	Indices []uint32
	Vertexs []SimpleVertex
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
