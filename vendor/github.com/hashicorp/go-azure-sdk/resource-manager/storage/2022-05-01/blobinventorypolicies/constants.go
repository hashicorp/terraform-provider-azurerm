package blobinventorypolicies

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Format string

const (
	FormatCsv     Format = "Csv"
	FormatParquet Format = "Parquet"
)

func PossibleValuesForFormat() []string {
	return []string{
		string(FormatCsv),
		string(FormatParquet),
	}
}

func parseFormat(input string) (*Format, error) {
	vals := map[string]Format{
		"csv":     FormatCsv,
		"parquet": FormatParquet,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Format(input)
	return &out, nil
}

type InventoryRuleType string

const (
	InventoryRuleTypeInventory InventoryRuleType = "Inventory"
)

func PossibleValuesForInventoryRuleType() []string {
	return []string{
		string(InventoryRuleTypeInventory),
	}
}

func parseInventoryRuleType(input string) (*InventoryRuleType, error) {
	vals := map[string]InventoryRuleType{
		"inventory": InventoryRuleTypeInventory,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := InventoryRuleType(input)
	return &out, nil
}

type ObjectType string

const (
	ObjectTypeBlob      ObjectType = "Blob"
	ObjectTypeContainer ObjectType = "Container"
)

func PossibleValuesForObjectType() []string {
	return []string{
		string(ObjectTypeBlob),
		string(ObjectTypeContainer),
	}
}

func parseObjectType(input string) (*ObjectType, error) {
	vals := map[string]ObjectType{
		"blob":      ObjectTypeBlob,
		"container": ObjectTypeContainer,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ObjectType(input)
	return &out, nil
}

type Schedule string

const (
	ScheduleDaily  Schedule = "Daily"
	ScheduleWeekly Schedule = "Weekly"
)

func PossibleValuesForSchedule() []string {
	return []string{
		string(ScheduleDaily),
		string(ScheduleWeekly),
	}
}

func parseSchedule(input string) (*Schedule, error) {
	vals := map[string]Schedule{
		"daily":  ScheduleDaily,
		"weekly": ScheduleWeekly,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Schedule(input)
	return &out, nil
}
