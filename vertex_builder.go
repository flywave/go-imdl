package imdl

import (
	"image/color"
	"math"
)

type VertexBuilder struct {
	data     []byte
	curIndex int
}

func (b *VertexBuilder) Advance(nBytes int) bool {
	b.curIndex += nBytes
	return b.curIndex <= len(b.data)
}

func (b *VertexBuilder) Append8(val byte) {
	b.data[b.curIndex] = byte(val)
	b.Advance(1)
}

func (b *VertexBuilder) Append16(val uint16) {
	b.Append8(uint8(val & 0x00ff))
	b.Append8(uint8(val >> 8))
}

func (b *VertexBuilder) Append32(val uint32) {
	b.Append16(uint16(val & 0x0000ffff))
	b.Append16(uint16(val >> 16))
}

func (b *VertexBuilder) AppendColor(color color.NRGBA) {
	color.A = 255 - color.A

	switch color.A {
	case 0:
		color.R = 0
		color.G = 0
		color.B = 0
	case 255:
		break
	default:
		{
			f := float64(color.A) / 255.0
			color.R = uint8(math.Floor(float64(color.R)*f + 0.5))
			color.G = uint8(math.Floor(float64(color.G)*f + 0.5))
			color.B = uint8(math.Floor(float64(color.B)*f + 0.5))
		}
	}

	b.Append8(color.R)
	b.Append8(color.G)
	b.Append8(color.B)
	b.Append8(color.A)
}

/**
 *  Each vertex consists of 12 bytes:
 *  pos.x           00
 *  pos.y           02
 *  pos.z           04
 *  colorIndex      06
 *  featureIndex    08
 */
type SimpleBuilder struct {
	VertexBuilder
}

func newSimpleBuilder(vertexCount int) *SimpleBuilder {
	data := make([]byte, vertexCount*12)
	return &SimpleBuilder{VertexBuilder{data: data, curIndex: 0}}
}

func (b *SimpleBuilder) AppendVertex(v *SimpleVertex) {
	b.AppendQuantizedPosition(v.Pos)
	b.AppendColorIndex(v.ColorIndex)
	b.AppendFeatureIndex(v.FeatureIndex)
}

func (b *SimpleBuilder) AppendQuantizedPosition(pos [3]uint16) {
	b.Append16(pos[0])
	b.Append16(pos[1])
	b.Append16(pos[2])
}

func (b *SimpleBuilder) AppendColorIndex(colorIndex *uint16) {
	if colorIndex != nil {
		b.Append16(*colorIndex)
	} else {
		b.Advance(2)
	}
}

func (b *SimpleBuilder) AppendFeatureIndex(featureIndex *uint32) {
	if featureIndex != nil {
		b.Append32(*featureIndex)
	} else {
		b.Advance(4)
	}
}

type SimplePolylineBuilder struct {
	SimpleBuilder
}

func newSimplePolylineBuilder(vertexCount int) *SimplePolylineBuilder {
	data := make([]byte, vertexCount*12)
	return &SimplePolylineBuilder{SimpleBuilder: SimpleBuilder{VertexBuilder{data: data, curIndex: 0}}}
}

func (b *SimplePolylineBuilder) Process(vertexs []SimpleVertex) {
	for i := range vertexs {
		b.AppendVertex(&vertexs[i])
	}
}

type MeshBuilder interface {
	Advance(nBytes int) bool
	Process(vertexs []MeshVertex)
	GetData() []byte
}

type SimpleMeshBuilder struct {
	SimpleBuilder
}

func newSimpleMeshBuilder(vertexCount int) *SimpleMeshBuilder {
	data := make([]byte, vertexCount*12)
	return &SimpleMeshBuilder{SimpleBuilder: SimpleBuilder{VertexBuilder{data: data, curIndex: 0}}}
}

func (b *SimpleMeshBuilder) AppendVertex(v *MeshVertex) {
	b.SimpleBuilder.AppendVertex(&v.SimpleVertex)
}

func (b *SimpleMeshBuilder) Process(vertexs []MeshVertex) {
	for i := range vertexs {
		b.AppendVertex(&vertexs[i])
	}
}

func (b *SimpleMeshBuilder) GetData() []byte {
	return b.data
}

type BaseMeshBuilder struct {
	SimpleMeshBuilder
}

func newBaseMeshBuilder(vertexCount int) *BaseMeshBuilder {
	return &BaseMeshBuilder{SimpleMeshBuilder: *newSimpleMeshBuilder(vertexCount)}
}

func (b *BaseMeshBuilder) AppendUV(uv *[2]uint16) {
	if uv != nil {
		b.Append16(uv[0])
		b.Append16(uv[1])
	} else {
		b.Advance(4)
	}
}

func (b *BaseMeshBuilder) AppendVertex(v *MeshVertex) {
	b.SimpleBuilder.AppendVertex(&v.SimpleVertex)
}

func (b *BaseMeshBuilder) Process(vertexs []MeshVertex) {
	for i := range vertexs {
		b.AppendVertex(&vertexs[i])
	}
}

type TexturedMeshBuilder struct {
	BaseMeshBuilder
}

func newTexturedMeshBuilder(vertexCount int) *TexturedMeshBuilder {
	b := &TexturedMeshBuilder{}
	b.data = make([]byte, vertexCount*16)
	b.curIndex = 0
	return b
}

func (b *TexturedMeshBuilder) AppendVertex(v *MeshVertex) {
	b.AppendQuantizedPosition(v.Pos)
	b.Advance(2)
	b.AppendFeatureIndex(v.FeatureIndex)
	b.AppendUV(v.UV)
}

func (b *TexturedMeshBuilder) Process(vertexs []MeshVertex) {
	for i := range vertexs {
		b.AppendVertex(&vertexs[i])
	}
}

type TexturedLitMeshBuilder struct {
	TexturedMeshBuilder
}

func newTexturedLitMeshBuilder(vertexCount int) *TexturedLitMeshBuilder {
	b := &TexturedLitMeshBuilder{}
	b.data = make([]byte, vertexCount*16)
	b.curIndex = 0
	return b
}

func (b *TexturedLitMeshBuilder) AppendNormal(normal *uint16) {
	if normal != nil {
		b.Append16(*normal)
	} else {
		b.Advance(2)
	}
}
func (b *TexturedLitMeshBuilder) AppendVertex(v *MeshVertex) {
	b.AppendQuantizedPosition(v.Pos)
	b.AppendNormal(v.Normals)
	b.AppendFeatureIndex(v.FeatureIndex)
	b.AppendUV(v.UV)
}

func (b *TexturedLitMeshBuilder) Process(vertexs []MeshVertex) {
	for i := range vertexs {
		b.AppendVertex(&vertexs[i])
	}
}

type LitMeshBuilder struct {
	BaseMeshBuilder
}

func newLitMeshBuilder(vertexCount int) *LitMeshBuilder {
	b := &LitMeshBuilder{}
	b.data = make([]byte, vertexCount*16)
	b.curIndex = 0
	return b
}

func (b *LitMeshBuilder) AppendVertex(v *MeshVertex) {
	b.BaseMeshBuilder.SimpleBuilder.AppendVertex(&v.SimpleVertex)
	b.Append16(*v.Normals)
	b.Advance(2)
}

func (b *LitMeshBuilder) Process(vertexs []MeshVertex) {
	for i := range vertexs {
		b.AppendVertex(&vertexs[i])
	}
}

func CreateMeshBuilder(st SurfaceType, vertexCount int) MeshBuilder {
	switch st {
	case ST_Unlit:
		return newBaseMeshBuilder(vertexCount)
	case ST_Lit:
		return newLitMeshBuilder(vertexCount)
	case ST_Textured:
		return newTexturedMeshBuilder(vertexCount)
	case ST_TexturedLit:
		return newTexturedLitMeshBuilder(vertexCount)
	case ST_VolumeClassifier:
		return newBaseMeshBuilder(vertexCount)
	}
	return nil
}
