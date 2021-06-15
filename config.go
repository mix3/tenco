package tenco

type Config struct {
	Events []*ConfigEvent `yaml:"events"`
}

type ConfigEvent struct {
	Name                  string                 `yaml:"name"`
	Description           string                 `yaml:"description"`
	Schedule              Schedule               `yaml:"schedule"`
	CloudwatchEventTarget map[string]interface{} `yaml:"cloudwatch_event_target"`
	IsEnabled             bool                   `yaml:"is_enabled"`
}
