// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package logic

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func TestMergeAppSettingsRemovesAllDeletedValues(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, LogicAppResource{}.Arguments(), map[string]any{})

	existing := []webapps.NameValuePair{
		{Name: pointer.To("preserve"), Value: pointer.To("existing")},
		{Name: pointer.To("keep"), Value: pointer.To("existing")},
		{Name: pointer.To("remove-one"), Value: pointer.To("one")},
		{Name: pointer.To("remove-two"), Value: pointer.To("two")},
		{Name: pointer.To("update"), Value: pointer.To("old")},
	}

	old := map[string]any{
		"keep":       "existing",
		"remove-one": "one",
		"remove-two": "two",
		"update":     "old",
	}
	new := map[string]any{
		"keep":   "configured",
		"update": "new",
		"add":    "value",
	}

	result := mergeAppSettings(existing, old, new, sdk.ResourceMetaData{ResourceData: resourceData})

	got := make(map[string]string, len(*result))
	for _, setting := range *result {
		got[pointer.From(setting.Name)] = pointer.From(setting.Value)
	}

	want := map[string]string{
		"preserve":                           "existing",
		"keep":                               "configured",
		"update":                             "new",
		"add":                                "value",
		functionVersionAppSettingName:        "~4",
		extensionBundleAppSettingName:        extensionBundleAppSettingValue,
		extensionBundleVersionAppSettingName: "[1.*, 2.0.0)",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected merged app settings: got %#v, want %#v", got, want)
	}
}
