package clients

import (
	"context"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/analysisservices"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/applicationinsights"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/authorization"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/automation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/batch"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/bot"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cdn"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cognitive"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/containers"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/databricks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datafactory"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datalake"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/devspace"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/devtestlabs"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/dns"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventgrid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/frontdoor"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/graph"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/hdinsight"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/healthcare"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iothub"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/kusto"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/logic"
)

type Client struct {
	// StopContext is used for propagating control from Terraform Core (e.g. Ctrl/Cmd+C)
	StopContext context.Context

	AnalysisServices *analysisservices.Client
	ApiManagement    *apimanagement.Client
	AppInsights      *applicationinsights.Client
	Automation       *automation.Client
	Authorization    *authorization.Client
	Batch            *batch.Client
	Bot              *bot.Client
	Cdn              *cdn.Client
	Cognitive        *cognitive.Client
	Containers       *containers.Client
	Cosmos           *cosmos.Client
	Compute          *ComputeClient
	DataBricks       *databricks.Client
	DataFactory      *datafactory.Client
	Datalake         *datalake.Client
	DevSpace         *devspace.Client
	DevTestLabs      *devtestlabs.Client
	Dns              *dns.Client
	EventGrid        *eventgrid.Client
	Eventhub         *eventhub.Client
	Frontdoor        *frontdoor.Client
	Graph            *graph.Client
	HDInsight        *hdinsight.Client
	Healthcare       *healthcare.Client
	IoTHub           *iothub.Client
	KeyVault         *keyvault.Client
	Kusto            *kusto.Client
	LogAnalytics     *loganalytics.Client
	Logic            *logic.Client
}
