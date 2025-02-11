/*
Copyright © 2025 Ken'ichiro Oyama <k1lowxb@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	cuedb "github.com/k1LoW/tbls-driver-tailordb/tailordb/cue"
	"github.com/k1LoW/tbls-driver-tailordb/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "tbls-driver-tailordb",
	Short:   "tbls driver for TailorDB schema definition",
	Long:    `tbls driver for TailorDB schema definition.`,
	Args:    cobra.NoArgs,
	Version: version.Version,
	RunE: func(cmd *cobra.Command, args []string) error {
		dsn := os.Getenv("TBLS_DSN")
		if dsn == "" {
			return fmt.Errorf("env TBLS_DSN is required")
		}
		u, err := url.Parse(dsn)
		if err != nil {
			return err
		}
		root := path.Join(u.Host, u.Path)
		modRoot, err := lookupModRoot(root)
		if err != nil {
			return err
		}

		ctx := cuecontext.New()
		cfg := &load.Config{
			ModuleRoot: modRoot,
		}
		insts := load.Instances([]string{root}, cfg)
		v := ctx.BuildInstance(insts[0])
		services := v.LookupPath(cue.MakePath(cue.Str("Services")))
		var database cue.Value
		if !services.Exists() {
			kind, err := v.LookupPath(cue.MakePath(cue.Str("Kind"))).String()
			if err != nil {
				return err
			}
			if kind != "tailordb" {
				return fmt.Errorf("no tailordb found")
			}
			database = v
		} else {
			databases, err := findTailordb(services)
			if err != nil {
				return err
			}
			switch len(databases) {
			case 0:
				return fmt.Errorf("no tailordb found")
			case 1:
			default:
				return fmt.Errorf("multiple tailordb found")
			}
			database = databases[0]
		}
		s, err := cuedb.Analyze(database)
		if err != nil {
			return err
		}
		b, err := json.MarshalIndent(s, "", "  ")
		if err != nil {
			return err
		}
		fmt.Print(string(b))
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {}

func findTailordb(v cue.Value) ([]cue.Value, error) {
	const kind = "tailordb"
	var result []cue.Value
	iter, err := v.List()
	if err != nil {
		return nil, err
	}
	for iter.Next() {
		value := iter.Value()
		kindValue := value.LookupPath(cue.MakePath(cue.Str("Kind")))
		if kindValue.Exists() {
			k, err := kindValue.String()
			if err != nil {
				return nil, err
			}
			if k == kind {
				result = append(result, value)
			}
		}
	}

	return result, nil
}

func lookupModRoot(root string) (string, error) {
	if fi, err := os.Stat(path.Join(root, "cue.mod")); err == nil && fi.IsDir() {
		return root, nil
	}
	if root == "." || root == "/" {
		return "", fmt.Errorf("module root not found")
	}
	return lookupModRoot(path.Dir(root))
}
