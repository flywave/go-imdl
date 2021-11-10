package imdl

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"math"
	"os"
	"unsafe"
)

const (
	defaultMaxExternalBufferCount = 10
	defaultMaxMemoryAllocation    = math.MaxUint32 // 4GB
)

func Open(name string) (*Document, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	dec := NewDecoder(f)
	doc := new(Document)
	if err = dec.Decode(doc); err != nil {
		doc = nil
	}
	return doc, err
}

type Decoder struct {
	MaxExternalBufferCount int
	MaxMemoryAllocation    uint64
	r                      *bufio.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		MaxExternalBufferCount: defaultMaxExternalBufferCount,
		MaxMemoryAllocation:    defaultMaxMemoryAllocation,
		r:                      bufio.NewReader(r),
	}
}

func (d *Decoder) Decode(doc *Document) error {
	_, err := d.decodeDocument(doc)
	if err != nil {
		return err
	}

	return nil
}

func (d *Decoder) validateDocumentQuotas(doc *Document, isBinary bool) error {
	var allocs uint64
	for _, b := range doc.Buffers {
		allocs += uint64(b.ByteLength)
	}
	if allocs > d.MaxMemoryAllocation {
		return errors.New("imdl: Memory allocation count quota exceeded")
	}
	return nil
}

func (d *Decoder) decodeDocument(doc *Document) (bool, error) {
	glbHeader, err := d.readGLBHeader()
	if err != nil {
		return false, err
	}
	var (
		jd       *json.Decoder
		isBinary bool
	)
	if glbHeader != nil {
		jd = json.NewDecoder(&io.LimitedReader{R: d.r, N: int64(glbHeader.JSONHeader.Length)})
		isBinary = true
	} else {
		jd = json.NewDecoder(d.r)
		isBinary = false
	}

	err = jd.Decode(doc)
	if err == nil {
		err = d.validateDocumentQuotas(doc, isBinary)
	}

	if data, err := d.decodeBinaryBuffer(glbHeader); err != nil {
		return isBinary, err
	} else {
		doc.decodeChunkData(data)
	}

	return isBinary, err
}

func (d *Decoder) readGLBHeader() (*glbHeader, error) {
	var header glbHeader
	chunk, err := d.r.Peek(int(unsafe.Sizeof(header)))
	if err != nil {
		return nil, nil
	}
	r := bytes.NewReader(chunk)
	binary.Read(r, binary.LittleEndian, &header)
	if header.Magic != glbHeaderMagic {
		return nil, nil
	}
	d.r.Read(chunk)
	return &header, d.validateGLBHeader(&header)
}

func (d *Decoder) validateGLBHeader(header *glbHeader) error {
	if header.Version != 1 || (header.JSONHeader.Length+uint32(unsafe.Sizeof(header))) > header.Length {
		return errors.New("imdl: Invalid imdl GLB JSON header")
	}
	return nil
}

func (d *Decoder) decodeBinaryBuffer(h *glbHeader) ([]byte, error) {
	byteLength := int(h.Length) - binary.Size(glbHeader{}) - int(h.JSONHeader.Length)
	data := make([]byte, byteLength)
	_, err := io.ReadFull(d.r, data)
	return data, err
}
