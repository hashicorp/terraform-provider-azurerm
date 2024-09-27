// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/tombuildsstuff/kermit/sdk/datafactory/2018-06-01/datafactory" // nolint: staticcheck
)

func importDataFactoryLinkedService(expectType datafactory.TypeBasicLinkedService) pluginsdk.ImporterFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) (data []*pluginsdk.ResourceData, err error) {
		id, err := parse.LinkedServiceID(d.Id())
		if err != nil {
			return nil, err
		}

		client := meta.(*clients.Client).DataFactory.LinkedServiceClient
		resp, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
		if err != nil {
			return nil, fmt.Errorf("retrieving Data Factory %s: %+v", *id, err)
		}

		byteArr, err := json.Marshal(resp.Properties)
		if err != nil {
			return nil, err
		}

		var m map[string]*json.RawMessage
		if err = json.Unmarshal(byteArr, &m); err != nil {
			return nil, err
		}

		t := ""
		if v, ok := m["type"]; ok && v != nil {
			if err := json.Unmarshal(*v, &t); err != nil {
				return nil, err
			}
			delete(m, "type")
		}

		if datafactory.TypeBasicLinkedService(t) != expectType {
			return nil, fmt.Errorf("data factory linked service has mismatched type, expected: %q, got %q", expectType, t)
		}

		return []*pluginsdk.ResourceData{d}, nil
	}
}
