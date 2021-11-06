package imdl

import (
	"testing"
)

func TestEncode(t *testing.T) {
	doc, err := Open("./testdata/-3-1-0-0-0-1.gltf")
	if err != nil || doc == nil {
		t.FailNow()
	}

}
