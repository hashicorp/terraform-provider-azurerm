package validate

import (
	"fmt"
	"regexp"
	"strings"
)

func MetaDataKeys(value interface{}, _ string) (warnings []string, errors []error) {
	v, ok := value.(map[string]interface{})
	if !ok {
		return
	}

	for k := range v {
		isCSharpKeyword := cSharpKeywords[strings.ToLower(k)] != nil
		if isCSharpKeyword {
			errors = append(errors, fmt.Errorf("%q is not a valid key (C# keyword)", k))
		}

		// must begin with a letter, underscore
		// the rest: letters, digits and underscores
		if !regexp.MustCompile(`^([a-z_]{1}[a-z0-9_]{1,})$`).MatchString(k) {
			errors = append(errors, fmt.Errorf("MetaData must start with letters or an underscores and be all lowercase. Got %q.", k))
		}
	}

	return
}

var cSharpKeywords = map[string]*struct{}{
	"abstract":   {},
	"as":         {},
	"base":       {},
	"bool":       {},
	"break":      {},
	"byte":       {},
	"case":       {},
	"catch":      {},
	"char":       {},
	"checked":    {},
	"class":      {},
	"const":      {},
	"continue":   {},
	"decimal":    {},
	"default":    {},
	"delegate":   {},
	"do":         {},
	"double":     {},
	"else":       {},
	"enum":       {},
	"event":      {},
	"explicit":   {},
	"extern":     {},
	"false":      {},
	"finally":    {},
	"fixed":      {},
	"float":      {},
	"for":        {},
	"foreach":    {},
	"goto":       {},
	"if":         {},
	"implicit":   {},
	"in":         {},
	"int":        {},
	"interface":  {},
	"internal":   {},
	"is":         {},
	"lock":       {},
	"long":       {},
	"namespace":  {},
	"new":        {},
	"null":       {},
	"object":     {},
	"operator":   {},
	"out":        {},
	"override":   {},
	"params":     {},
	"private":    {},
	"protected":  {},
	"public":     {},
	"readonly":   {},
	"ref":        {},
	"return":     {},
	"sbyte":      {},
	"sealed":     {},
	"short":      {},
	"sizeof":     {},
	"stackalloc": {},
	"static":     {},
	"string":     {},
	"struct":     {},
	"switch":     {},
	"this":       {},
	"throw":      {},
	"true":       {},
	"try":        {},
	"typeof":     {},
	"uint":       {},
	"ulong":      {},
	"unchecked":  {},
	"unsafe":     {},
	"ushort":     {},
	"using":      {},
	"void":       {},
	"volatile":   {},
	"while":      {},
}
