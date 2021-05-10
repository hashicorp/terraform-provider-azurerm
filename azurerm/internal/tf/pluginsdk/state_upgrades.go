package pluginsdk

import (
	"context"
	"fmt"
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

// StateUpgrades is a wrapper around the Plugin SDK's State Upgraders
// which allows us to upgrade the Plugin SDK without breaking all open
// PR's and attempts to make this interface a little less verbose.
func StateUpgrades(upgrades map[int]StateUpgrade) []StateUpgrader {
	versions := make([]int, 0)
	for version := range upgrades {
		versions = append(versions, version)
	}
	sort.Ints(versions)

	out := make([]StateUpgrader, 0)
	expectedVersion := 0
	for _, version := range versions {
		// sanity check that there's no missing versions
		if expectedVersion != version {
			panic(fmt.Sprintf("missing state upgrade for version %d", expectedVersion))
		}
		expectedVersion++

		upgrade := upgrades[version]
		resource := Resource{
			Schema: upgrade.Schema(),
		}
		// TODO: with Plugin SDK 1.x we'll need to add a wrapper here to inject ctx

		out = append(out, StateUpgrader{
			Type: resource.CoreConfigSchema().ImpliedType(),
			Upgrade: func(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
				return upgrade.UpgradeFunc()(context.TODO(), rawState, meta)
			},
			Version: version,
		})
	}
	return out
}
