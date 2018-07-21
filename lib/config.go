package lib

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// Conf is the global config
var Conf *Config

// The Config says how certain things should behave, and is created from the
// config file.
type Config struct {
	MapWidth                   int     `yaml:"map-width"`
	MapHeight                  int     `yaml:"map-height"`
	GridSpacing                int     `yaml:"grid-spacing"`
	RoadWidth                  int     `yaml:"road-width"`
	RoomWidth                  int     `yaml:"room-radius"`
	RoomWidthVariance          float64 `yaml:"room-radius-variance"`
	NodeChance                 float64 `yaml:"node-chance"`
	RoomProbabilityCoefficient float64 `yaml:"room-prob-coefficient"`
	NumThreads                 int     `yaml:"num-threads"`
	BoxChance                  float64 `yaml:"box-chance"`
	ChestChance                float64 `yaml:"chest-chance"`
	NumMerchants               int     `yaml:"num-merchants"`
}

// LoadConfig creates a new Config instance from the given file.
func LoadConfig(filename string) (*Config, error) {
	cfg := &Config{}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(bytes, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func init() {
	c, err := LoadConfig("cfg.yaml")
	if err != nil {
		panic(err)
	}

	Conf = c
}
