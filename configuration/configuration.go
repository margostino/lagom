package configuration

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Params struct {
	InitialLoad             int           `yaml:"initial_load" json:"initial_load"`
	StepLoad                int           `yaml:"step_load" json:"step_load"`
	StepTime                time.Duration `yaml:"step_time" json:"step_time"`
	BufferTime              int           `yaml:"buffer_time" json:"buffer_time"`
	SpawnRate               int           `yaml:"spawn_rate" json:"spawn_rate"`
	MaxLoad                 int           `yaml:"max_load" json:"max_load"`
	RunTime                 time.Duration `yaml:"run_time" json:"run_time"`
	SpikeMultiplierIncrease int           `yaml:"spike_multiplier_increase" json:"spike_multiplier_increase"`
	MaxSpikeMultiplier      int           `yaml:"max_spike_multiplier" json:"max_spike_multiplier"`
	SpikeDuration           time.Duration `yaml:"spike_duration" json:"spike_duration"`
	SpikeCooldownDuration   time.Duration `yaml:"spike_cooldown_duration" json:"spike_cooldown_duration"`
}

type Http struct {
	Url         string `yaml:"url"`
	RequestFile string `yaml:"requestFile"`
	Method      string `yaml:"method"`
	ContentType string `yaml:"contentType"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
}

type Configuration struct {
	Name        string  `yaml:"name"`
	Description string  `yaml:"description"`
	Enabled     bool    `yaml:"enabled"`
	Params      *Params `yaml:"params"`
	Http        *Http   `yaml:"http"`
}

func GetConfig(resource string) *Configuration {
	var configuration Configuration
	ymlFile, err := ioutil.ReadFile(resource)

	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	ymlFile = []byte(os.ExpandEnv(string(ymlFile)))
	err = yaml.Unmarshal(ymlFile, &configuration)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return &configuration
}
