package tf

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/tenntenn/golden"
)

func TestAnalyze(t *testing.T) {
	root := "../../testdata/manifests/typetf"
	s, err := Analyze(root)
	if err != nil {
		t.Fatal(err)
	}
	f := "typetf.json"
	b, err := json.MarshalIndent(s, "", "	 ")
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	if os.Getenv("UPDATE_GOLDEN") != "" {
		golden.Update(t, "../../testdata", f, got)
		return
	}
	if diff := golden.Diff(t, "../../testdata", f, got); diff != "" {
		t.Error(diff)
	}
}
