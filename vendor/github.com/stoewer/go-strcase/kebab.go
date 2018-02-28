// Copyright (c) 2017, A. Stoewer <adrian.stoewer@rz.ifi.lmu.de>
// All rights reserved.

package strcase

// KebabCase converts a string into kebab case.
func KebabCase(s string) string {
	return lowerDelimiterCase(s, '-')
}
