package tenco

import (
	"encoding/json"
	"fmt"
	"io"
)

type Generator interface {
	Generate(w io.Writer, conf Config, offset int) error
}

type JsonGenerator struct {
}

func (g *JsonGenerator) Generate(w io.Writer, conf Config, offset int) error {
	v := map[string]map[string]map[string]interface{}{
		"resource": {
			"aws_cloudwatch_event_rule":   {},
			"aws_cloudwatch_event_target": {},
		},
	}
	type EventRule struct {
		Name               string `json:"name"`
		Description        string `json:"description"`
		ScheduleExpression string `json:"schedule_expression"`
		IsEnabled          bool   `json:"is_enabled"`
	}
	type EventTarget struct {
		CloudwatchEventTarget
		TargetID string `json:"target_id"`
		Input    string `json:"input"`
	}
	for _, event := range conf.Events {
		for i, expr := range event.Schedule.CronExprs(offset) {
			eventName := fmt.Sprintf("%s-%d", event.Name, i)
			v["resource"]["aws_cloudwatch_event_rule"][eventName] = EventRule{
				Name:               eventName,
				Description:        fmt.Sprintf("%s\norig cron(%s)", event.Description, event.Schedule.OrigSchedule),
				ScheduleExpression: fmt.Sprintf("cron(%s)", expr),
				IsEnabled:          event.IsEnabled,
			}
			input := map[string][]ContainerOverride{
				"containerOverrides": event.ContainerOverridesWithEnv([]Environment{
					{
						Name:  "CRON_SCHEDULE",
						Value: event.Schedule.OrigSchedule.String(),
					},
					{
						Name:  "CRON_DESCRIPTION",
						Value: event.Description,
					},
				}),
			}
			b, err := json.Marshal(input)
			if err != nil {
				return fmt.Errorf("event.ContainerOverride json Marshal failed. %w", err)
			}
			v["resource"]["aws_cloudwatch_event_target"][eventName] = EventTarget{
				CloudwatchEventTarget: CloudwatchEventTarget{
					Arn:       event.CloudwatchEventTarget.Arn,
					Rule:      fmt.Sprintf("${aws_cloudwatch_event_rule.%s.name}", eventName),
					RoleArn:   event.CloudwatchEventTarget.RoleArn,
					ECSTarget: event.CloudwatchEventTarget.ECSTarget,
				},
				TargetID: event.Name,
				Input:    string(b),
			}
		}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
