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
	"path/filepath"

	cuedb "github.com/k1LoW/tbls-driver-tailordb/tailordb/cue"
	"github.com/k1LoW/tbls-driver-tailordb/tailordb/tf"
	"github.com/k1LoW/tbls-driver-tailordb/version"
	"github.com/k1LoW/tbls/schema"
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
		fi, err := os.Stat(root)
		if err != nil {
			return err
		}
		var s *schema.SchemaJSON
		switch {
		case filepath.Ext(root) == ".cue":
			s, err = cuedb.Analyze(root)
			if err != nil {
				return err
			}
		case fi.IsDir():
			s, err = tf.Analyze(root)
			if err != nil {
				return err
			}
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
