// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"strings"
)

func DatabaseCharset(v interface{}, k string) (warnings []string, errors []error) {
	value, ok := v.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	// https://www.postgresql.org/docs/13/multibyte.html#CHARSET-TABLE
	charsets := map[string]bool{
		"BIG5":           true,
		"EUC_CN":         true,
		"EUC_JP":         true,
		"EUC_JIS_2004":   true,
		"EUC_KR":         true,
		"EUC_TW":         true,
		"GB18030":        true,
		"GBK":            true,
		"ISO_8859_5":     true,
		"ISO_8859_6":     true,
		"ISO_8859_7":     true,
		"ISO_8859_8":     true,
		"JOHAB":          true,
		"KOI8R":          true,
		"KOI8U":          true,
		"LATIN1":         true,
		"LATIN2":         true,
		"LATIN3":         true,
		"LATIN4":         true,
		"LATIN5":         true,
		"LATIN6":         true,
		"LATIN7":         true,
		"LATIN8":         true,
		"LATIN9":         true,
		"LATIN10":        true,
		"MULE_INTERNAL":  true,
		"SJIS":           true,
		"SHIFT_JIS_2004": true,
		"SQL_ASCII":      true,
		"UHC":            true,
		"UTF8":           true,
		"WIN866":         true,
		"WIN874":         true,
		"WIN1250":        true,
		"WIN1251":        true,
		"WIN1252":        true,
		"WIN1253":        true,
		"WIN1254":        true,
		"WIN1255":        true,
		"WIN1256":        true,
		"WIN1257":        true,
		"WIN1258":        true,
	}

	if !charsets[strings.ToUpper(value)] {
		errors = append(errors, fmt.Errorf("%s contains unknown charset %s, see https://www.postgresql.org/docs/13/multibyte.html#CHARSET-TABLE for available charsets in PostgreSQL", k, v))
	}

	return warnings, errors
}
