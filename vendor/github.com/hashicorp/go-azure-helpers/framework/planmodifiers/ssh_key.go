package planmodifiers

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

type sshKeyPlanModifier struct{}

func SSHKey() sshKeyPlanModifier {
	return sshKeyPlanModifier{}
}

func (s sshKeyPlanModifier) Description(_ context.Context) string {
	return "this plan modifier normalises the planned and stored values for SSH keys and suppresses the diff if the normalised values are equal, otherwise triggers RequiresReplace"
}

func (s sshKeyPlanModifier) MarkdownDescription(ctx context.Context) string {
	return s.Description(ctx)
}

func (s sshKeyPlanModifier) PlanModifyString(_ context.Context, request planmodifier.StringRequest, response *planmodifier.StringResponse) {
	planKey, err := NormalizeSSHKey(request.PlanValue.ValueString())
	if err != nil {
		response.Diagnostics.AddError("normalising planned SSH key failed", "Failed to normalise planned SSH key.")
		return
	}

	if request.StateValue.IsNull() || request.StateValue.IsUnknown() {
		return
	}

	stateKey, err := NormalizeSSHKey(request.StateValue.ValueString())
	if err != nil {
		response.Diagnostics.AddError("normalising state SSH key failed", "Failed to normalise state key.")
		return
	}

	if strings.EqualFold(*planKey, *stateKey) {
		response.PlanValue = request.StateValue
		response.RequiresReplace = false
		return
	}

	response.RequiresReplace = true
}

var _ planmodifier.String = &sshKeyPlanModifier{}

// NormalizeSSHKey attempts to rationalise the SSH Key data to account for some of the interesting implementations in Azure...
func NormalizeSSHKey(input string) (*string, error) {
	if input == "" {
		return nil, fmt.Errorf("empty string supplied")
	}

	output := input
	output = strings.ReplaceAll(output, "<<~EOT", "")
	output = strings.ReplaceAll(output, "EOT", "")
	output = strings.ReplaceAll(output, "\r", "")

	lines := make([]string, 0)
	for _, line := range strings.Split(output, "\n") {
		lines = append(lines, strings.TrimSpace(line))
	}

	normalised := strings.Join(lines, "")

	return pointer.To(normalised), nil
}
