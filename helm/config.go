package helm

import (
	"errors"
)

type Config struct {
	ReleaseNamePrefix string `json:"release_name_prefix"`
	DefaultNamespace  string `json:"default_namespace"`
	BinaryLocation    string `json:"binary_location"`
	TillerHost        string `json:"tiller_host,omitempty"`
	TillerNamespace   string `json:"tiller_namespace,omitempty"`
	KubeContext       string `json:"kube_context,omitempty"`
	Debug             bool   `json:"debug"`
}

func (c Config) Validate() error {
	if c.ReleaseNamePrefix == "" {
		return errors.New("Must provide a non-empty Release Name Prefix")
	}

	if c.DefaultNamespace == "" {
		return errors.New("Must provide a non-empty Default Namespace")
	}

	if c.BinaryLocation == "" {
		return errors.New("Must provide a non-empty Binary Location")
	}

	return nil
}
