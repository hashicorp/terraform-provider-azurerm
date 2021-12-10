package datafactory

func expandDataFactoryExpressionResultType(str string, isDynamic bool) interface{} {
	if !isDynamic {
		return str
	}
	return map[string]string{
		"type":  "Expression",
		"value": str,
	}
}

func flattenDataFactoryExpressionResultType(obj interface{}) (result string, isDynamic bool) {
	switch v := obj.(type) {
	case string:
		result, isDynamic = v, false
	case map[string]interface{}:
		isDynamic = true
		if value, ok := v["value"]; ok {
			result = value.(string)
		}
	}
	return
}
