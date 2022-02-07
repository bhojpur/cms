package cmd

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/bhojpur/application/pkg/utils"
	"github.com/spf13/cobra"
)

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "To compile your application template files into Go source code",
	Run: func(cmd *cobra.Command, args []string) {
		exitAfterCompile := flag.Bool("exit-after-compile", false, "Exit after compiling application templates")
		flag.Parse()

		if len(args) == 0 {
			fmt.Println("invalid argument")
			os.Exit(1)
		}

		fmt.Println("Initializing Bhojpur CMS application template compiler...")

		destPath := args[0]
		funcMap := map[string]interface{}{
			"package_path": func() string {
				return destPath
			},
			"package_name": func() string {
				return path.Base(destPath)
			},
			"exit_after_compile": func() bool {
				return *exitAfterCompile
			},
		}

		hasExists := false
		for _, gopath := range utils.GOPATH() {
			sourcePath := filepath.Join(gopath, "templates")
			_, err := os.Stat(sourcePath)
			if err == nil {
				hasExists = true
			}
			err = filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
				if err == nil {
					var relativePath = strings.TrimPrefix(path, sourcePath)

					if info.IsDir() {
						err = os.MkdirAll(filepath.Join(destPath, relativePath), os.ModePerm)
					} else if info.Mode().IsRegular() {
						if source, err := ioutil.ReadFile(path); err == nil {
							var tmpl *template.Template
							if tmpl, err = template.New("").Funcs(funcMap).Parse(string(source)); err == nil {
								var result = bytes.NewBufferString("")
								if err = tmpl.Execute(result, ""); err != nil {
									return err
								}
								source = result.Bytes()
							} else {
								return err
							}
							if err = ioutil.WriteFile(filepath.Join(destPath, strings.TrimSuffix(relativePath, ".template")), source, os.ModePerm); err != nil {
								fmt.Println(err)
							}
						}
					}
				}
				return err
			})

			if hasExists && err == nil {
				fmt.Printf("copy application template from %s to %s\n", sourcePath, destPath)
				break
			}

			if err != nil {
				fmt.Println("failed to copy application template files:", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(templateCmd)
}
