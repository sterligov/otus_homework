package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strings"
)

var ErrEmptyParams = errors.New("empty params")

var tagReg = regexp.MustCompile("validate:\"(.*?)\"(\\s|`)")

type (
	Validator struct {
		Name   string
		Params []string
	}

	Field struct {
		Name       string
		Type       string
		IsArray    bool
		Validators []*Validator
	}

	StructType struct {
		Name    string
		Package string
		Fields  []*Field
	}
)

func parseStructs(filename string) ([]*StructType, error) {
	fset := token.NewFileSet()
	fast, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parse file failed: %w", err)
	}

	structs := make([]*StructType, 0)
	var inspectError error

	ast.Inspect(fast, func(x ast.Node) bool {
		ts, ok := x.(*ast.TypeSpec)
		if !ok {
			return true
		}

		s, ok := ts.Type.(*ast.StructType)
		if !ok {
			return true
		}

		fields, err := parseStructFields(s)
		if err != nil {
			inspectError = err
		}

		if len(fields) > 0 {
			st := &StructType{
				Name:    ts.Name.Name,
				Fields:  fields,
				Package: fast.Name.Name,
			}
			structs = append(structs, st)
		}

		return false
	})

	if inspectError != nil {
		return nil, inspectError
	}

	return structs, nil
}

func parseStructFields(curStruct *ast.StructType) ([]*Field, error) {
	if curStruct.Fields == nil {
		return nil, nil
	}

	fields := make([]*Field, 0)

	for _, f := range curStruct.Fields.List {
		if f.Tag == nil || len(f.Names) == 0 || f.Tag.Value == "" {
			continue
		}

		validators, err := parseTag(f.Tag.Value)
		if err != nil {
			return nil, err
		}
		if len(validators) == 0 {
			continue
		}

		field := &Field{
			Name:       f.Names[0].Name,
			Validators: validators,
		}

		expr := f.Type

	loop:
		for {
			switch fieldType := expr.(type) {
			case *ast.ArrayType:
				field.IsArray = true
				expr = fieldType.Elt
			case *ast.Ident:
				if fieldType.Obj != nil {
					if ts, ok := fieldType.Obj.Decl.(*ast.TypeSpec); ok {
						expr = ts.Type.(*ast.Ident)
						continue
					}
				}
				field.Type = fieldType.Name
				break loop
			}
		}

		fields = append(fields, field)
	}

	return fields, nil
}

func parseTag(tag string) ([]*Validator, error) {
	matches := tagReg.FindStringSubmatch(tag)
	if len(matches) != 3 {
		return nil, nil
	}

	vals := strings.Split(matches[1], "|")
	validators := make([]*Validator, 0)

	for _, v := range vals {
		validatorValue := strings.SplitN(v, ":", 2)
		if len(validatorValue[1]) == 0 {
			return nil, fmt.Errorf("%s: %w", validatorValue[0], ErrEmptyParams)
		}

		validators = append(validators, &Validator{
			Name:   validatorValue[0],
			Params: strings.Split(validatorValue[1], ","),
		})
	}

	return validators, nil
}
