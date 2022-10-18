package configuration

import "time"

const (
	SuccessEnabled = iota
	FailureEnabled
	RandomnessEnabled
)

type Params struct {
	InitialLoad             int           `yaml:"initial_load"`
	StepLoad                int           `yaml:"step_load"`
	StepTime                time.Duration `yaml:"step_time"`
	SpawnRate               int           `yaml:"spawn_rate"`
	MaxLoad                 int           `yaml:"max_load"`
	RunTime                 time.Duration `yaml:"run_time"`
	SpikeMultiplierIncrease int           `yaml:"spike_multiplier_increase"`
	MaxSpikeMultiplier      int           `yaml:"max_spike_multiplier"`
	SpikeDuration           time.Duration `yaml:"spike_duration"`
	SpikeCooldownDuration   time.Duration `yaml:"spike_cooldown_duration"`
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
	Http        *Http   `yaml:"maxStepTime"`
}
