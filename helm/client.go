package helm

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"code.cloudfoundry.org/lager"
)

const (
	instanceIDLogKey = "instance-id"
	chartLogKey      = "chart"
	repositoryLogKey = "repository"
	versionLogKey    = "version"
	programLogKey    = "program"
	argumentsLogKey  = "arguments"
	outputLogKey     = "output"
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

func (c *Client) InstallRelease(instanceID string, chart string, repository string, version string) error {
	c.logger.Debug("install-release-parameters", lager.Data{
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

func (c *Client) UpgradeRelease(instanceID string, chart string, repository string, version string) error {
	c.logger.Debug("upgrade-release-parameters", lager.Data{
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

func (c *Client) DeleteRelease(instanceID string) error {
	c.logger.Debug("delete-release-parameters", lager.Data{
		instanceIDLogKey: instanceID,
	})

	cmd := fmt.Sprintf("delete --purge %s", c.releaseName(instanceID))
	if _, err := c.helm(cmd); err != nil {
		return fmt.Errorf("Error deleting Helm release `%s`", c.releaseName(instanceID))
	}

	return nil
}

func (c *Client) ReleaseStatus(instanceID string) (string, string, error) {
	c.logger.Debug("release-status-parameters", lager.Data{
		instanceIDLogKey: instanceID,
	})

	status := "FAILED"
	description := ""

	cmd := fmt.Sprintf("status %s", c.releaseName(instanceID))
	out, err := c.helm(cmd)
	if err != nil {
		return status, description, fmt.Errorf("Error getting status for Helm release `%s`", c.releaseName(instanceID))
	}

	statusRe := regexp.MustCompile(`\nSTATUS: ([A-Z]+)\n`)
	capturedStatus := statusRe.FindStringSubmatch(out)
	if capturedStatus != nil {
		switch capturedStatus[1] {
		case "PENDING_INSTALL":
			status = "INPROGRESS"
		case "PENDING_UPGRADE":
			status = "INPROGRESS"
		case "DEPLOYED":
			status = "SUCCEEDED"
		case "DELETING":
			status = "INPROGRESS"
		case "DELETED":
			status = "SUCCEEDED"
		}
	}

	lastDeployedRe := regexp.MustCompile(`LAST DEPLOYED: (.+)\n`)
	capturedlastDeployed := lastDeployedRe.FindStringSubmatch(out)
	if capturedlastDeployed != nil {
		description = fmt.Sprintf("Last deployed: %s", capturedlastDeployed[1])
	}

	return status, description, nil
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

	c.logger.Debug("exec", lager.Data{
		programLogKey:   c.config.BinaryLocation,
		argumentsLogKey: args,
	})

	out, err := exec.Command(c.config.BinaryLocation, args...).CombinedOutput()
	if err != nil {
		c.logger.Error("exec", err)
		c.logger.Debug("exec", lager.Data{
			outputLogKey: string(out),
		})
		return "", err
	}

	c.logger.Debug("exec", lager.Data{
		outputLogKey: string(out),
	})

	return string(out), nil
}
