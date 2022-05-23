package apps

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AppSku string

const (
	AppSkuFOne   AppSku = "F1"
	AppSkuSOne   AppSku = "S1"
	AppSkuSTOne  AppSku = "ST1"
	AppSkuSTTwo  AppSku = "ST2"
	AppSkuSTZero AppSku = "ST0"
)

func PossibleValuesForAppSku() []string {
	return []string{
		string(AppSkuFOne),
		string(AppSkuSOne),
		string(AppSkuSTOne),
		string(AppSkuSTTwo),
		string(AppSkuSTZero),
	}
}

func parseAppSku(input string) (*AppSku, error) {
	vals := map[string]AppSku{
		"f1":  AppSkuFOne,
		"s1":  AppSkuSOne,
		"st1": AppSkuSTOne,
		"st2": AppSkuSTTwo,
		"st0": AppSkuSTZero,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AppSku(input)
	return &out, nil
}
