package helm

import (
	"encoding/base64"
	"fmt"
	"os/exec"
	"strings"

	"code.cloudfoundry.org/lager"
)

const (
	instanceIDLogKey      = "instance-id"
	programLogKey         = "program"
	argumentsLogKey       = "arguments"
	execCmdResponseLogKey = "exec-cmd-response"
)

var (
	escaper   = strings.NewReplacer("9", "99", "-", "90", "_", "91")
	unescaper = strings.NewReplacer("99", "9", "90", "-", "91", "_")
)

type Client struct {
	config Config
	logger lager.Logger
}

func New(config Config, logger lager.Logger) *Client {
	return &Client{
		config: config,
		logger: logger.Session("helm"),
	}
}

func (c *Client) Install(instanceID string) error {
	c.logger.Debug("install", lager.Data{
		instanceIDLogKey: instanceID,
	})

	// TODO

	return nil
}

func (c *Client) Update(instanceID string) error {
	c.logger.Debug("install", lager.Data{
		instanceIDLogKey: instanceID,
	})

	// TODO

	return nil
}

func (c *Client) Delete(instanceID string) error {
	c.logger.Debug("delete", lager.Data{
		instanceIDLogKey: instanceID,
	})

	// TODO

	return nil
}

func (c *Client) releaseName(instanceID string) string {
	return fmt.Sprintf("%s-%s", c.config.ReleaseNamePrefix, escaper.Replace(base64.RawURLEncoding.EncodeToString([]byte(strings.Replace(instanceID, "-", "", -1)))))
}

func (c *Client) helm(cmd string) (string, error) {
	args := []string{}
	if c.config.TillerHost != "" {
		args = append(args, fmt.Sprintf("--host %s", c.config.TillerHost))
	}
	if c.config.TillerNamespace != "" {
		args = append(args, fmt.Sprintf("--tiller-namespace %s", c.config.TillerNamespace))
	}
	if c.config.KubeContext != "" {
		args = append(args, fmt.Sprintf("--kube-context %s", c.config.KubeContext))
	}
	if c.config.Debug {
		args = append(args, "--debug")
	}

	args = append(args, strings.Fields(cmd)...)

	c.logger.Debug("helm", lager.Data{
		"program":   c.config.BinaryLocation,
		"arguments": args,
	})

	out, err := exec.Command(c.config.BinaryLocation, args...).CombinedOutput()
	if err != nil {
		c.logger.Error("helm", err)
		c.logger.Debug("helm", lager.Data{
			"output": string(out),
		})
		return "", err
	}

	c.logger.Debug("helm", lager.Data{
		"output": string(out),
	})

	return string(out), nil
}
