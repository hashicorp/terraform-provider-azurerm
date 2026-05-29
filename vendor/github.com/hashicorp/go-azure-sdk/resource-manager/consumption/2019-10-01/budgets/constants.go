package budgets

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BudgetOperatorType string

const (
	BudgetOperatorTypeIn BudgetOperatorType = "In"
)

func PossibleValuesForBudgetOperatorType() []string {
	return []string{
		string(BudgetOperatorTypeIn),
	}
}

func (s *BudgetOperatorType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBudgetOperatorType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBudgetOperatorType(input string) (*BudgetOperatorType, error) {
	vals := map[string]BudgetOperatorType{
		"in": BudgetOperatorTypeIn,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BudgetOperatorType(input)
	return &out, nil
}

type CategoryType string

const (
	CategoryTypeCost CategoryType = "Cost"
)

func PossibleValuesForCategoryType() []string {
	return []string{
		string(CategoryTypeCost),
	}
}

func (s *CategoryType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCategoryType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCategoryType(input string) (*CategoryType, error) {
	vals := map[string]CategoryType{
		"cost": CategoryTypeCost,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CategoryType(input)
	return &out, nil
}

type CultureCode string

const (
	CultureCodeCsNegativecz CultureCode = "cs-cz"
	CultureCodeDaNegativedk CultureCode = "da-dk"
	CultureCodeDeNegativede CultureCode = "de-de"
	CultureCodeEnNegativegb CultureCode = "en-gb"
	CultureCodeEnNegativeus CultureCode = "en-us"
	CultureCodeEsNegativees CultureCode = "es-es"
	CultureCodeFrNegativefr CultureCode = "fr-fr"
	CultureCodeHuNegativehu CultureCode = "hu-hu"
	CultureCodeItNegativeit CultureCode = "it-it"
	CultureCodeJaNegativejp CultureCode = "ja-jp"
	CultureCodeKoNegativekr CultureCode = "ko-kr"
	CultureCodeNbNegativeno CultureCode = "nb-no"
	CultureCodeNlNegativenl CultureCode = "nl-nl"
	CultureCodePlNegativepl CultureCode = "pl-pl"
	CultureCodePtNegativebr CultureCode = "pt-br"
	CultureCodePtNegativept CultureCode = "pt-pt"
	CultureCodeRuNegativeru CultureCode = "ru-ru"
	CultureCodeSvNegativese CultureCode = "sv-se"
	CultureCodeTrNegativetr CultureCode = "tr-tr"
	CultureCodeZhNegativecn CultureCode = "zh-cn"
	CultureCodeZhNegativetw CultureCode = "zh-tw"
)

func PossibleValuesForCultureCode() []string {
	return []string{
		string(CultureCodeCsNegativecz),
		string(CultureCodeDaNegativedk),
		string(CultureCodeDeNegativede),
		string(CultureCodeEnNegativegb),
		string(CultureCodeEnNegativeus),
		string(CultureCodeEsNegativees),
		string(CultureCodeFrNegativefr),
		string(CultureCodeHuNegativehu),
		string(CultureCodeItNegativeit),
		string(CultureCodeJaNegativejp),
		string(CultureCodeKoNegativekr),
		string(CultureCodeNbNegativeno),
		string(CultureCodeNlNegativenl),
		string(CultureCodePlNegativepl),
		string(CultureCodePtNegativebr),
		string(CultureCodePtNegativept),
		string(CultureCodeRuNegativeru),
		string(CultureCodeSvNegativese),
		string(CultureCodeTrNegativetr),
		string(CultureCodeZhNegativecn),
		string(CultureCodeZhNegativetw),
	}
}

func (s *CultureCode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCultureCode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCultureCode(input string) (*CultureCode, error) {
	vals := map[string]CultureCode{
		"cs-cz": CultureCodeCsNegativecz,
		"da-dk": CultureCodeDaNegativedk,
		"de-de": CultureCodeDeNegativede,
		"en-gb": CultureCodeEnNegativegb,
		"en-us": CultureCodeEnNegativeus,
		"es-es": CultureCodeEsNegativees,
		"fr-fr": CultureCodeFrNegativefr,
		"hu-hu": CultureCodeHuNegativehu,
		"it-it": CultureCodeItNegativeit,
		"ja-jp": CultureCodeJaNegativejp,
		"ko-kr": CultureCodeKoNegativekr,
		"nb-no": CultureCodeNbNegativeno,
		"nl-nl": CultureCodeNlNegativenl,
		"pl-pl": CultureCodePlNegativepl,
		"pt-br": CultureCodePtNegativebr,
		"pt-pt": CultureCodePtNegativept,
		"ru-ru": CultureCodeRuNegativeru,
		"sv-se": CultureCodeSvNegativese,
		"tr-tr": CultureCodeTrNegativetr,
		"zh-cn": CultureCodeZhNegativecn,
		"zh-tw": CultureCodeZhNegativetw,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CultureCode(input)
	return &out, nil
}

type OperatorType string

const (
	OperatorTypeEqualTo              OperatorType = "EqualTo"
	OperatorTypeGreaterThan          OperatorType = "GreaterThan"
	OperatorTypeGreaterThanOrEqualTo OperatorType = "GreaterThanOrEqualTo"
)

func PossibleValuesForOperatorType() []string {
	return []string{
		string(OperatorTypeEqualTo),
		string(OperatorTypeGreaterThan),
		string(OperatorTypeGreaterThanOrEqualTo),
	}
}

func (s *OperatorType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOperatorType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOperatorType(input string) (*OperatorType, error) {
	vals := map[string]OperatorType{
		"equalto":              OperatorTypeEqualTo,
		"greaterthan":          OperatorTypeGreaterThan,
		"greaterthanorequalto": OperatorTypeGreaterThanOrEqualTo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OperatorType(input)
	return &out, nil
}

type ThresholdType string

const (
	ThresholdTypeActual ThresholdType = "Actual"
)

func PossibleValuesForThresholdType() []string {
	return []string{
		string(ThresholdTypeActual),
	}
}

func (s *ThresholdType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseThresholdType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseThresholdType(input string) (*ThresholdType, error) {
	vals := map[string]ThresholdType{
		"actual": ThresholdTypeActual,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ThresholdType(input)
	return &out, nil
}

type TimeGrainType string

const (
	TimeGrainTypeAnnually       TimeGrainType = "Annually"
	TimeGrainTypeBillingAnnual  TimeGrainType = "BillingAnnual"
	TimeGrainTypeBillingMonth   TimeGrainType = "BillingMonth"
	TimeGrainTypeBillingQuarter TimeGrainType = "BillingQuarter"
	TimeGrainTypeMonthly        TimeGrainType = "Monthly"
	TimeGrainTypeQuarterly      TimeGrainType = "Quarterly"
)

func PossibleValuesForTimeGrainType() []string {
	return []string{
		string(TimeGrainTypeAnnually),
		string(TimeGrainTypeBillingAnnual),
		string(TimeGrainTypeBillingMonth),
		string(TimeGrainTypeBillingQuarter),
		string(TimeGrainTypeMonthly),
		string(TimeGrainTypeQuarterly),
	}
}

func (s *TimeGrainType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTimeGrainType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTimeGrainType(input string) (*TimeGrainType, error) {
	vals := map[string]TimeGrainType{
		"annually":       TimeGrainTypeAnnually,
		"billingannual":  TimeGrainTypeBillingAnnual,
		"billingmonth":   TimeGrainTypeBillingMonth,
		"billingquarter": TimeGrainTypeBillingQuarter,
		"monthly":        TimeGrainTypeMonthly,
		"quarterly":      TimeGrainTypeQuarterly,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TimeGrainType(input)
	return &out, nil
}
