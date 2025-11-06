// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-04-01/redisenterprise"
)

func possibleValuesForAofFrequency() []string {
	res := make([]string, 0, len(redisenterprise.PossibleValuesForAofFrequency())-1)
	for _, freq := range redisenterprise.PossibleValuesForAofFrequency() {
		// `always` has been deprecated but not yet marked on the OpenAPI spec. It is no longer listed in the docs / portal:
		// https://learn.microsoft.com/azure/redis/how-to-persistence
		if freq != string(redisenterprise.AofFrequencyAlways) {
			res = append(res, freq)
		}
	}
	return res
}

func AofBackupFrequency(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(int)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %s to be int", k)}
	}

	validValues := []int{}
	for _, freq := range possibleValuesForAofFrequency() {
		freqDuration, err := time.ParseDuration(freq)
		if err != nil {
			return nil, []error{fmt.Errorf("unable to parse AOF frequency duration from SDK %q: %s, this is likely due to SDK update", freq, err)}
		}
		validValues = append(validValues, int(freqDuration.Seconds()))
		if v == int(freqDuration.Seconds()) {
			return nil, nil
		}
	}
	slices.Sort(validValues)
	return nil, []error{fmt.Errorf("expected %q to be one of [%s], got %d", k, intsJoin(validValues, ", "), v)}
}

func RdbBackupFrequency(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(int)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %s to be int", k)}
	}

	validValues := []int{}
	for _, freq := range redisenterprise.PossibleValuesForRdbFrequency() {
		freqDuration, err := time.ParseDuration(freq)
		if err != nil {
			return nil, []error{fmt.Errorf("unable to parse RDB frequency duration from SDK %q: %s, this is likely due to SDK update", freq, err)}
		}
		validValues = append(validValues, int(freqDuration.Hours()))
		if v == int(freqDuration.Hours()) {
			return nil, nil
		}
	}
	slices.Sort(validValues)
	return nil, []error{fmt.Errorf("expected %q to be one of [%s], got %d", k, intsJoin(validValues, ", "), v)}
}

func intsJoin(ints []int, sep string) string {
	strs := make([]string, len(ints))
	for i, v := range ints {
		strs[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(strs, sep)
}
