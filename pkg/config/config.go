package config

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
)

const (
	// WorkDir is the directory used by thetool
	WorkDir = "repositories"
	// EnvoyHash is the commit hash of the version of Envoy used
	EnvoyHash = "29989a38c017d3be5aa3c735a797fcf58b754fe5"
	// EnvoyBuilderHash
	EnvoyBuilderHash = "52f6880ffbf761c9b809fc3ac208900956ff16b4"
	// GlooHash is the commit hash of the version of Gloo used
	GlooHash = "c88c90c332e5528a070a1c800bc65b2c39f8ca24"
	// GlooRepo is the repository URL for Gloo
	GlooRepo = "https://github.com/solo-io/gloo.git"
	//GlooChartHash is the commit hash of the Gloo chart used
	GlooChartHash = "a2f12f82fd41d7b7eab91a2f825fee8f9fdf6ec5"
	//GlooChartRepo is the repository URL for Gloo chart
	GlooChartRepo = "https://github.com/solo-io/gloo-chart.git"

	// DockerUser is the default Docker registry user used for publishing the images
	DockerUser = "soloio"
	// ConfigFile is the name of the configuraiton file
	ConfigFile = "thetool.json"
)

// Config contains the configuration used by thetool
type Config struct {
	WorkDir          string `json:"workDir"`
	EnvoyHash        string `json:"envoyHash"`
	EnvoyBuilderHash string `json:"envoyBuilderHash"`
	GlooHash         string `json:"glooHash"`
	GlooRepo         string `json:"glooRepo"`
	GlooChartHash    string `json:"glooChartHash"`
	GlooChartRepo    string `json:"glooChartRepo"`
	DockerUser       string `json:"dockerUser,omitempty"`
}

// Save the current configuration used by thetool to a file
func (c *Config) Save(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return saveToWriter(c, f)
}

func saveToWriter(c *Config, w io.Writer) error {
	b, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

// Load the configuration for thetool from a file
func Load(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return loadFromReader(f)
}

func loadFromReader(r io.Reader) (*Config, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	err = json.Unmarshal(buf, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}