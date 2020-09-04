package tfexec

import (
	"bytes"
	"context"
	"encoding/json"
	"os/exec"

	tfjson "github.com/hashicorp/terraform-json"
)

func (tf *Terraform) ProvidersSchema(ctx context.Context) (*tfjson.ProviderSchemas, error) {
	schemaCmd := tf.providersSchemaCmd(ctx)

	var ret tfjson.ProviderSchemas
	var outBuf bytes.Buffer
	schemaCmd.Stdout = &outBuf

	err := tf.runTerraformCmd(schemaCmd)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(outBuf.Bytes(), &ret)
	if err != nil {
		return nil, err
	}

	err = ret.Validate()
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func (tf *Terraform) providersSchemaCmd(ctx context.Context, args ...string) *exec.Cmd {
	allArgs := []string{"providers", "schema", "-json", "-no-color"}
	allArgs = append(allArgs, args...)

	return tf.buildTerraformCmd(ctx, allArgs...)
}
