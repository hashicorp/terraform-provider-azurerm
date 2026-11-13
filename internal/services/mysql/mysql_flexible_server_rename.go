// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package mysql

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/checknameavailability"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func mysqlFlexibleServerRenameCustomizeDiff(_ context.Context, d *pluginsdk.ResourceDiff, _ any) error {
	if !d.HasChange("name") {
		return nil
	}
	if d.Id() != "" && !d.NewValueKnown("rename_enabled") {
		return fmt.Errorf("`rename_enabled` must be known before planning a MySQL Flexible Server name change")
	}
	if !d.Get("rename_enabled").(bool) {
		return d.ForceNew("name")
	}
	if d.Id() == "" {
		return nil
	}
	if !d.NewValueKnown("name") {
		return fmt.Errorf("`name` must be known before planning a MySQL Flexible Server rename")
	}
	old, new := d.GetChange("name")
	if strings.EqualFold(old.(string), new.(string)) {
		return fmt.Errorf("renaming a MySQL Flexible Server requires a different name, ignoring case")
	}
	if mysqlFlexibleServerHasOtherChanges(d) {
		return fmt.Errorf("a MySQL Flexible Server rename must be applied separately from changes to other server properties")
	}
	if err := d.SetNewComputed("id"); err != nil {
		return err
	}
	return d.SetNewComputed("fqdn")
}

func mysqlFlexibleServerHasOtherChanges(d interface{ HasChange(string) bool }) bool {
	for name, field := range resourceMysqlFlexibleServer().Schema {
		if name == "name" || name == "rename_enabled" || field.WriteOnly || (!field.Optional && !field.Required) {
			continue
		}
		if d.HasChange(name) {
			return true
		}
	}
	return false
}

func resourceMysqlFlexibleServerRename(ctx context.Context, d *pluginsdk.ResourceData, meta *clients.Client, sourceID servers.FlexibleServerId) (err error) {
	// Until completion is verified, the source remains the only identity we can safely retain.
	defer func() {
		if err != nil && d.Id() == sourceID.ID() {
			if setErr := d.Set("name", sourceID.FlexibleServerName); setErr != nil {
				err = fmt.Errorf("%w; restoring source name in state: %+v", err, setErr)
			}
			oldFQDN, _ := d.GetChange("fqdn")
			if setErr := d.Set("fqdn", oldFQDN); setErr != nil {
				err = fmt.Errorf("%w; restoring source FQDN in state: %+v", err, setErr)
			}
		}
	}()

	if !d.Get("rename_enabled").(bool) {
		return fmt.Errorf("`rename_enabled` must be true to rename a MySQL Flexible Server")
	}
	if mysqlFlexibleServerHasOtherChanges(d) {
		return fmt.Errorf("a MySQL Flexible Server rename must be applied separately from changes to other server properties")
	}
	destinationID := servers.NewFlexibleServerID(sourceID.SubscriptionId, sourceID.ResourceGroupName, d.Get("name").(string))
	if strings.EqualFold(sourceID.FlexibleServerName, destinationID.FlexibleServerName) {
		return fmt.Errorf("renaming a MySQL Flexible Server requires a different name, ignoring case")
	}

	client := meta.MySQL.FlexibleServers.Servers
	source, err := client.Get(ctx, sourceID)
	if err != nil {
		return fmt.Errorf("retrieving rename source %s: %+v", sourceID, err)
	}
	if source.Model == nil || source.Model.Id == nil || source.Model.Location == "" || source.Model.Properties == nil {
		return fmt.Errorf("retrieving rename source %s: response is missing the server ID, location or properties", sourceID)
	}
	canonicalID, err := servers.ParseFlexibleServerIDInsensitively(*source.Model.Id)
	if err != nil || !strings.EqualFold(canonicalID.ID(), sourceID.ID()) {
		return fmt.Errorf("retrieving rename source %s: response contains an unexpected server ID %q", sourceID, *source.Model.Id)
	}
	if pointer.From(source.Model.Properties.State) != servers.ServerStateReady {
		return fmt.Errorf("renaming %s: the source server must be Ready", sourceID)
	}

	existing, err := client.Get(ctx, destinationID)
	if err == nil {
		return fmt.Errorf("renaming %s: destination %s already exists; it will not be overwritten or adopted", sourceID, destinationID)
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return fmt.Errorf("checking rename destination %s: %+v", destinationID, err)
	}

	availabilityClient := meta.MySQL.FlexibleServers.CheckNameAvailability
	availabilityRequest := checknameavailability.NameAvailabilityRequest{
		Name: destinationID.FlexibleServerName,
		Type: pointer.To("Microsoft.DBforMySQL/flexibleServers"),
	}
	availability, err := availabilityClient.Execute(ctx,
		checknameavailability.NewLocationID(sourceID.SubscriptionId, source.Model.Location), availabilityRequest)
	if err != nil && response.WasStatusCode(availability.HttpResponse, http.StatusForbidden) &&
		availability.OData != nil && availability.OData.Error != nil &&
		pointer.From(availability.OData.Error.Code) == "AuthorizationFailed" {
		fallback, fallbackErr := availabilityClient.WithoutLocationExecute(ctx, commonids.NewSubscriptionID(sourceID.SubscriptionId), availabilityRequest)
		availability.Model, err = fallback.Model, fallbackErr
	}
	if err != nil {
		return fmt.Errorf("checking name availability for %s: %+v", destinationID, err)
	}
	if availability.Model == nil || availability.Model.NameAvailable == nil {
		return fmt.Errorf("checking name availability for %s: response is missing nameAvailable", destinationID)
	}
	if !*availability.Model.NameAvailable {
		return fmt.Errorf("name %q is unavailable: %s", destinationID.FlexibleServerName, pointer.From(availability.Model.Message))
	}

	result, err := meta.MySQL.FlexibleServerRename.Create(ctx, destinationID, servers.Server{
		Location: source.Model.Location,
		Properties: &servers.ServerProperties{
			CreateMode:             pointer.To(servers.CreateMode("Rename")),
			SourceServerResourceId: source.Model.Id,
		},
	})
	operationURL := ""
	if result.HttpResponse != nil {
		operationURL = result.HttpResponse.Header.Get("Azure-AsyncOperation")
		if operationURL == "" {
			operationURL = result.HttpResponse.Header.Get("Location")
		}
	}
	recoveryError := func(cause error) error {
		return fmt.Errorf("renaming %s to %s (operation URL %q): %+v. The source identity has been retained; inspect the operation and both server IDs before retrying. If rename completed, reconcile the destination with Terraform state instead of submitting another rename", sourceID, destinationID, operationURL, cause)
	}
	if err != nil {
		return recoveryError(err)
	}
	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return recoveryError(err)
	}

	var destination *servers.Server
	poller := custompollers.NewEventualConsistencyPoller(1, func(pollerCtx context.Context) (*http.Response, error) {
		resp, err := client.Get(pollerCtx, destinationID)
		destination = resp.Model
		return resp.HttpResponse, err
	}, custompollers.DefaultCreationEventualConsistencyPollerOptions())
	if err := poller.PollUntilDone(ctx); err != nil {
		return recoveryError(err)
	}
	if destination == nil || destination.Id == nil || !strings.EqualFold(*destination.Id, destinationID.ID()) ||
		destination.Properties == nil || pointer.From(destination.Properties.State) != servers.ServerStateReady ||
		pointer.From(destination.Properties.FullyQualifiedDomainName) == "" {
		return recoveryError(fmt.Errorf("destination did not return the expected ID, Ready state and FQDN"))
	}

	d.SetId(destinationID.ID())
	if err := d.Set("name", destinationID.FlexibleServerName); err != nil {
		return err
	}
	if err := pluginsdk.SetResourceIdentityData(d, &destinationID); err != nil {
		return err
	}
	if err := resourceMysqlFlexibleServerFlatten(d, &destinationID, destination); err != nil {
		return err
	}
	return nil
}
