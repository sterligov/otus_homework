package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/tools/imports"
)

var ErrUnexpectedValidator = errors.New("unexpected validator")

const (
	_minValidator      = "min"
	_maxValidator      = "max"
	_inIntValidator    = "in_int"
	_inStringValidator = "in_string"
	_regexpValidator   = "regexp"
	_lenValidator      = "len"
	_arrayValidator    = "array"
)

const _validatorFuncGenerator = "validate_func_gen"

type GenValidatorFunc func(fieldName string, params []string) (string, error)

type Generator struct {
	template   *template.Template
	generators map[string]GenValidatorFunc
}

func NewGenerator() *Generator {
	funcs := template.FuncMap{
		"join":         strings.Join,
		"getFieldName": getFieldName,
		"replace":      strings.Replace,
	}
	t := template.New("validator").Funcs(funcs)

	g := &Generator{
		template: t,
	}

	g.registerGenerators()

	return g
}

func (g *Generator) registerGenerators() {
	g.generators = map[string]GenValidatorFunc{
		_maxValidator:           g.maxValidator(),
		_minValidator:           g.minValidator(),
		_inIntValidator:         g.inIntValidator(),
		_inStringValidator:      g.inStringValidator(),
		_lenValidator:           g.lenValidator(),
		_regexpValidator:        g.regexpValidator(),
		_arrayValidator:         g.arrayValidator(),
		_validatorFuncGenerator: g.genValidatorFuncBody(),
	}
}

func (g *Generator) genStructValidators(s *StructType) (string, error) {
	validatorsCode := make([]string, 0, 10)
	for _, f := range s.Fields {
		for _, v := range f.Validators {
			code, err := g.generateValidatorCode(v, f)
			if err != nil {
				return "", err
			}
			validatorsCode = append(validatorsCode, code)
		}
	}

	return g.generators[_validatorFuncGenerator](s.Name, validatorsCode)
}

func (g *Generator) generateValidatorCode(v *Validator, f *Field) (string, error) {
	var (
		code      string
		err       error
		fieldName = "this." + f.Name
	)

	if vfunc, ok := g.generators[v.Name+"_"+f.Type]; ok {
		code, err = vfunc(fieldName, v.Params)
	} else if vfunc, ok := g.generators[v.Name]; ok {
		code, err = vfunc(fieldName, v.Params)
	} else {
		err = fmt.Errorf("%s: %w", v.Name, ErrUnexpectedValidator)
	}

	if err != nil {
		return "", err
	}

	if f.IsArray {
		return g.generators[_arrayValidator](fieldName, []string{code})
	}

	return code, nil
}

func (g *Generator) genValidatorFileBody() func(pack string, structValidators []string) ([]byte, error) {
	templ := `
// Code generated by cool valid tool; DO NOT EDIT.
package {{.package}}

type ValidationError struct {
	Field string
	Err string
}

{{- range .structValidators}}
	{{.}}
{{- end}}`

	return func(pack string, structValidators []string) ([]byte, error) {
		return executeInSlice(g.cloneAndParse(templ), map[string]interface{}{
			"package":          pack,
			"structValidators": structValidators,
		})
	}
}

func (g *Generator) genValidatorFuncBody() GenValidatorFunc {
	templ := `
func (this {{.structName}}) Validate() ([]ValidationError, error) {
	ve := make([]ValidationError, 0)
	{{- range .validators}}
		{{. -}}
	{{end}}
	
	return ve, nil
}`

	return func(structName string, validators []string) (string, error) {
		return executeInString(g.cloneAndParse(templ), map[string]interface{}{
			"structName": structName,
			"validators": validators,
		})
	}
}

func (g *Generator) lenValidator() GenValidatorFunc {
	templ := `
if len({{.fieldName}}) != {{.value}} {
	err := ValidationError{
		Field: "{{getFieldName .fieldName}}",
		Err: "Len must be equal to {{.value}}",
	}
	ve = append(ve, err)
}`

	return g.oneParamGeneratorFunc(templ)
}

func (g *Generator) inIntValidator() GenValidatorFunc {
	templ := `
if {{$delim := .fieldName | printf " != %s && "}} {{join .values $delim}} != {{.fieldName}} {
	err := ValidationError{
		Field: "{{getFieldName .fieldName}}",
		Err: "Value must be in list [{{join .values ","}}]",
	}
	ve = append(ve, err)
}`

	return g.allParamsGeneratorFunc(templ)
}

func (g *Generator) inStringValidator() GenValidatorFunc {
	templ := `
if {{$delim := .fieldName | printf "\" != %s && \""}} "{{join .values $delim}}" != {{.fieldName}} {
	err := ValidationError{
		Field: "{{getFieldName .fieldName}}",
		Err: "Value must be in list [{{join .values ","}}]",
	}
	ve = append(ve, err)
}`

	return g.allParamsGeneratorFunc(templ)
}

func (g *Generator) minValidator() GenValidatorFunc {
	templ := `
if {{.fieldName}} < {{.value}} {
	err := ValidationError{
		Field: "{{getFieldName .fieldName}}",
		Err: "Value must be greater or equal than {{.value}}",
	}
	ve = append(ve, err)
}`

	return g.oneParamGeneratorFunc(templ)
}

func (g *Generator) maxValidator() GenValidatorFunc {
	templ := `
if {{.fieldName}} > {{.value}} {
	err := ValidationError{
		Field: "{{getFieldName .fieldName}}",
		Err: "Value must be less or equal to {{.value}}",
	}
	ve = append(ve, err)
}`

	return g.oneParamGeneratorFunc(templ)
}

func (g *Generator) regexpValidator() GenValidatorFunc {
	templ := `
if match, err := regexp.MatchString("{{.value}}", {{.fieldName}}); err != nil {
	return nil, err
} else if !match {
	err := ValidationError{
		Field: "{{getFieldName .fieldName}}",
		Err: "Value must satisfy regular expression {{.value}}",
	}
	ve = append(ve, err)
}`

	return g.oneParamGeneratorFunc(templ)
}

func (g *Generator) arrayValidator() GenValidatorFunc {
	templ := `
for _, v := range {{.fieldName}} {
	{{- replace .value .fieldName "v" -1}}
}`

	return g.oneParamGeneratorFunc(templ)
}

func (g *Generator) oneParamGeneratorFunc(templ string) GenValidatorFunc {
	t := g.cloneAndParse(templ)

	return func(fieldName string, params []string) (string, error) {
		return executeInString(t, map[string]interface{}{
			"fieldName": fieldName,
			"value":     params[0],
		})
	}
}

func (g *Generator) allParamsGeneratorFunc(templ string) GenValidatorFunc {
	t := g.cloneAndParse(templ)

	return func(fieldName string, params []string) (string, error) {
		return executeInString(t, map[string]interface{}{
			"fieldName": fieldName,
			"values":    params,
		})
	}
}

func (g *Generator) cloneAndParse(templ string) *template.Template {
	t := template.Must(g.template.Clone())

	return template.Must(t.Parse(templ))
}

func GenerateValidators(filename string) error {
	structs, err := parseStructs(filename)
	if err != nil {
		return err
	}

	if len(structs) == 0 {
		return nil
	}

	g := NewGenerator()
	structValidators := make([]string, 0)

	for _, s := range structs {
		code, err := g.genStructValidators(s)
		if err != nil {
			return err
		}

		structValidators = append(structValidators, code)
	}

	code, err := g.genValidatorFileBody()(structs[0].Package, structValidators)
	if err != nil {
		return err
	}

	code, err = imports.Process("", code, nil)
	if err != nil {
		return fmt.Errorf("goimports failed: %w", err)
	}

	validatorFilename := getValidatorFilename(filename)
	if err := ioutil.WriteFile(validatorFilename, code, os.FileMode(0755)); err != nil {
		return fmt.Errorf("write file failed: %w", err)
	}

	return nil
}

func executeInString(t *template.Template, data interface{}) (string, error) {
	buf, err := executeInBuffer(t, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func executeInSlice(t *template.Template, data interface{}) ([]byte, error) {
	buf, err := executeInBuffer(t, data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func executeInBuffer(t *template.Template, data interface{}) (*bytes.Buffer, error) {
	buf := make([]byte, 0, 1024)
	buffer := bytes.NewBuffer(buf)

	if err := t.Execute(buffer, data); err != nil {
		return nil, fmt.Errorf("execute template failed: %w", err)
	}

	return buffer, nil
}

func getFieldName(name string) string {
	s := strings.Split(name, ".")
	if len(s) < 2 {
		return ""
	}

	return s[1]
}

func getValidatorFilename(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename)) + "_validation_generated.go"
}