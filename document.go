package imdl

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"

	"github.com/flywave/gltf"
)

const NODE_ROOT = "Node_Root"

type FeatureIndexType uint32

const (
	Empty      FeatureIndexType = 0
	NonUniform FeatureIndexType = 2
	Uniform    FeatureIndexType = 1
)

type PrimitiveType uint32

const (
	PT_Mesh     PrimitiveType = 0
	PT_Polyline PrimitiveType = 1
	PT_Point    PrimitiveType = 2
)

type MaterialAtlas struct {
	NumMaterials    uint32 `json:"numMaterials,omitempty"`
	HasTranslucency *bool  `json:"hasTranslucency,omitempty"`
	OverridesAlpha  *bool  `json:"overridesAlpha,omitempty"`
}

type Instances struct {
	Count              uint32    `json:"count,omitempty"`
	TransformCenter    []float32 `json:"transformCenter,omitempty"`
	FeatureIds         string    `json:"featureIds,omitempty"`
	Transforms         string    `json:"transforms,omitempty"`
	SymbologyOverrides string    `json:"symbologyOverrides,omitempty"`
}

type VertexTable struct {
	BufferView       string           `json:"bufferView"`
	Count            uint32           `json:"count"`
	NumRgbaPerVertex uint32           `json:"numRgbaPerVertex"`
	NumColors        *uint32          `json:"numColors,omitempty"`
	Width            uint32           `json:"width"`
	Height           uint32           `json:"height"`
	HasTranslucency  bool             `json:"hasTranslucency"`
	FeatureIndexType FeatureIndexType `json:"featureIndexType"`
	FeatureId        *uint32          `json:"featureID,omitempty"`
	UniformColor     uint32           `json:"uniformColor,omitempty"`
	Params           struct {
		DecodeMatrix []float32 `json:"decodeMatrix"`
		DecodedMin   []float32 `json:"decodedMin"`
		DecodedMax   []float32 `json:"decodedMax"`
	} `json:"params"`
	MaterialAtlas *MaterialAtlas `json:"materialAtlas,omitempty"`
	VertexData    []byte         `json:"-"`
}

func (v *VertexTable) GetPosQParams3d() *QParams3d {
	ra := &Range3d{Low: [3]float32{v.Params.DecodedMin[0], v.Params.DecodedMin[1], v.Params.DecodedMin[2]}, High: [3]float32{v.Params.DecodedMax[0], v.Params.DecodedMax[1], v.Params.DecodedMax[2]}}
	qparams := &QParams3d{}
	qparams.SetFromRange(ra, rangeScale16)
	return qparams
}

type Primitive struct {
	Material              string      `json:"material,omitempty"`
	Vertices              VertexTable `json:"vertices,omitempty"`
	IsPlanar              *bool       `json:"isPlanar,omitempty"`
	ViewIndependentOrigin *[3]float32 `json:"viewIndependentOrigin,omitempty"`
	Instances             *Instances  `json:"instances,omitempty"`
}

type SegmentEdges struct {
	Indices                string `json:"indices"`
	EndPointAndQuadIndices string `json:"endPointAndQuadIndices"`
}

type SilhouetteEdges struct {
	SegmentEdges
	NormalPairs string `json:"normalPairs"`
}

type Polyline struct {
	Indices              string `json:"indices"`
	PrevIndices          string `json:"prevIndices"`
	NextIndicesAndParams string `json:"nextIndicesAndParams"`
}

type MeshEdges struct {
	Segments    *SegmentEdges    `json:"segments,omitempty"`
	Silhouettes *SilhouetteEdges `json:"silhouettes,omitempty"`
	Polylines   *Polyline        `json:"polylines,omitempty"`
}

type ClipPlane struct {
	Normal    *[3]float32 `json:"normal,omitempty"`
	Dist      *float32    `json:"dist,omitempty"`
	Invisible *bool       `json:"invisible,omitempty"`
	Interior  *bool       `json:"interior,omitempty"`
}

type ConvexClipPlaneSet []ClipPlane

type UnionOfConvexClipPlaneSets []ConvexClipPlaneSet

type ClipPrimitive interface{}

type ClipPrimitivePlanes struct {
	ClipPrimitive
	Planes *struct {
		Clips     *UnionOfConvexClipPlaneSets `json:"clips,omitempty"`
		Invisible *bool                       `json:"invisible,omitempty"`
	} `json:"planes,omitempty"`
}

type ClipPrimitiveShape struct {
	ClipPrimitive
	Shape *struct {
		Points    [][3]float32 `json:"points,omitempty"`
		Trans     *[16]float32 `json:"trans,omitempty"`
		ZLow      *float32     `json:"zlow,omitempty"`
		ZHigh     *float32     `json:"zhigh,omitempty"`
		Mask      *bool        `json:"mask,omitempty"`
		Invisible *bool        `json:"invisible,omitempty"`
	} `json:"shape,omitempty"`
}

type ClipVector []ClipPrimitive

type Range2d struct {
	Low  [2]float32 `json:"low"`
	High [2]float32 `json:"high"`
}

func (r *Range2d) ExtendXY(x float32, y float32) {
	if x < r.Low[0] {
		r.Low[0] = x
	}
	if x > r.High[0] {
		r.High[0] = x
	}

	if y < r.Low[1] {
		r.Low[1] = y
	}
	if y > r.High[1] {
		r.High[1] = y
	}
}

func (r *Range2d) Extend(xyz [2]float32) {
	r.ExtendXY(xyz[0], xyz[1])
}

func CreateRange2d(points [][2]float32) *Range2d {
	result := &Range2d{}
	for _, point := range points {
		result.Extend(point)
	}
	return result
}

type Range3d struct {
	Low  [3]float32 `json:"low"`
	High [3]float32 `json:"high"`
}

func (r *Range3d) ExtendXYZ(x float32, y float32, z float32) {
	if x < r.Low[0] {
		r.Low[0] = x
	}
	if x > r.High[0] {
		r.High[0] = x
	}

	if y < r.Low[1] {
		r.Low[1] = y
	}
	if y > r.High[1] {
		r.High[1] = y
	}

	if z < r.Low[2] {
		r.Low[2] = z
	}
	if z > r.High[2] {
		r.High[2] = z
	}
}

func (r *Range3d) Extend(xyz [3]float32) {
	r.ExtendXYZ(xyz[0], xyz[1], xyz[2])
}

func CreateRange3d(points [][3]float32) *Range3d {
	result := &Range3d{}
	for _, point := range points {
		result.Extend(point)
	}
	return result
}

type AreaPattern struct {
	Type                  string      `json:"type"` //"areaPattern"
	SymbolName            string      `json:"symbolName"`
	Clip                  ClipVector  `json:"clip"`
	Scale                 float32     `json:"scale"`
	Spacing               [2]float32  `json:"spacing"`
	OrgTransform          [16]float32 `json:"orgTransform"`
	Origin                [2]float32  `json:"origin"`
	XYOffsets             string      `json:"xyOffsets"`
	FeatureId             uint32      `json:"featureId"`
	ModelTransform        [16]float32 `json:"modelTransform"`
	Range                 Range3d     `json:"range"`
	SymbolTranslation     [3]float32  `json:"symbolTranslation"`
	ViewIndependentOrigin *[3]float32 `json:"viewIndependentOrigin"`
}

type SurfaceType uint32

const (
	ST_Unlit            SurfaceType = 0
	ST_Lit              SurfaceType = 1
	ST_Textured         SurfaceType = 2
	ST_TexturedLit      SurfaceType = 3
	ST_VolumeClassifier SurfaceType = 4
)

type Surface struct {
	Type                 SurfaceType `json:"type,omitempty"`
	Indices              string      `json:"indices,omitempty"`
	AlwaysDisplayTexture *bool       `json:"alwaysDisplayTexture,omitempty"`
	UVParams             *struct {
		DecodedMin []float32 `json:"decodedMin"`
		DecodedMax []float32 `json:"decodedMax"`
	} `json:"uvParams"`
}

func (v *Surface) GetUvQParams2d() *QParams2d {
	if v.UVParams == nil {
		return nil
	}
	ra := &Range2d{Low: [2]float32{v.UVParams.DecodedMin[0], v.UVParams.DecodedMin[1]}, High: [2]float32{v.UVParams.DecodedMax[0], v.UVParams.DecodedMax[1]}}
	qparams := &QParams2d{}
	qparams.SetFromRange(ra, rangeScale16)
	return qparams
}

type AuxChannel struct {
	Name    string   `json:"name"`
	Inputs  []uint32 `json:"inputs"`
	Indices []uint32 `json:"indices"`
}

type QuantizedAuxChannel struct {
	AuxChannel
	QOrigin []float32 `json:"qOrigin"`
	QScale  []float32 `json:"qScale"`
}

type AuxChannelTable struct {
	BufferView        string                `json:"bufferView"`
	Width             uint32                `json:"width"`
	Height            uint32                `json:"height"`
	Count             uint32                `json:"count"`
	NumBytesPerVertex uint32                `json:"numBytesPerVertex"`
	Displacements     []QuantizedAuxChannel `json:"displacements,omitempty"`
	Normals           []AuxChannel          `json:"normals,omitempty"`
	Params            []QuantizedAuxChannel `json:"params,omitempty"`
	AuxChannelData    []byte                `json:"-"`
}

type MeshPrimitive struct {
	Primitive
	Type        PrimitiveType    `json:"type"` // Mesh
	Surface     Surface          `json:"surface"`
	Edges       *MeshEdges       `json:"edges,omitempty"`
	AuxChannels *AuxChannelTable `json:"auxChannels,omitempty"`
	AreaPattern *AreaPattern     `json:"areaPattern,omitempty"`
	Data        *MeshData        `json:"-"`
}

type PolylinePrimitive struct {
	Primitive
	Polyline
	Type PrimitiveType `json:"type"` // Polyline
	Data *PolylineData `json:"-"`
}

type PointStringPrimitive struct {
	Primitive
	Type    PrimitiveType    `json:"type"` // Point
	Indices string           `json:"indices,omitempty"`
	Data    *PointStringData `json:"-"`
}

type Mesh struct {
	Primitives interface{} `json:"primitives,omitempty"`
	Layer      string      `json:"layer,omitempty"`
}

func (p *Mesh) UnmarshalJSON(data []byte) error {
	type _mesh struct {
		Primitives []struct {
			Type interface{} `json:"type,omitempty"`
		} `json:"primitives,omitempty"`
		Layer string `json:"layer,omitempty"`
	}
	var m _mesh
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	if len(m.Primitives) > 0 {
		if t, ok := m.Primitives[0].Type.(string); ok && t == "areaPattern" {
			type a_mesh struct {
				Primitives []AreaPattern `json:"primitives,omitempty"`
				Layer      string        `json:"layer,omitempty"`
			}
			var m a_mesh
			err := json.Unmarshal(data, &m)
			if err != nil {
				return err
			}
			p.Primitives = m.Primitives
			p.Layer = m.Layer
		} else if t, ok := m.Primitives[0].Type.(float64); ok {
			switch PrimitiveType(t) {
			case PT_Mesh:
				type m_mesh struct {
					Primitives []MeshPrimitive `json:"primitives,omitempty"`
					Layer      string          `json:"layer,omitempty"`
				}
				var m m_mesh
				err := json.Unmarshal(data, &m)
				if err != nil {
					return err
				}
				p.Primitives = m.Primitives
				p.Layer = m.Layer
			case PT_Polyline:
				type p_mesh struct {
					Primitives []PolylinePrimitive `json:"primitives,omitempty"`
					Layer      string              `json:"layer,omitempty"`
				}
				var m p_mesh
				err := json.Unmarshal(data, &m)
				if err != nil {
					return err
				}
				p.Primitives = m.Primitives
				p.Layer = m.Layer
			case PT_Point:
				type p_mesh struct {
					Primitives []PointStringPrimitive `json:"primitives,omitempty"`
					Layer      string                 `json:"layer,omitempty"`
				}
				var m p_mesh
				err := json.Unmarshal(data, &m)
				if err != nil {
					return err
				}
				p.Primitives = m.Primitives
				p.Layer = m.Layer
			}
		}
	}

	return nil
}

func (p *Mesh) MarshalJSON() ([]byte, error) {
	switch privs := p.Primitives.(type) {
	case []AreaPattern:
		type m_mesh struct {
			Primitives []AreaPattern `json:"primitives,omitempty"`
			Layer      string        `json:"layer,omitempty"`
		}
		m := &m_mesh{Primitives: privs, Layer: p.Layer}
		b, err := json.Marshal(m)
		if err != nil {
			return nil, err
		}
		return b, nil
	case []MeshPrimitive:
		type m_mesh struct {
			Primitives []MeshPrimitive `json:"primitives,omitempty"`
			Layer      string          `json:"layer,omitempty"`
		}
		m := &m_mesh{Primitives: privs, Layer: p.Layer}
		b, err := json.Marshal(m)
		if err != nil {
			return nil, err
		}
		return b, nil
	case []PolylinePrimitive:
		type p_mesh struct {
			Primitives []PolylinePrimitive `json:"primitives,omitempty"`
			Layer      string              `json:"layer,omitempty"`
		}
		m := &p_mesh{Primitives: privs, Layer: p.Layer}
		b, err := json.Marshal(m)
		if err != nil {
			return nil, err
		}
		return b, nil
	case []PointStringPrimitive:
		type p_mesh struct {
			Primitives []PointStringPrimitive `json:"primitives,omitempty"`
			Layer      string                 `json:"layer,omitempty"`
		}
		m := &p_mesh{Primitives: privs, Layer: p.Layer}
		b, err := json.Marshal(m)
		if err != nil {
			return nil, err
		}
		return b, nil
	}
	return nil, errors.New("error")
}

type Buffer struct {
	Extensions gltf.Extensions `json:"extensions,omitempty"`
	Extras     interface{}     `json:"extras,omitempty"`
	Name       string          `json:"name,omitempty"`
	Type       string          `json:"type,omitempty"`
	ByteLength uint32          `json:"byteLength" validate:"required"`
}

type BufferView struct {
	Extensions gltf.Extensions `json:"extensions,omitempty"`
	Extras     interface{}     `json:"extras,omitempty"`
	Buffer     string          `json:"buffer"`
	ByteOffset uint32          `json:"byteOffset,omitempty"`
	ByteLength uint32          `json:"byteLength" validate:"required"`
	ByteStride uint32          `json:"byteStride,omitempty" validate:"omitempty,gte=4,lte=252"`
}

type AreaPatternSymbol struct {
	Primitives interface{} `json:"primitives,omitempty"`
}

func (p *AreaPatternSymbol) UnmarshalJSON(data []byte) error {
	type _mesh struct {
		Primitives []struct {
			Type PrimitiveType `json:"type,omitempty"`
		} `json:"primitives,omitempty"`
	}
	var m _mesh
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}
	if len(m.Primitives) > 0 {
		switch m.Primitives[0].Type {
		case PT_Mesh:
			type m_mesh struct {
				Primitives []MeshPrimitive `json:"primitives,omitempty"`
				Layer      string          `json:"layer,omitempty"`
			}
			var m m_mesh
			err := json.Unmarshal(data, &m)
			if err != nil {
				return err
			}
			p.Primitives = m.Primitives
		case PT_Polyline:
			type p_mesh struct {
				Primitives []PolylinePrimitive `json:"primitives,omitempty"`
				Layer      string              `json:"layer,omitempty"`
			}
			var m p_mesh
			err := json.Unmarshal(data, &m)
			if err != nil {
				return err
			}
			p.Primitives = m.Primitives
		case PT_Point:
			type p_mesh struct {
				Primitives []PointStringPrimitive `json:"primitives,omitempty"`
				Layer      string                 `json:"layer,omitempty"`
			}
			var m p_mesh
			err := json.Unmarshal(data, &m)
			if err != nil {
				return err
			}
			p.Primitives = m.Primitives
		}
	}
	return nil
}

func (p *AreaPatternSymbol) MarshalJSON() ([]byte, error) {
	switch privs := p.Primitives.(type) {
	case []MeshPrimitive:
		type m_mesh struct {
			Primitives []MeshPrimitive `json:"primitives,omitempty"`
		}
		m := &m_mesh{Primitives: privs}
		b, err := json.Marshal(m)
		if err != nil {
			return nil, err
		}
		return b, nil
	case []PolylinePrimitive:
		type p_mesh struct {
			Primitives []PolylinePrimitive `json:"primitives,omitempty"`
		}
		m := &p_mesh{Primitives: privs}
		b, err := json.Marshal(m)
		if err != nil {
			return nil, err
		}
		return b, nil
	case []PointStringPrimitive:
		type p_mesh struct {
			Primitives []PointStringPrimitive `json:"primitives,omitempty"`
		}
		m := &p_mesh{Primitives: privs}
		b, err := json.Marshal(m)
		if err != nil {
			return nil, err
		}
		return b, nil
	}
	return nil, errors.New("error")
}

type RenderTexture struct {
	BufferView    string      `json:"bufferView"`
	Format        uint32      `json:"format"`
	Width         uint32      `json:"width"`
	Height        uint32      `json:"height"`
	IsGlyph       bool        `json:"isGlyph"`
	IsTileSection bool        `json:"isTileSection"`
	TextureData   image.Image `json:"-"`
}

type TextureMappingMode int32

const (
	TM_None             TextureMappingMode = -1
	TM_Parametric       TextureMappingMode = 0
	TM_ElevationDrape   TextureMappingMode = 1
	TM_Planar           TextureMappingMode = 2
	TM_DirectionalDrape TextureMappingMode = 3
	TM_Cubic            TextureMappingMode = 4
	TM_Spherical        TextureMappingMode = 5
	TM_Cylindrical      TextureMappingMode = 6
	TM_Solid            TextureMappingMode = 7
	TM_FrontProject     TextureMappingMode = 8
)

type Texture struct {
	Name   string `json:"name"`
	Params struct {
		Mode          TextureMappingMode `json:"mode"`
		TextureMatrix [][3]float64       `json:"transform"`
		Weight        float64            `json:"weight"`
		WorldMapping  bool               `json:"worldMapping"`
	} `json:"params"`
}

type TextureMapping struct {
	Texture Texture `json:"texture"`
}

type RenderMaterial struct {
	Ambient          float32         `json:"ambient"`
	Diffuse          float32         `json:"diffuse"`
	DiffuseColor     *[3]float32     `json:"diffuseColor,omitempty"`
	EmissiveColor    *[3]float32     `json:"emissiveColor,omitempty"`
	Key              string          `json:"key,omitempty"`
	Reflect          float32         `json:"reflect"`
	ReflectColor     *[3]float32     `json:"reflectColor,omitempty"`
	Refract          float32         `json:"refract"`
	Shadows          bool            `json:"shadows"`
	Specular         float32         `json:"specular"`
	SpecularColor    *[3]float32     `json:"specularColor,omitempty"`
	SpecularExponent float32         `json:"specularExponent"`
	Transparency     *float32        `json:"transparency,omitempty"`
	TextureMapping   *TextureMapping `json:"textureMapping"`
}

type Scene struct {
	Name           string   `json:"name,omitempty"`
	Nodes          []string `json:"nodes,omitempty"`
	AnimationNodes []string `json:"animationNodes,omitempty"`
}

type Material struct {
	CategoryId     string   `json:"categoryId"`
	FillColor      *uint32  `json:"fillColor,omitempty"`
	FillFlags      *uint32  `json:"fillFlags,omitempty"`
	IgnoreLighting *bool    `json:"ignoreLighting,omitempty"`
	LineColor      *uint32  `json:"lineColor,omitempty"`
	LinePixels     *uint32  `json:"linePixels,omitempty"`
	LineWidth      *uint32  `json:"lineWidth,omitempty"`
	MaterialId     string   `json:"materialId"`
	SubCategoryId  string   `json:"subCategoryId"`
	Texture        *Texture `json:"texture,omitempty"`
	Type           uint32   `json:"type"`
}

type AnimationNodes struct {
	BufferView string `json:"bufferView"`
	BytesPerId uint32 `json:"bytesPerId"`
}

type Document struct {
	ExtensionsUsed   []string                   `json:"extensionsUsed,omitempty"`
	GLExtensionsUsed []string                   `json:"glExtensionsUsed,omitempty"`
	Buffers          map[string]*Buffer         `json:"buffers,omitempty" validate:"dive"`
	BufferViews      map[string]*BufferView     `json:"bufferViews,omitempty" validate:"dive"`
	Materials        map[string]*Material       `json:"materials,omitempty" validate:"dive"`
	Meshes           map[string]*Mesh           `json:"meshes,omitempty" validate:"dive"`
	Nodes            map[string]string          `json:"nodes,omitempty" validate:"dive"`
	Scene            *string                    `json:"scene,omitempty"`
	Scenes           map[string]*Scene          `json:"scenes,omitempty" validate:"dive"`
	NamedTextures    map[string]*RenderTexture  `json:"namedTextures,omitempty" validate:"dive"`
	RenderMaterials  map[string]*RenderMaterial `json:"renderMaterials,omitempty" validate:"dive"`
	AnimationNodes   *AnimationNodes            `json:"animationNodes,omitempty"`
	chunks           []chunkData                `json:"-"`
}

func newString(s string) *string {
	return &s
}

func NewDocument() *Document {
	return &Document{
		Scene:  newString("defaultScene"),
		Scenes: map[string]*Scene{"defaultScene": {Nodes: []string{"rootNode"}}},
	}
}

func (doc *Document) FindBuffer(bufferView string) []byte {
	for i := range doc.chunks {
		if doc.chunks[i].name == bufferView {
			return doc.chunks[i].data
		}
	}
	return nil
}

func (doc *Document) decodeChunkData(data []byte) {
	chunkMap := make(map[string]*chunkData)
	for k, v := range doc.BufferViews {
		byteOffset := int(v.ByteOffset)
		byteLength := int(v.ByteLength)
		cd := chunkData{name: k, data: data[byteOffset : byteOffset+byteLength]}
		doc.chunks = append(doc.chunks, cd)
		chunkMap[k] = &cd
	}

	for _, m := range doc.Meshes {
		switch privs := m.Primitives.(type) {
		case []MeshPrimitive:
			for i := range privs {
				privs[i].Data = &MeshData{Type: privs[i].Surface.Type}

				uvq := privs[i].Surface.GetUvQParams2d()
				posq := privs[i].Vertices.GetPosQParams3d()

				if cd, ok := chunkMap[privs[i].Surface.Indices]; ok {
					privs[i].Data.DecodeIndices(cd.data)
				}

				if cd, ok := chunkMap[privs[i].Vertices.BufferView]; ok {
					privs[i].Data.DecodeVertexs(cd.data, privs[i].Vertices.Count)
				}

				privs[i].Data.UnQuantize(posq, uvq)
			}
		case []PolylinePrimitive:
			for i := range privs {
				privs[i].Data = &PolylineData{}

				posq := privs[i].Vertices.GetPosQParams3d()

				if cd, ok := chunkMap[privs[i].Indices]; ok {
					privs[i].Data.DecodeIndices(cd.data)
				}

				if cd, ok := chunkMap[privs[i].Vertices.BufferView]; ok {
					privs[i].Data.DecodeVertexs(cd.data)
				}

				privs[i].Data.UnQuantize(posq)
			}
		case []PointStringPrimitive:
			for i := range privs {
				privs[i].Data = &PointStringData{}

				posq := privs[i].Vertices.GetPosQParams3d()

				if cd, ok := chunkMap[privs[i].Indices]; ok {
					privs[i].Data.DecodeIndices(cd.data)
				}

				if cd, ok := chunkMap[privs[i].Vertices.BufferView]; ok {
					privs[i].Data.DecodeVertexs(cd.data)
				}

				privs[i].Data.UnQuantize(posq)
			}
		}
	}

	for _, t := range doc.NamedTextures {
		if cd, ok := chunkMap[t.BufferView]; ok {
			t.TextureData = DecodeTexture(cd.data, TextureFormat(t.Format))
		}
	}
}

func (doc *Document) encodeChunkData() ([][]byte, uint32) {
	doc.Buffers = make(map[string]*Buffer)
	doc.BufferViews = make(map[string]*BufferView)

	const bufferName = "binary_glTF"

	chunkid := 0

	for _, m := range doc.Meshes {
		switch privs := m.Primitives.(type) {
		case []MeshPrimitive:
			for i := range privs {
				if privs[i].Data != nil {
					posr, uvr := privs[i].Data.Quantize()

					if privs[i].Surface.Indices == "" {
						privs[i].Surface.Indices = fmt.Sprintf("buffer-%d", chunkid)
						chunkid++
					}
					doc.chunks = append(doc.chunks, chunkData{name: privs[i].Surface.Indices, data: privs[i].Data.EncodeIndices()})

					if uvr != nil {
						privs[i].Surface.UVParams.DecodedMin = uvr.Low[:]
						privs[i].Surface.UVParams.DecodedMax = uvr.High[:]
					}

					if privs[i].Vertices.BufferView == "" {
						privs[i].Vertices.BufferView = fmt.Sprintf("buffer-%d", chunkid)
						chunkid++
					}
					doc.chunks = append(doc.chunks, chunkData{name: privs[i].Vertices.BufferView, data: privs[i].Data.EncodeVertexs()})

					if posr != nil {
						privs[i].Vertices.Params.DecodedMin = posr.Low[:]
						privs[i].Vertices.Params.DecodedMax = posr.High[:]
					}
				}
			}
		case []PolylinePrimitive:
			for i := range privs {
				if privs[i].Data != nil {
					posr := privs[i].Data.Quantize()

					if privs[i].Indices == "" {
						privs[i].Indices = fmt.Sprintf("buffer-%d", chunkid)
						chunkid++
					}
					doc.chunks = append(doc.chunks, chunkData{name: privs[i].Indices, data: privs[i].Data.EncodeIndices()})
					if privs[i].Vertices.BufferView == "" {
						privs[i].Vertices.BufferView = fmt.Sprintf("buffer-%d", chunkid)
						chunkid++
					}
					doc.chunks = append(doc.chunks, chunkData{name: privs[i].Vertices.BufferView, data: privs[i].Data.EncodeVertexs()})

					if posr != nil {
						privs[i].Vertices.Params.DecodedMin = posr.Low[:]
						privs[i].Vertices.Params.DecodedMax = posr.High[:]
					}
				}
			}
		case []PointStringPrimitive:
			for i := range privs {
				if privs[i].Data != nil {
					posr := privs[i].Data.Quantize()

					if privs[i].Indices == "" {
						privs[i].Indices = fmt.Sprintf("buffer-%d", chunkid)
						chunkid++
					}
					doc.chunks = append(doc.chunks, chunkData{name: privs[i].Indices, data: privs[i].Data.EncodeIndices()})
					if privs[i].Vertices.BufferView == "" {
						privs[i].Vertices.BufferView = fmt.Sprintf("buffer-%d", chunkid)
						chunkid++
					}
					doc.chunks = append(doc.chunks, chunkData{name: privs[i].Vertices.BufferView, data: privs[i].Data.EncodeVertexs()})

					if posr != nil {
						privs[i].Vertices.Params.DecodedMin = posr.Low[:]
						privs[i].Vertices.Params.DecodedMax = posr.High[:]
					}
				}
			}
		}
	}

	for _, t := range doc.NamedTextures {
		if t.BufferView == "" {
			t.BufferView = fmt.Sprintf("buffer-%d", chunkid)
			chunkid++
		}
		doc.chunks = append(doc.chunks, chunkData{name: t.BufferView, data: EncodeTexture(t.TextureData, TextureFormat(t.Format))})
	}

	offset := uint32(0)

	out := make([][]byte, len(doc.chunks))

	for i, ck := range doc.chunks {
		dataLen := uint32(len(ck.data))
		doc.BufferViews[ck.name] = &BufferView{Buffer: bufferName, ByteOffset: offset, ByteLength: dataLen}
		out[i] = createPaddingBytes(ck.data, dataLen, 8, 0x20)
		offset += uint32(len(out[i]))
	}

	doc.Buffers[bufferName] = &Buffer{ByteLength: offset}

	return out, offset
}
