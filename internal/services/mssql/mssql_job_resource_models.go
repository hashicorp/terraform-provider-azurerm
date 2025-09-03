package mssql

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MsSqlJobResourceModel struct {
	ID       types.String   `tfsdk:"id"`
	Timeouts timeouts.Value `tfsdk:"timeouts"`

	Name        types.String `tfsdk:"name"`
	JobAgentID  types.String `tfsdk:"job_agent_id"`
	Description types.String `tfsdk:"description"`
}

type MsSqlJobResourceIdentityModel struct {
	SubscriptionId    string `tfsdk:"subscription_id"`
	ResourceGroupName string `tfsdk:"resource_group_name"`
	ServerName        string `tfsdk:"server_name"`
	JobAgentName      string `tfsdk:"job_agent_name"`
	Name              string `tfsdk:"name"`
}
