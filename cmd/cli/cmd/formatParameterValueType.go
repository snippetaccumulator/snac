package cmd

import (
	"fmt"
	"strings"
)

type formatParameterValueType int

const (
	FormatParameterValueTypeDefault formatParameterValueType = iota
	FormatParameterValueTypeJSON
	FormatParameterValueTypeYAML
)

type formatParameterValue struct {
	format formatParameterValueType
}

func (f *formatParameterValue) String() string {
	switch f.format {
	case FormatParameterValueTypeJSON:
		return "json"
	case FormatParameterValueTypeYAML:
		return "yaml"
	default:
		return "default"
	}
}

func (f *formatParameterValue) Set(input string) error {
	lowercaseInput := strings.ToLower(input)
	switch lowercaseInput {
	case "json":
		f.format = FormatParameterValueTypeJSON
	case "yaml":
		f.format = FormatParameterValueTypeYAML
	case "default":
		f.format = FormatParameterValueTypeDefault
	default:
		return fmt.Errorf("invalid format: ''%s' (allowed values: 'json', 'yaml', 'default')", input)
	}
	return nil
}

func (f *formatParameterValue) Type() string {
	return "json/yaml"
}
