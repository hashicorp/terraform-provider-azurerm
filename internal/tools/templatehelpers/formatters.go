// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package templatehelpers

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var TplFuncMap = template.FuncMap{
	"ToLower":                       strings.ToLower,
	"ToTitle":                       ToTitle,
	"ToCamel":                       strcase.ToCamel,
	"ToSnake":                       strcase.ToSnake,
	"TfName":                        TerraformResourceName,
	"ToString":                      ToString,
	"ToDelim":                       strcase.ToDelimited,
	"ToDelimTitle":                  ToDelimTitle,
	"PrefixedDescriptionString":     PrefixedDescriptionString,
	"PrefixedLabelString":           PrefixedLabelString,
	"IdToID":                        IdToID,
	"NewIDResourceIdentityFormater": NewIDResourceIdentityFormater,
	"NewIDCreateFormater":           NewIDCreateFormater,
	"ClientToPackageName":           ClientToPackageName,
}

// TerraformResourceName generates a Terraform-compliant resource name by combining the provider and resource name.
func TerraformResourceName(provider, resourceName string) string {
	fmtStr := "%s_%s"
	return fmt.Sprintf(fmtStr, strings.ToLower(provider), strcase.ToSnake(resourceName))
}

func ToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	default:
		return ""
	}
}

func ToTitle(input string) string {
	return cases.Title(language.BritishEnglish).String(input)
}

// PrefixedDescriptionString returns a string prefixed with "a" or "an" based on whether the input starts with a vowel.
func PrefixedDescriptionString(input string) string {
	prefix := "a"
	first := input[0:1]
	vowel, _ := regexp.Match(first, []byte(`aeiouAEIOU`))

	if vowel {
		prefix = "an"
	}
	return fmt.Sprintf("%s %s", prefix, cases.Title(language.English).String(strcase.ToDelimited(input, ' ')))
}

// ToDelimTitle converts the input string to a title-cased string with words delimited by spaces.
func ToDelimTitle(input string) string {
	return cases.Title(language.English).String(strcase.ToDelimited(input, ' '))
}

// PrefixedLabelString determines whether a given label should use "A" or "An" as its prefix based on its starting letter.
func PrefixedLabelString(input string) string {
	prefix := "A"
	vowel, _ := regexp.Match(input[0:1], []byte(`aeiouAEIOU`))

	if vowel {
		prefix = "An"
	}

	return fmt.Sprintf("%s `%s`", prefix, input)
}

func IdToID(id string) string {
	if strings.HasSuffix(id, "Id") {
		id = strings.TrimSuffix(id, "d")
		id += "D"
	}

	return id
}

func NewIDResourceIdentityFormater(idType []string, idSegments []string, prefix string) string {
	if len(idType) != 2 {
		return "// TODO - ID Type provided to scaff tool did not match the expected format"
	}

	f := "%s.New%s(%s)"
	out := make([]string, 0)
	out = append(out, idSegments...)

	output := fmt.Sprintf(f, idType[0], IdToID(idType[1]), strings.Join(out, ", "))

	return output
}

func NewIDCreateFormater(idType []string, idSegments []string, prefix string) string {
	if len(idType) != 2 {
		return "// TODO - ID Type provided to scaff tool did not match the expected format"
	}

	f := "%s.New%s(%s)"
	out := make([]string, 0)
	out = append(out, idSegments...)

	if slices.Contains(out, "SubscriptionId") {
		f = "%s.New%s(metadata.Client.Account.SubscriptionId, %s)"
		out = slices.Delete(out, slices.Index(out, "SubscriptionId"), 1)
		out = out[:len(out)-1]
	}

	if slices.Contains(out, "ResourceGroup") {
		idx := slices.Index(out, "ResourceGroup")
		out = slices.Replace(out, idx, idx, "ResourceGroupName")
	}

	for i, v := range out {
		out = slices.Replace(out, i, i, fmt.Sprintf("%s.%s.ValueString()", prefix, v))
	}

	return fmt.Sprintf(f, idType[0], IdToID(idType[1]), strings.Join(out, ", "))
}

func ClientToPackageName(input string) string {
	return strings.ToLower(strings.TrimSuffix(input, "Client"))
}
