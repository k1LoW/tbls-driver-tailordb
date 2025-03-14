package cue

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/tenntenn/golden"
)

func TestAnalyze(t *testing.T) {
	root := "../../testdata/manifests/typecue/tailordb.cue"
	s, err := Analyze(root)
	if err != nil {
		t.Fatal(err)
	}
	f := "typecue.json"
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(s); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if os.Getenv("UPDATE_GOLDEN") != "" {
		golden.Update(t, "../../testdata", f, got)
		return
	}
	if diff := golden.Diff(t, "../../testdata", f, got); diff != "" {
		t.Error(diff)
	}
}
