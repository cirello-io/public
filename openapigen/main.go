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
	"encoding/json"
	"flag"
	tplHTML "html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	tplText "text/template"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/iancoleman/strcase"
)

var (
	spec     = flag.String("spec", ".", "swagger (openAPI v2) json filename")
	isHTML   = flag.Bool("html", false, "use html/template")
	template = flag.String("template", "", "location of the template file")
	output   = flag.String("output", "", "filename of the expected output")
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
	var swagger openapi2.Swagger
	if err := json.NewDecoder(fd).Decode(&swagger); err != nil {
		log.Fatal("cannot parse swagger json file:", err)
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
			log.Fatalf("cannot calculate relative directory for %s:", path, err)
		}
		log.Println("rendering", relpath)
		tplRaw, err := readFile(path)
		if err != nil {
			log.Fatal("cannot load template:", err)
		}
		var tpl interface {
			Execute(wr io.Writer, data interface{}) error
		}
		funcs := map[string]interface{}{
			"camel":      strcase.ToCamel,
			"lowerCamel": strcase.ToLowerCamel,
			"snake":      strcase.ToSnake,
			"stripDefinitionPrefix": func(s string) string {
				return strings.TrimPrefix(s, "#/definitions/")
			},
		}
		switch {
		case *isHTML:
			tpl, err = tplHTML.New("openapigen").Funcs(tplHTML.FuncMap(funcs)).Option("missingkey=zero").Parse(tplRaw)
			if err != nil {
				log.Fatal("cannot parse template (html mode):", err)
			}
		default:
			tpl, err = tplText.New("openapigen").Funcs(tplText.FuncMap(funcs)).Option("missingkey=zero").Parse(tplRaw)
			if err != nil {
				log.Fatal("cannot parse template (text mode):", err)
			}
		}
		dir := filepath.Dir(filepath.Join(outputDir, strings.TrimPrefix(path, templateDir)))
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, os.ModePerm&0755); err != nil {
				log.Fatalf("cannot create directory %s: %v", dir, err)
			}
		}
		fd, err := os.Create(strings.TrimSuffix(filepath.Join(dir, filepath.Base(path)), ".tpl"))
		if err != nil {
			log.Fatal("cannot create output file:", err)
		}
		defer fd.Close()
		if err := tpl.Execute(fd, swagger); err != nil {
			log.Fatal("cannot render output:", err)
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
