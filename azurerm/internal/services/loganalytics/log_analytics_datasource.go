package loganalytics

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/operationalinsights/mgmt/2015-11-01-preview/operationalinsights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

type collectionStateIndicator struct {
	enabled  string
	disabled string
}

// collectionStateIndicatorLookup is a table for collection data sources to lookup the enabled/disabled
// indicator.
var collectionStateIndicatorLookup = map[operationalinsights.DataSourceKind]collectionStateIndicator{
	operationalinsights.LinuxSyslogCollection: {
		"Enabled",
		"Disabled",
	},
}

type dataSourceCollectionProperty struct {
	State string `json:"state"`
}

func filterForDataSourceKind(kind operationalinsights.DataSourceKind) string {
	return fmt.Sprintf("kind eq '%s'", kind)
}

func dataSourceCollectionName(kind operationalinsights.DataSourceKind) string {
	// This is the default name for collection resource when creating via Portal
	return fmt.Sprintf("DataSource_%s", kind)
}

// createDataSourceCollection guarantees to create a certain kind of data source collections with expected state (enabled/disabled).
func createDataSourceCollection(ctx context.Context, client *operationalinsights.DataSourcesClient, resourceGroupName, workspaceName string, kind operationalinsights.DataSourceKind, toEnable bool) error {
	// Check whether collection already being setup.
	ds, err := getDataSourceCollection(ctx, client, kind, resourceGroupName, workspaceName)
	if err != nil {
		return err
	}
	name := dataSourceCollectionName(kind)
	// If there is already a data source collection, then update the state of that collection.
	if ds != nil {
		if ds.Name == nil {
			return fmt.Errorf("unexpected nil name of the existed Log Analytics DataSource %s (Resource Group %q / Workspace: %q)", kind, resourceGroupName, workspaceName)
		}
		name = *ds.Name
	}

	state := collectionStateIndicatorLookup[kind].disabled
	if toEnable {
		state = collectionStateIndicatorLookup[kind].enabled
	}
	param := operationalinsights.DataSource{
		Kind: kind,
		Properties: &dataSourceCollectionProperty{
			State: state,
		},
	}
	if _, err := client.CreateOrUpdate(ctx, resourceGroupName, workspaceName, name, param); err != nil {
		return fmt.Errorf("failed to create Log Analytics DataSource %s: %s (Resource Group %q / Workspace: %q): %+v", kind, name, resourceGroupName, workspaceName, err)
	}
	return nil
}

// getDataSourceCollectionState gets current state of a certain kind of data source collection. In case of absent of collection data source, we regards it as disabled.
func getDataSourceCollectionState(ctx context.Context, client *operationalinsights.DataSourcesClient, resourceGroupName, workspaceName string, kind operationalinsights.DataSourceKind) (bool, error) {
	ds, err := getDataSourceCollection(ctx, client, kind, resourceGroupName, workspaceName)
	if err != nil {
		return false, err
	}
	if ds == nil {
		return false, nil
	}
	propStr, err := structure.FlattenJsonToString(ds.Properties.(map[string]interface{}))
	if err != nil {
		return false, fmt.Errorf("failed to flatten properties map to json for Log Analytics DataSource %s (Resource Group %q / Workspace: %q): %+v", kind, resourceGroupName, workspaceName, err)
	}
	prop := dataSourceCollectionProperty{}
	if err := json.Unmarshal([]byte(propStr), &prop); err != nil {
		return false, fmt.Errorf("failed to decode properties json for Log Analytics DataSource %s (Resource Group %q / Workspace: %q): %+v", kind, resourceGroupName, workspaceName, err)
	}
	if prop.State == collectionStateIndicatorLookup[kind].enabled {
		return true, nil
	}
	if prop.State == collectionStateIndicatorLookup[kind].disabled {
		return false, nil
	}
	return false, fmt.Errorf("unknown state of Log Analytics DataSource %s (Resource Group %q / Workspace: %q): %s", kind, resourceGroupName, workspaceName, prop.State)
}

func getDataSourceCollection(ctx context.Context, client *operationalinsights.DataSourcesClient, kind operationalinsights.DataSourceKind, resourceGroupName, workspaceName string) (ds *operationalinsights.DataSource, err error) {
	// As there is at most one collection defiend for each kind among the workspace, so we will only iterate for one time.
	iterator, err := client.ListByWorkspace(ctx, resourceGroupName, workspaceName, filterForDataSourceKind(kind), "")
	if err != nil {
		return nil, fmt.Errorf("failed to list Log Analytics DataSource for %s (Resource Group: %q / Workspace: %q): %+v", kind, resourceGroupName, workspaceName, err)
	}
	collections := iterator.Values()
	if len(collections) == 0 {
		return nil, nil
	}
	return &collections[0], nil
}

// importLogAnalyticsDataSource returns a StateFunc to be used in
// `ValidateresourceIDPriorToImportThen` so as to guarantee the data source kind belogns
// to the current resource.
func importLogAnalyticsDataSource(kind operationalinsights.DataSourceKind) func(d *schema.ResourceData, meta interface{}) (data []*schema.ResourceData, err error) {
	return func(d *schema.ResourceData, meta interface{}) (data []*schema.ResourceData, err error) {
		id, err := parse.LogAnalyticsDataSourceID(d.Id())
		if err != nil {
			return nil, err
		}

		client := meta.(*clients.Client).LogAnalytics.DataSourcesClient
		ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
		defer cancel()

		resp, err := client.Get(ctx, id.ResourceGroup, id.Workspace, id.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve Log Analytics Data Source %q (Resource Group %q / Workspace: %q): %+v", id.Name, id.ResourceGroup, id.Workspace, err)
		}

		if resp.Kind != kind {
			return nil, fmt.Errorf(`Log Analytics Data Source "kind" mismatch, expected "%s", got "%s"`, kind, resp.Kind)
		}
		return []*schema.ResourceData{d}, nil
	}
}
