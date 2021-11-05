package imdl

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

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
		err = e.encodeBinary(doc)
	} else {
		err = json.NewEncoder(e.w).Encode(doc)
	}
	if err != nil {
		return err
	}

	return err
}

func (e *Encoder) encodeBinary(doc *Document) error {
	chunks := doc.encodeChunkData()

	jsonText, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	jsonHeader := JSONHeader{
		Length: uint32(((len(jsonText) + 3) / 4) * 4),
		Type:   glbChunkJSON,
	}
	header := glbHeader{
		Magic:      glbHeaderMagic,
		Version:    1,
		Length:     12 + 8 + jsonHeader.Length,
		JSONHeader: jsonHeader,
	}
	headerPadding := make([]byte, header.JSONHeader.Length-uint32(len(jsonText)))
	for i := range headerPadding {
		headerPadding[i] = ' '
	}

	err = binary.Write(e.w, binary.LittleEndian, &header)
	if err != nil {
		return err
	}
	e.w.Write(jsonText)
	e.w.Write(headerPadding)

	for i := range chunks {
		e.w.Write(chunks[i])
	}

	return nil
}
