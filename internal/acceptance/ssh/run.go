// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ssh

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"golang.org/x/crypto/ssh"
)

type Runner struct {
	Hostname      string
	Port          int
	Username      string
	Password      string
	CommandsToRun []string
}

func (r Runner) Run(ctx context.Context) error {
	if err := retry.RetryContext(ctx, 5*time.Minute, r.tryRun); err != nil {
		return err
	}

	return nil
}

func (r Runner) tryRun() *pluginsdk.RetryError {
	config := &ssh.ClientConfig{
		User: r.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(r.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	hostAddress := fmt.Sprintf("%s:%d", r.Hostname, r.Port)
	log.Printf("[INFO] SSHing to %q...", hostAddress)
	client, err := ssh.Dial("tcp", hostAddress, config)
	if err != nil {
		return pluginsdk.RetryableError(fmt.Errorf("connecting to host: %+v", err))
	}

	session, err := client.NewSession()
	if err != nil {
		return pluginsdk.RetryableError(fmt.Errorf("creating session: %+v", err))
	}
	defer session.Close()

	for _, cmd := range r.CommandsToRun {
		log.Printf("[DEBUG] Running %q..", cmd)
		var b bytes.Buffer
		session.Stdout = &b
		if err := session.Run(cmd); err != nil {
			return pluginsdk.NonRetryableError(fmt.Errorf("failure running command %q: %+v", cmd, err))
		}
	}

	return nil
}
