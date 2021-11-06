package imdl

import (
	"os"
	"testing"
)

func TestEncode(t *testing.T) {
	doc, err := Open("./testdata/-3-1-0-0-0-1.gltf")
	if err != nil || doc == nil {
		t.FailNow()
	}

	saveImdl(doc, "./test.gltf", true)

	doc2, err := Open("./test.gltf")
	if err != nil || doc2 == nil {
		t.FailNow()
	}

	os.Remove("./test.gltf")

}
