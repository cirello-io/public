// Copyright 2019 cirello.io and github.com/ucirello
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Command openapigen is an OpenAPI v2 renderer. Internally it uses Go's
// template engine to render the output.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	tplHTML "html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	tplText "text/template"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/iancoleman/strcase"
)

var (
	spec        = flag.String("spec", ".", "openAPI json filename")
	isHTML      = flag.Bool("html", false, "use html/template")
	template    = flag.String("template", "", "location of the template file")
	output      = flag.String("output", "", "filename of the expected output")
	isOpenAPIV2 = flag.Bool("v2mode", false, "indicates the spec is an openAPI v2 file")
	view        = flag.Bool("view", false, "print parsed spec file")
)

func main() {
	flag.Parse()
	log.SetFlags(0)
	log.SetPrefix("openapigen: ")
	fd, err := os.Open(*spec)
	if err != nil {
		log.Fatal("cannot open swagger json file:", err)
	}
	log.Println("Decoding spec file with https://godoc.org/github.com/getkin/kin-openapi/openapi2#Swagger")
	var swagger *openapi3.Swagger
	if *isOpenAPIV2 {
		var swaggerV2 openapi2.Swagger
		err := json.NewDecoder(fd).Decode(&swaggerV2)
		if err != nil {
			log.Fatal("cannot parse swaggerV2 json file:", err)
		}
		swagger, err = openapi2conv.ToV3Swagger(&swaggerV2)
		if err != nil {
			log.Fatal("cannot convert from v2 to v3:", err)
		}
	}
	if *view {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "	")
		err := enc.Encode(swagger)
		if err != nil {
			log.Fatal("cannot encode spec file")
		}
		os.Exit(0)
	}
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("cannot detect current working directory:", err)
	}
	templateDir, err := filepath.Abs(*template)
	if err != nil {
		log.Fatal("cannot calculate absolute directory for template:", err)
	}
	outputDir, err := filepath.Abs(*output)
	if err != nil {
		log.Fatal("cannot calculate absolute directory for output:", err)
	}
	err = filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ".tpl" || (filepath.Ext(path) == ".tpl" && info.IsDir()) {
			return nil
		}
		relpath, err := filepath.Rel(wd, path)
		if err != nil {
			return fmt.Errorf("cannot calculate relative directory for %s: %w", path, err)
		}
		log.Println("rendering", relpath)
		tplRaw, err := readFile(path)
		if err != nil {
			return fmt.Errorf("cannot load template: %w", err)
		}
		var tpl interface {
			Execute(wr io.Writer, data interface{}) error
		}
		funcs := map[string]interface{}{
			"firstLetter": func(s string) string {
				if len(s) == 0 {
					return ""
				}
				return string(s[0])
			},
			"toLower":    strings.ToLower,
			"camel":      strcase.ToCamel,
			"lowerCamel": strcase.ToLowerCamel,
			"snake":      strcase.ToSnake,
			"stripDefinitionPrefix": func(s string) string {
				return strings.TrimPrefix(s, "#/definitions/")
			},
			"debug": func(v interface{}) (string, error) {
				var buf bytes.Buffer
				enc := json.NewEncoder(&buf)
				enc.SetIndent("", "	")
				err := enc.Encode(v)
				if err != nil {
					return "", fmt.Errorf("cannot marshal: %w", err)
				}
				return buf.String(), nil
			},
			"uniquePathTags": func() []string {
				var tags []string
				for _, pathItem := range swagger.Paths {
					if pathItem.Connect != nil {
						tags = append(tags, pathItem.Connect.Tags...)
					}
					if pathItem.Delete != nil {
						tags = append(tags, pathItem.Delete.Tags...)
					}
					if pathItem.Get != nil {
						tags = append(tags, pathItem.Get.Tags...)
					}
					if pathItem.Head != nil {
						tags = append(tags, pathItem.Head.Tags...)
					}
					if pathItem.Options != nil {
						tags = append(tags, pathItem.Options.Tags...)
					}
					if pathItem.Patch != nil {
						tags = append(tags, pathItem.Patch.Tags...)
					}
					if pathItem.Post != nil {
						tags = append(tags, pathItem.Post.Tags...)
					}
					if pathItem.Put != nil {
						tags = append(tags, pathItem.Put.Tags...)
					}
					if pathItem.Trace != nil {
						tags = append(tags, pathItem.Trace.Tags...)
					}
				}
				tagsDict := make(map[string]struct{})
				for _, tag := range tags {
					tagsDict[tag] = struct{}{}
				}
				uniqTags := []string{}
				for tag := range tagsDict {
					uniqTags = append(uniqTags, tag)
				}
				sort.Strings(uniqTags)
				return uniqTags
			},
		}
		switch {
		case *isHTML:
			tpl, err = tplHTML.New("openapigen").Funcs(tplHTML.FuncMap(funcs)).Option("missingkey=zero").Parse(tplRaw)
			if err != nil {
				return fmt.Errorf("cannot parse template (html mode): %w", err)
			}
		default:
			tpl, err = tplText.New("openapigen").Funcs(tplText.FuncMap(funcs)).Option("missingkey=zero").Parse(tplRaw)
			if err != nil {
				return fmt.Errorf("cannot parse template (text mode): %w", err)
			}
		}
		dir := filepath.Dir(filepath.Join(outputDir, strings.TrimPrefix(path, templateDir)))
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, os.ModePerm&0755); err != nil {
				return fmt.Errorf("cannot create directory %s: %w", dir, err)
			}
		}
		fd, err := os.Create(strings.TrimSuffix(filepath.Join(dir, filepath.Base(path)), ".tpl"))
		if err != nil {
			return fmt.Errorf("cannot create output file: %w", err)
		}
		defer fd.Close()
		if err := tpl.Execute(fd, swagger); err != nil {
			return fmt.Errorf("cannot render output: %w", err)
		}
		return nil
	})
	if err != nil {
		log.Fatal("cannot iterate through template files:", err)
	}
}

func readFile(fn string) (string, error) {
	b, err := ioutil.ReadFile(fn)
	return string(b), err
}
