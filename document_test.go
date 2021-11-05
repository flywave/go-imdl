package imdl

import "testing"

func TestEncodeVertexIndices(t *testing.T) {
	indices := make([]uint32, 10)

	for i := 0; i < 10; i++ {
		indices[i] = uint32(i)
	}

	data := EncodeVertexIndices(indices)

	if data == nil {
		t.FailNow()
	}

	oindices := DecodeVertexIndices(data)

	if oindices == nil {
		t.FailNow()
	}

	for i := range indices {
		if indices[i] != oindices[i] {
			t.FailNow()
		}
	}
}

func TestMeshBuilder(t *testing.T) {
	vers := make([]MeshVertex, 10)

	for i := range vers {
		find := uint32(i)
		vers[i] = MeshVertex{SimpleVertex: SimpleVertex{Pos: [3]uint16{uint16(i), uint16(i), uint16(i)}, FeatureIndex: &find}}
	}

	builder := newSimpleMeshBuilder(10)

	builder.Process(vers)

	data := builder.data

	decoder := &SimpleMeshDecoder{}
	decoder.data = data
	decoder.curIndex = 0

	if decoder.VertexCount() != 10 {
		t.FailNow()
	}

	dver := make([]MeshVertex, 10)

	j := 0
	for decoder.HasNext() {
		dver[j] = *decoder.Next()
		j++
	}

	for i := range vers {
		if vers[i].FeatureIndex == nil || dver[i].FeatureIndex == nil || *vers[i].FeatureIndex != *dver[i].FeatureIndex {
			t.FailNow()
		}
	}
}

func TestTexturedMeshDecoder(t *testing.T) {
	vers := make([]MeshVertex, 10)

	for i := range vers {
		find := uint32(i)
		vers[i] = MeshVertex{SimpleVertex: SimpleVertex{Pos: [3]uint16{uint16(i), uint16(i), uint16(i)}, FeatureIndex: &find}, UV: &[2]uint16{uint16(i), uint16(i)}}
	}

	builder := newTexturedMeshBuilder(10)

	builder.Process(vers)

	data := builder.data

	decoder := &TexturedMeshDecoder{}
	decoder.data = data
	decoder.curIndex = 0

	if decoder.VertexCount() != 10 {
		t.FailNow()
	}

	dver := make([]MeshVertex, 10)

	j := 0
	for decoder.HasNext() {
		dver[j] = *decoder.Next()
		j++
	}

	for i := range vers {
		if vers[i].UV == nil || dver[i].UV == nil || *vers[i].UV != *dver[i].UV {
			t.FailNow()
		}
	}
}
