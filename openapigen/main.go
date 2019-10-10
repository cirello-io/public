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
	tplText "text/template"

	"github.com/getkin/kin-openapi/openapi2"
)

var (
	input    = flag.String("file", "", "swagger (openAPI v2) json filename")
	isHTML   = flag.Bool("html", false, "use html/template")
	template = flag.String("template", "", "location of the template file")
	output   = flag.String("output", "", "filename of the expected output")
)

func main() {
	flag.Parse()
	log.SetFlags(0)
	log.SetPrefix("openapigen: ")
	var err error
	tplRaw := defaultRawTemplate
	if *template != "" {
		tplRaw, err = readFile(*template)
		if err != nil {
			log.Fatal("cannot load template:", err)
		}
	} else {
		log.Println("using default template")
	}
	var tpl interface {
		Execute(wr io.Writer, data interface{}) error
	}
	switch {
	case *isHTML:
		tpl, err = tplHTML.New("openapigen").Option("missingkey=zero").Parse(tplRaw)
		if err != nil {
			log.Fatal("cannot parse template (html mode):", err)
		}
	default:
		tpl, err = tplText.New("openapigen").Option("missingkey=zero").Parse(tplRaw)
		if err != nil {
			log.Fatal("cannot parse template (text mode):", err)
		}
	}
	fd, err := os.Open(*input)
	if err != nil {
		log.Fatal("cannot open swagger json file:", err)
	}
	log.Println("Decoding input file with https://godoc.org/github.com/getkin/kin-openapi/openapi2#Swagger")
	var swagger openapi2.Swagger
	if err := json.NewDecoder(fd).Decode(&swagger); err != nil {
		log.Fatal("cannot parse swagger json file:", err)
	}
	var out io.Writer = os.Stdout
	if *output != "" {
		fd, err := os.Create(*output)
		if err != nil {
			log.Fatal("cannot create output file:", err)
		}
		defer fd.Close()
		out = fd
	}
	if err := tpl.Execute(out, swagger); err != nil {
		log.Fatal("cannot render output:", err)
	}
}

func readFile(fn string) (string, error) {
	b, err := ioutil.ReadFile(fn)
	return string(b), err
}

const defaultRawTemplate = `
{{- block "info" .Info }}
Info
	Title: {{ .Title }}
	Description: {{ .Description }}
	TermsOfService: {{ .TermsOfService }}
	{{- with .Contact}}
	Contact: {{ .Name }} | {{ .URL }} | {{ .Email }}
	{{ end -}}
	{{- with .License }}
	License: {{ .Name }} | {{ .URL }}
	{{ end -}}
	Version: {{ .Version }}
{{ end -}}

{{- block "externalDocs" .ExternalDocs }}
External Docs: {{ .Description }} - {{ .URL }}
{{ end -}}

{{- block "schemes" .Schemes }}
Schemes: {{ range . }}{{ . }} {{ end }}
{{ end -}}

{{- block "basePath" .BasePath }}
BasePath: {{ . }}
{{ end -}}

{{- block "paths" .Paths }}
Paths
	{{ range $url, $item := . }}
		{{ $url }}:
			- Ref: {{ $item.Ref }}
			- Parameters: {{ $item.Parameters }}
			{{- with $item.Delete }}
			- Delete: {{ template "Operation" . -}}
			{{ end -}}
			{{- with $item.Get }}
			- Get: {{ template "Operation" . -}}
			{{ end -}}
			{{- with $item.Head }}
			- Head: {{ template "Operation" . -}}
			{{ end -}}
			{{- with $item.Options }}
			- Options: {{ template "Operation" . -}}
			{{ end -}}
			{{- with $item.Patch }}
			- Patch: {{ template "Operation" . -}}
			{{ end -}}
			{{- with $item.Post }}
			- Post: {{ template "Operation" . -}}
			{{ end -}}
			{{- with $item.Put }}
			- Put: {{ template "Operation" . -}}
			{{ end -}}
	{{ end }}
{{ end -}}

{{- define "Operation" }}
				Summary: {{ .Summary }}
				Description: {{ .Description }}
				ExternalDocs: {{ .ExternalDocs }}
				Tags: {{ .Tags }}
				OperationID: {{ .OperationID }}
				Parameters: {{ template "Parameters" .Parameters }}
				Responses: {{ template "Responses" .Responses }}
				Consumes: {{ .Consumes }}
				Produces: {{ .Produces }}
				Security: {{ .Security }}
{{ end -}}

{{- define "Parameters" }}
	{{- range . }}
					Parameter:
						- Ref: {{ .Ref }}
						- In: {{ .In }}
						- Name: {{ .Name }}
						- Description: {{ .Description }}
						- Required: {{ .Required }}
						- UniqueItems: {{ .UniqueItems }}
						- ExclusiveMin: {{ .ExclusiveMin }}
						- ExclusiveMax: {{ .ExclusiveMax }}
						- Schema: {{ .Schema }}
						- Type: {{ .Type }}
						- Format: {{ .Format }}
						- Enum: {{ .Enum }}
						- Minimum: {{ .Minimum }}
						- Maximum: {{ .Maximum }}
						- MinLength: {{ .MinLength }}
						- MaxLength: {{ .MaxLength }}
						- Pattern: {{ .Pattern }}
						- Default: {{ .Default }}
	{{ end -}}
{{ end -}}

{{- define "Responses" }}
	{{- range $key, $item := . }}
					Response {{ $key }}:
					- Ref: {{ $item.Ref }}
					- Description: {{ $item.Description }}
					- Schema: {{ $item.Schema }}
					- Headers: {{ $item.Headers }}
					- Examples: {{ $item.Examples }}
	{{ end -}}
{{ end -}}
`
