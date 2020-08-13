package tenco

import (
	"fmt"
	"strings"
)

type Config struct {
	Events []*ConfigEvent `yaml:"events"`
}

type ConfigEvent struct {
	Name                  string                `yaml:"name"`
	Description           string                `yaml:"description"`
	Schedule              Schedule              `yaml:"schedule"`
	CloudwatchEventTarget CloudwatchEventTarget `yaml:"cloudwatch_event_target"`
	ContainerOverrides    []ContainerOverride   `yaml:"container_overrides"`
	IsEnabled             bool                  `yaml:"is_enabled"`
}

func (c ConfigEvent) ContainerOverridesWithEnv(addEnvs []Environment) []ContainerOverride {
	ret := make([]ContainerOverride, 0, len(c.ContainerOverrides))
	for _, co := range c.ContainerOverrides {
		ret = append(ret, ContainerOverride{
			Name:        co.Name,
			Command:     co.Command,
			Environment: append(co.Environment, addEnvs...),
		})
	}
	return ret
}

type CloudwatchEventTarget struct {
	Arn       string    `yaml:"arn"        json:"arn"`
	Rule      string    `yaml:"-"          json:"rule"`
	RoleArn   string    `yaml:"role_arn"   json:"role_arn"`
	ECSTarget ECSTarget `yaml:"ecs_target" json:"ecs_target"`
}

type ECSTarget struct {
	LaunchType           string               `yaml:"launch_type"           json:"launch_type"`
	PlatformVersion      string               `yaml:"platform_version"      json:"platform_version"`
	TaskDefinitionArn    string               `yaml:"task_definition_arn"   json:"task_definition_arn"`
	NetworkConfiguration NetworkConfiguration `yaml:"network_configuration" json:"network_configuration"`
}

type NetworkConfiguration struct {
	Subnets        *StringWithArray `yaml:"subnets"          json:"subnets"`
	SecurityGroups *StringWithArray `yaml:"security_groups"  json:"security_groups"`
	AssignPublicIP bool             `yaml:"assign_public_ip" json:"assign_public_ip"`
}

type ContainerOverride struct {
	Name        string        `json:"name"`
	Command     []string      `json:"command"`
	Environment []Environment `json:"environment"`
}

type containerOverride struct {
	Name        string        `yaml:"name"`
	Entrypoint  []string      `yaml:"entrypoint"`
	Command     []string      `yaml:"command"`
	Environment []Environment `yaml:"environment"`
}

type Environment struct {
	Name  string `yaml:"name"  json:"name"`
	Value string `yaml:"value" json:"value"`
}

func (c *ContainerOverride) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var ret containerOverride
	if err := unmarshal(&ret); err != nil {
		return err
	}
	c.Name = ret.Name
	c.Command = append(ret.Entrypoint, ret.Command...)
	c.Environment = ret.Environment
	return nil
}

type StringWithArray struct {
	String string
	Array  []string
}

func (s *StringWithArray) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var ret interface{}
	if err := unmarshal(&ret); err != nil {
		return err
	}

	switch t := ret.(type) {
	case string:
		s.String = t
	case []interface{}:
		for _, v := range t {
			s.Array = append(s.Array, v.(string))
		}
	default:
		return fmt.Errorf("Unexpected parse error. %v", ret)
	}
	return nil
}

func (s *StringWithArray) MarshalJSON() ([]byte, error) {
	if len(s.Array) == 0 {
		return []byte(fmt.Sprintf(`"%s"`, s.String)), nil
	}
	return []byte(fmt.Sprintf(`["%s"]`, strings.Join(s.Array, `","`))), nil
}
