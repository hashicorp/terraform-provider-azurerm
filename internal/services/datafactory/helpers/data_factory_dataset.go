// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

func ExpandDataFactoryExpressionResultType(str string, isDynamic bool) any {
	if !isDynamic {
		return str
	}
	return map[string]string{
		"type":  "Expression",
		"value": str,
	}
}

func FlattenDataFactoryExpressionResultType(obj any) (result string, isDynamic bool) {
	switch v := obj.(type) {
	case string:
		result, isDynamic = v, false
	case map[string]any:
		isDynamic = true
		if value, ok := v["value"]; ok {
			result = value.(string)
		}
	}
	return
}
