// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package timeouts

import (
	"sync"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// dataSourceReadTimeouts tracks the effective read timeout for each in-flight Data Source
// read, keyed by its *pluginsdk.ResourceData.
//
// This exists to work around https://github.com/hashicorp/terraform-plugin-sdk/issues/1038:
// the Plugin SDK never attaches the schema-declared Timeouts (or the user's `timeouts` block)
// to the ResourceData for a Data Source read, so `d.Timeout(TimeoutRead)` always returns the
// SDK-global default of 20 minutes. Data Source read functions use the legacy signature without
// a context argument, so the effective timeout cannot be threaded through a context either -
// instead the provider registers it here before invoking the read function (see
// internal/provider) and ForRead prefers it over `d.Timeout`.
var dataSourceReadTimeouts sync.Map

// RegisterDataSourceReadTimeout records the effective read timeout for an in-flight Data
// Source read. Callers must call DeregisterDataSourceReadTimeout once the read completes.
func RegisterDataSourceReadTimeout(d *pluginsdk.ResourceData, timeout time.Duration) {
	dataSourceReadTimeouts.Store(d, timeout)
}

// DeregisterDataSourceReadTimeout removes the effective read timeout recorded for an
// in-flight Data Source read.
func DeregisterDataSourceReadTimeout(d *pluginsdk.ResourceData) {
	dataSourceReadTimeouts.Delete(d)
}

func dataSourceReadTimeout(d *pluginsdk.ResourceData) (time.Duration, bool) {
	if v, ok := dataSourceReadTimeouts.Load(d); ok {
		return v.(time.Duration), true
	}
	return 0, false
}
