package main

import (
	"strings"
	"unicode"
)

// camelCaseToUnderscore copy from https://github.com/asaskevich/govalidator/blob/master/utils.go#L107-L119
// Ex: AbcDef => abc_dev
func camelCaseToUnderscore(str string) string {
	var output []rune
	var segment []rune
	for _, r := range str {

		// not treat number as separate segment
		if !unicode.IsLower(r) && string(r) != "_" && !unicode.IsNumber(r) {
			output = addSegment(output, segment)
			segment = nil
		}
		segment = append(segment, unicode.ToLower(r))
	}
	output = addSegment(output, segment)
	return string(output)
}

// underscoreToCamelCase copy from https://github.com/asaskevich/govalidator/blob/master/utils.go
// Ex.: my_func => MyFunc
func underscoreToCamelCase(s string) string {
	return strings.Replace(strings.Title(strings.Replace(s, "_", " ", -1)), " ", "", -1)
}

func addSegment(inrune, segment []rune) []rune {
	if len(segment) == 0 {
		return inrune
	}
	if len(inrune) != 0 {
		inrune = append(inrune, '_')
	}
	inrune = append(inrune, segment...)
	return inrune
}
