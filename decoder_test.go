package imdl

import (
	"encoding/json"
	"os"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	doc := &Document{}
	f, err := os.Open("./testdata/-3-1-0-0-0-1.json")
	if err != nil {
		t.FailNow()
	}
	defer f.Close()
	jd := json.NewDecoder(f)
	err = jd.Decode(doc)
	if err != nil {
		t.FailNow()
	}
}

func TestDecode(t *testing.T) {
	doc, err := Open("./testdata/-3-1-0-0-0-1.gltf")
	if err != nil || doc == nil {
		t.FailNow()
	}
}
