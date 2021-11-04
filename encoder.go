package imdl

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/flywave/gltf"
)

func Save(doc *Document, name string) error {
	return saveImdl(doc, name, false)
}

func SaveBinary(doc *Document, name string) error {
	return saveImdl(doc, name, true)
}

func saveImdl(doc *Document, name string, asBinary bool) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	e := NewEncoder(f).WithWriteHandler(&gltf.RelativeFileHandler{Dir: filepath.Dir(name)})
	e.AsBinary = asBinary
	if err := e.Encode(doc); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

type Encoder struct {
	AsBinary     bool
	WriteHandler gltf.WriteHandler
	w            io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		AsBinary:     true,
		WriteHandler: new(gltf.RelativeFileHandler),
		w:            w,
	}
}

func (e *Encoder) WithWriteHandler(h gltf.WriteHandler) *Encoder {
	e.WriteHandler = h
	return e
}

func (e *Encoder) Encode(doc *Document) error {
	var err error
	if e.AsBinary {
		_, err = e.encodeBinary(doc)
	} else {
		err = json.NewEncoder(e.w).Encode(doc)
	}
	if err != nil {
		return err
	}

	for k := range doc.Buffers {
		buffer := doc.Buffers[k]
		if len(buffer.Data) == 0 || buffer.IsEmbeddedResource() {
			continue
		}
		if err = e.encodeBuffer(buffer); err != nil {
			return err
		}
	}

	return err
}

func (e *Encoder) encodeBuffer(buffer *gltf.Buffer) error {
	if err := validateBufferURI(buffer.URI); err != nil {
		return err
	}

	return e.WriteHandler.WriteResource(buffer.URI, buffer.Data)
}

func validateBufferURI(uri string) error {
	if uri == "" || strings.Contains(uri, "..") || strings.HasPrefix(uri, "/") || strings.HasPrefix(uri, "\\") {
		return fmt.Errorf("gltf: Invalid buffer.uri value '%s'", uri)
	}
	return nil
}

func (e *Encoder) encodeBinary(doc *Document) (bool, error) {
	jsonText, err := json.Marshal(doc)
	if err != nil {
		return false, err
	}
	jsonHeader := chunkHeader{
		Length: uint32(((len(jsonText) + 3) / 4) * 4),
		Type:   glbChunkJSON,
	}
	header := glbHeader{
		Magic:      glbHeaderMagic,
		Version:    2,
		Length:     12 + 8 + jsonHeader.Length,
		JSONHeader: jsonHeader,
	}
	headerPadding := make([]byte, header.JSONHeader.Length-uint32(len(jsonText)))
	for i := range headerPadding {
		headerPadding[i] = ' '
	}

	var firstBuffers *gltf.Buffer

	for k := range doc.Buffers {
		firstBuffers = doc.Buffers[k]
		break
	}

	hasBinChunk := len(doc.Buffers) > 0 && firstBuffers.URI == ""
	var binPaddedLength uint32
	if hasBinChunk {
		binPaddedLength = ((firstBuffers.ByteLength + 3) / 4) * 4
		header.Length += uint32(8) + binPaddedLength
	}

	err = binary.Write(e.w, binary.LittleEndian, &header)
	if err != nil {
		return false, err
	}
	e.w.Write(jsonText)
	e.w.Write(headerPadding)

	if hasBinChunk {
		binBuffer := firstBuffers
		binPadding := make([]byte, binPaddedLength-binBuffer.ByteLength)
		for i := range binPadding {
			binPadding[i] = 0
		}
		binHeader := chunkHeader{Length: binPaddedLength, Type: glbChunkBIN}
		binary.Write(e.w, binary.LittleEndian, &binHeader)
		e.w.Write(binBuffer.Data)
		_, err = e.w.Write(binPadding)
	}

	return hasBinChunk, err
}

func (p *Mesh) UnmarshalJSON(data []byte) error {
	type alias Mesh
	tmp := alias(Mesh{})
	err := json.Unmarshal(data, &tmp)
	if err == nil {
		*p = Mesh(tmp)
	}
	return err
}

func (p *Mesh) MarshalJSON() ([]byte, error) {
	type alias Mesh
	tmp := &struct {
		alias
	}{
		alias: (alias)(*p),
	}

	return json.Marshal(tmp)
}
