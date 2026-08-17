// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"strings"
	"testing"
)

func TestSL012(t *testing.T) {
	fs := runRule(t, sl012, wrap(`
			"network_config":     {Type: pluginsdk.TypeString, Optional: true, Description: "d", ValidateFunc: validation.StringIsNotEmpty},
			"firewall_properties": {Type: pluginsdk.TypeString, Optional: true, Description: "d", ValidateFunc: validation.StringIsNotEmpty},
			"os_profile":         {Type: pluginsdk.TypeString, Optional: true, Description: "d", ValidateFunc: validation.StringIsNotEmpty},
			"name":               {Type: pluginsdk.TypeString, Optional: true, Description: "d", ValidateFunc: validation.StringIsNotEmpty},
	`))

	if f := at(fs, "network_config"); f == nil {
		t.Error("expected SL012 on network_config")
	} else if !strings.Contains(f.fix, `"network"`) {
		t.Errorf("expected rename to network, got %q", f.fix)
	}
	if !flagged(fs, "firewall_properties") || !flagged(fs, "os_profile") {
		t.Error("expected SL012 on _properties and _profile suffixes")
	}
	if flagged(fs, "name") {
		t.Error("did not expect SL012 on a name without a redundant suffix")
	}
}
