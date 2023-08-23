// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
	"strings"
)

func PostgresqlDatabaseCollation(v interface{}, k string) (warnings []string, errors []error) {
	originalValue, ok := v.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if _, isSpecialCase := specialCases[originalValue]; isSpecialCase {
		return
	}

	value := strings.ToLower(originalValue)
	value = strings.TrimSuffix(value, ".utf8")

	// based on the string format, determine what we're validating
	if _, isLanguageCode := languageCodes[value]; isLanguageCode {
		return
	}

	// either `en-GB`, `en_GB`, ja_001` or `jv_java_id`
	languageCodeAndLocaleRegex := regexp.MustCompile("^[a-z]{1,}([-_]{1}([a-z0-9]{1,})){1,2}$")
	if !languageCodeAndLocaleRegex.MatchString(value) {
		errors = append(errors, databaseCollationDidNotMatchError(k, originalValue))
	}

	containsDash := strings.Contains(value, "-")
	containsUnderscore := strings.Contains(value, "_")
	if containsDash && containsUnderscore {
		errors = append(errors, databaseCollationDidNotMatchError(k, originalValue))
		return
	}

	var split []string
	if containsDash {
		split = strings.Split(value, "-")
	} else {
		split = strings.Split(value, "_")
	}
	if len(split) == 0 {
		errors = append(errors, databaseCollationDidNotMatchError(k, originalValue))
		return
	}

	// validate the language code is valid
	languageCode := split[0]
	if _, languageCodeIsValid := languageCodes[languageCode]; !languageCodeIsValid {
		errors = append(errors, databaseCollationDidNotMatchError(k, originalValue))
		return
	}

	// we can't do much about the locale, so we'll assume that's fine for now

	return warnings, errors
}

var languageCodes = map[string]struct{}{
	"aa":  {},
	"af":  {},
	"agq": {},
	"ak":  {},
	"am":  {},
	"ar":  {},
	"arn": {},
	"as":  {},
	"asa": {},
	"ast": {},
	"bas": {},
	"bem": {},
	"bez": {},
	"bin": {},
	"bm":  {},
	"bn":  {},
	"bo":  {},
	"br":  {},
	"brx": {},
	"byn": {},
	"ca":  {},
	"ce":  {},
	"cgg": {},
	"chr": {},
	"co":  {},
	"cu":  {},
	"cy":  {},
	"da":  {},
	"dav": {},
	"de":  {},
	"dje": {},
	"dsb": {},
	"dua": {},
	"dv":  {},
	"dyo": {},
	"dz":  {},
	"ebu": {},
	"ee":  {},
	"el":  {},
	"en":  {},
	"eo":  {},
	"es":  {},
	"eu":  {},
	"ewo": {},
	"ff":  {},
	"fi":  {},
	"fil": {},
	"fo":  {},
	"fr":  {},
	"fur": {},
	"fy":  {},
	"ga":  {},
	"gd":  {},
	"gl":  {},
	"gsw": {},
	"gu":  {},
	"guz": {},
	"gv":  {},
	"ha":  {},
	"haw": {},
	"hi":  {},
	"hsb": {},
	"hy":  {},
	"ia":  {},
	"ibb": {},
	"id":  {},
	"ig":  {},
	"ii":  {},
	"is":  {},
	"it":  {},
	"iu":  {},
	"ja":  {},
	"jgo": {},
	"jmc": {},
	"jv":  {},
	"ka":  {},
	"kab": {},
	"kam": {},
	"kde": {},
	"kea": {},
	"khq": {},
	"ki":  {},
	"kk":  {},
	"kkj": {},
	"kl":  {},
	"kln": {},
	"km":  {},
	"kn":  {},
	"ko":  {},
	"kok": {},
	"kr":  {},
	"ks":  {},
	"ksb": {},
	"ksf": {},
	"ksh": {},
	"ku":  {},
	"kw":  {},
	"la":  {},
	"lag": {},
	"lb":  {},
	"lg":  {},
	"lkt": {},
	"ln":  {},
	"lo":  {},
	"lrc": {},
	"lu":  {},
	"luo": {},
	"luy": {},
	"mas": {},
	"mer": {},
	"mfe": {},
	"mgh": {},
	"mgo": {},
	"mi":  {},
	"ml":  {},
	"mn":  {},
	"mni": {},
	"moh": {},
	"mr":  {},
	"ms":  {},
	"mt":  {},
	"mua": {},
	"mzn": {},
	"naq": {},
	"nb":  {},
	"nd":  {},
	"nds": {},
	"ne":  {},
	"nl":  {},
	"nmg": {},
	"nn":  {},
	"nnh": {},
	"no":  {},
	"nr":  {},
	"nso": {},
	"nus": {},
	"nyn": {},
	"oc":  {},
	"om":  {},
	"or":  {},
	"os":  {},
	"pa":  {},
	"pap": {},
	"prg": {},
	"ps":  {},
	"pt":  {},
	"quc": {},
	"quz": {},
	"rm":  {},
	"rn":  {},
	"rof": {},
	"ru":  {},
	"rw":  {},
	"rwk": {},
	"sa":  {},
	"saq": {},
	"sbp": {},
	"sd":  {},
	"se":  {},
	"seh": {},
	"ses": {},
	"sg":  {},
	"shi": {},
	"si":  {},
	"sma": {},
	"smj": {},
	"smn": {},
	"sms": {},
	"so":  {},
	"sq":  {},
	"sr":  {},
	"ss":  {},
	"ssy": {},
	"st":  {},
	"sv":  {},
	"sw":  {},
	"syr": {},
	"ta":  {},
	"te":  {},
	"teo": {},
	"ti":  {},
	"tig": {},
	"tn":  {},
	"to":  {},
	"tr":  {},
	"twq": {},
	"tzm": {},
	"uk":  {},
	"uz":  {},
	"vai": {},
	"ve":  {},
	"vo":  {},
	"vun": {},
	"wae": {},
	"wal": {},
	"wo":  {},
	"xh":  {},
	"xog": {},
	"yav": {},
	"yi":  {},
	"yo":  {},
	"zh":  {},
	"zu":  {},
}

var specialCases = map[string]struct{}{
	".utf8":                      {},
	"C":                          {},
	"POSIX":                      {},
	"English_United States.1252": {},
	"ucs_basic":                  {},
	"default":                    {},
}

var databaseCollationDidNotMatchError = func(fieldName, value string) error {
	specialCasedValues := make([]string, 0)
	for specialCasedValue := range specialCases {
		specialCasedValues = append(specialCasedValues, fmt.Sprintf("%q", specialCasedValue))
	}

	return fmt.Errorf(`the value %[2]q is not valid for the field %[1]q, expected either:

* A language code (e.g. 'de', 'en' or 'ja')
* A language code and locale with a dash (e.g. 'de-DE', 'en-GB' or 'ja-JP')
* A language code and locale with an underscore (e.g. 'de_DE', 'en_GB' or 'ja_JP')
* One of the following values: %[3]s

This value can optionally end in '.utf8'

but got %[2]q
`, fieldName, value, strings.Join(specialCasedValues, ", "))
}
