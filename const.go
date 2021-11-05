package imdl

const (
	glbHeaderMagic = 0x46546c67
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
