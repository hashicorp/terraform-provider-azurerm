package ssh

import (
	"bytes"
	"fmt"
	"log"

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
		return fmt.Errorf("connecting to host: %+v", err)
	}

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("creating session: %+v", err)
	}
	defer session.Close()

	for _, cmd := range r.CommandsToRun {
		log.Printf("[DEBUG] Running %q..", cmd)
		var b bytes.Buffer
		session.Stdout = &b
		if err := session.Run(cmd); err != nil {
			return fmt.Errorf("failure running command %q: %+v", cmd, err)
		}
	}

	return nil
}
