package pluginsdk

import (
	"context"
	"sort"
)

type StateUpgraderFunc = func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error)

type StateUpgrade interface {
	// Schema is a point-in-time reference to the Schema at the time of this version
	//
	// This only needs to include the Fields and Optional/Required/Computed/Defaults
	// it does not need to include validation functions so can/should be simplified
	//
	// NOTE: This also shouldn't reference the existing schema since it's a
	// point-in-time reference
	Schema() map[string]*Schema

	// UpgradeFunc is the upgrade function which should be run when upgrading this
	// from the "current" version to the "next" version of this Resource.
	UpgradeFunc() StateUpgraderFunc
}

func StateUpgrades(upgrades map[int]StateUpgrade) []StateUpgrader {
	versions := make([]int, len(upgrades))
	for version := range versions {
		versions = append(versions, version)
	}
	sort.Ints(versions)

	out := make([]StateUpgrader, 0)
	for _, version := range versions {
		upgrade := upgrades[version]
		resource := Resource{
			Schema: upgrade.Schema(),
		}
		// TODO: with Plugin SDK 1.x we'll need to add a wrapper here to inject ctx

		out = append(out, StateUpgrader{
			Type:    resource.CoreConfigSchema().ImpliedType(),
			Upgrade: upgrade.UpgradeFunc(),
			Version: version,
		})
	}
	return out
}
