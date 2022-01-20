package application

import "strings"

type FailureAction string

const (
	FailureActionManual   FailureAction = "Manual"
	FailureActionRollback FailureAction = "Rollback"
)

func PossibleValuesForFailureAction() []string {
	return []string{
		string(FailureActionManual),
		string(FailureActionRollback),
	}
}

func parseFailureAction(input string) (*FailureAction, error) {
	vals := map[string]FailureAction{
		"manual":   FailureActionManual,
		"rollback": FailureActionRollback,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FailureAction(input)
	return &out, nil
}

type RollingUpgradeMode string

const (
	RollingUpgradeModeMonitored       RollingUpgradeMode = "Monitored"
	RollingUpgradeModeUnmonitoredAuto RollingUpgradeMode = "UnmonitoredAuto"
)

func PossibleValuesForRollingUpgradeMode() []string {
	return []string{
		string(RollingUpgradeModeMonitored),
		string(RollingUpgradeModeUnmonitoredAuto),
	}
}

func parseRollingUpgradeMode(input string) (*RollingUpgradeMode, error) {
	vals := map[string]RollingUpgradeMode{
		"monitored":       RollingUpgradeModeMonitored,
		"unmonitoredauto": RollingUpgradeModeUnmonitoredAuto,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RollingUpgradeMode(input)
	return &out, nil
}
