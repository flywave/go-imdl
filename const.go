package imdl

const (
	glbHeaderMagic = 0x46546c67
	glbChunkJSON   = 0x4e4f534a
	glbChunkBIN    = 0x004e4942
)

type JSONHeader struct {
	Length uint32
	Type   uint32
}

type glbHeader struct {
	Magic      uint32
	Version    uint32
	Length     uint32
	JSONHeader JSONHeader
}
