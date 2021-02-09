package ssh

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"golang.org/x/crypto/ssh"
)

type Runner struct {
	Hostname      string
	Port          int
	Username      string
	Password      string
	CommandsToRun []string
}

func (r Runner) Run() error {
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"NotFound"},
		Target:                    []string{"Succeeded"},
		Refresh:                   r.tryRun,
		MinTimeout:                1 * time.Minute,
		PollInterval:              15 * time.Second,
		ContinuousTargetOccurence: 1,
		Timeout:                   5 * time.Minute,
	}
	result, err := stateConf.WaitForState()
	if err != nil {
		return err
	}

	v, ok := result.(bool)
	if !ok || !v {
		return fmt.Errorf("failure connecting to host/running commands")
	}

	return nil
}

func (r Runner) tryRun() (result interface{}, state string, err error) {
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
		return nil, "NotFound", fmt.Errorf("connecting to host: %+v", err)
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, "", fmt.Errorf("creating session: %+v", err)
	}
	defer session.Close()

	for _, cmd := range r.CommandsToRun {
		log.Printf("[DEBUG] Running %q..", cmd)
		var b bytes.Buffer
		session.Stdout = &b
		if err := session.Run(cmd); err != nil {
			return false, "Failed", fmt.Errorf("failure running command %q: %+v", cmd, err)
		}
	}

	return true, "Succeeded", nil
}
