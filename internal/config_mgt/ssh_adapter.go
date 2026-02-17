package config_mgt

import (
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

type SSHAdapter struct{}

func NewSSHAdapter() *SSHAdapter {
	return &SSHAdapter{}
}

func (a *SSHAdapter) Execute(ip, user, password, command string) (string, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	client, err := ssh.Dial("tcp", ip+":22", config)
	if err != nil {
		return "", fmt.Errorf("failed to dial: %w", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	if err != nil {
		return string(output), fmt.Errorf("failed to run command: %w", err)
	}

	return string(output), nil
}
