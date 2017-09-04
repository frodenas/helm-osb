package helm

import (
	"fmt"
	"os/exec"
	"strings"

	"code.cloudfoundry.org/lager"
)

const (
	instanceIDLogKey = "instance-id"
	chartLogKey      = "chart"
	repositoryLogKey = "repository"
	versionLogKey    = "version"
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

func (c *Client) Install(instanceID string, chart string, repository string, version string) error {
	c.logger.Debug("install", lager.Data{
		instanceIDLogKey: instanceID,
		chartLogKey:      chart,
		repositoryLogKey: repository,
		versionLogKey:    version,
	})

	cmd := fmt.Sprintf("install %s --name %s --namespace %s", chart, c.releaseName(instanceID), c.config.DefaultNamespace)
	if repository != "" {
		cmd = cmd + fmt.Sprintf(" --repo %s", repository)
	}
	if version != "" {
		cmd = cmd + fmt.Sprintf(" --version %s", version)
	}
	if _, err := c.helm(cmd); err != nil {
		return fmt.Errorf("Error installing Helm release `%s`", c.releaseName(instanceID))
	}

	return nil
}

func (c *Client) Upgrade(instanceID string, chart string, repository string, version string) error {
	c.logger.Debug("upgrade", lager.Data{
		instanceIDLogKey: instanceID,
		chartLogKey:      chart,
		repositoryLogKey: repository,
		versionLogKey:    version,
	})

	cmd := fmt.Sprintf("upgrade %s %s --namespace %s", c.releaseName(instanceID), chart, c.config.DefaultNamespace)
	if repository != "" {
		cmd = cmd + fmt.Sprintf(" --repo %s", repository)
	}
	if version != "" {
		cmd = cmd + fmt.Sprintf(" --version %s", version)
	}
	if _, err := c.helm(cmd); err != nil {
		return fmt.Errorf("Error upgrading Helm release `%s`", c.releaseName(instanceID))
	}

	return nil
}

func (c *Client) Delete(instanceID string) error {
	c.logger.Debug("delete", lager.Data{
		instanceIDLogKey: instanceID,
	})

	cmd := fmt.Sprintf("delete --purge %s", c.releaseName(instanceID))
	if _, err := c.helm(cmd); err != nil {
		return fmt.Errorf("Error deleting Helm release `%s`", c.releaseName(instanceID))
	}

	return nil
}

func (c *Client) releaseName(instanceID string) string {
	return fmt.Sprintf("%s-%s", c.config.ReleaseNamePrefix, strings.Replace(instanceID, "-", "", -1))
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
	if c.config.Home != "" {
		args = append(args, fmt.Sprintf("--home %s", c.config.Home))
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
