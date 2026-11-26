// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package analysis_services_server

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/analysisservices/2017-08-01/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAnalysisServicesServerRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AnalysisServices.Servers
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := servers.ParseServerID(d.Id())
	if err != nil {
		return err
	}

	server, err := client.GetDetails(ctx, *id)
	if err != nil {
		if response.WasNotFound(server.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return resourceAnalysisServicesServerFlatten(d, id, server.Model)
}

func hashAnalysisServicesServerIPv4FirewallRule(v any) int {
	var buf bytes.Buffer
	m := v.(map[string]any)

	fmt.Fprintf(&buf, "%s-", strings.ToLower(m["name"].(string)))
	fmt.Fprintf(&buf, "%s-", m["range_start"].(string))
	buf.WriteString(m["range_end"].(string))

	return pluginsdk.HashString(buf.String())
}
