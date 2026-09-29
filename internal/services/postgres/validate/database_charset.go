// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// https://www.postgresql.org/docs/13/multibyte.html#CHARSET-TABLE
var charsets = []string{
	"BIG5",
	"EUC_CN",
	"EUC_JP",
	"EUC_JIS_2004",
	"EUC_KR",
	"EUC_TW",
	"GB18030",
	"GBK",
	"ISO_8859_5",
	"ISO_8859_6",
	"ISO_8859_7",
	"ISO_8859_8",
	"JOHAB",
	"KOI8R",
	"KOI8U",
	"LATIN1",
	"LATIN2",
	"LATIN3",
	"LATIN4",
	"LATIN5",
	"LATIN6",
	"LATIN7",
	"LATIN8",
	"LATIN9",
	"LATIN10",
	"MULE_INTERNAL",
	"SJIS",
	"SHIFT_JIS_2004",
	"SQL_ASCII",
	"UHC",
	"UTF8",
	"WIN866",
	"WIN874",
	"WIN1250",
	"WIN1251",
	"WIN1252",
	"WIN1253",
	"WIN1254",
	"WIN1255",
	"WIN1256",
	"WIN1257",
	"WIN1258",
}

func DatabaseCharset(v interface{}, k string) ([]string, []error) {
	return validation.StringInSlice(charsets, true)(v, k)
}
