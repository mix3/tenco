package tenco

import (
	"fmt"
)

type Config struct {
	Events []*ConfigEvent `yaml:"events"`
}

type ConfigEvent struct {
	Name                  string                 `yaml:"name"`
	Description           string                 `yaml:"description"`
	Schedule              *Schedule              `yaml:"schedule"`
	CloudwatchEventTarget map[string]interface{} `yaml:"cloudwatch_event_target"`
	IsEnabled             bool                   `yaml:"is_enabled"`
}

func (s *ConfigEvent) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type configEvent ConfigEvent
	var ret configEvent
	if err := unmarshal(&ret); err != nil {
		return err
	}
	if ret.Name == "" {
		return fmt.Errorf("Validation failed. field `name` required")
	}
	if ret.Schedule == nil {
		return fmt.Errorf("Validation failed. field `schedule` required")
	}
	*s = ConfigEvent(ret)
	return nil
}
