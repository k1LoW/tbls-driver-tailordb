package cue

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/tenntenn/golden"
)

func TestAnalyze(t *testing.T) {
	root := "../../testdata/manifests/typecue/tailordb.cue"
	modRoot := filepath.Dir(root)
	ctx := cuecontext.New()
	cfg := &load.Config{
		ModuleRoot: modRoot,
	}
	insts := load.Instances([]string{root}, cfg)
	v := ctx.BuildInstance(insts[0])
	s, err := Analyze(v)
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
