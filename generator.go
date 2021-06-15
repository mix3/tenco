package tenco

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

type Generator interface {
	Generate(w io.Writer, conf Config, offset int) error
}

type JSONGenerator struct {
}

func (g *JSONGenerator) Generate(w io.Writer, conf Config, offset int) error {
	type EventRule struct {
		Name               string `json:"name"`
		Description        string `json:"description"`
		ScheduleExpression string `json:"schedule_expression"`
		IsEnabled          bool   `json:"is_enabled"`
	}
	type Resource struct {
		AWSCloudWatchEventRule   map[string]EventRule              `json:"aws_cloudwatch_event_rule"`
		AWSCloudWatchEventTarget map[string]map[string]interface{} `json:"aws_cloudwatch_event_target"`
	}

	r := Resource{
		AWSCloudWatchEventRule:   map[string]EventRule{},
		AWSCloudWatchEventTarget: map[string]map[string]interface{}{},
	}
	for _, event := range conf.Events {
		for i, expr := range event.Schedule.CronExprs(offset) {
			eventName := fmt.Sprintf("%s-%d", event.Name, i)
			r.AWSCloudWatchEventRule[eventName] = EventRule{
				Name:               eventName,
				Description:        fmt.Sprintf("%s\norig cron(%s)", event.Description, event.Schedule.OrigSchedule),
				ScheduleExpression: fmt.Sprintf("cron(%s)", expr),
				IsEnabled:          event.IsEnabled,
			}
			event.CloudwatchEventTarget["rule"] = fmt.Sprintf("${aws_cloudwatch_event_rule.%s.name}", eventName)
			r.AWSCloudWatchEventTarget[eventName] = convertMapToKeyString(event.CloudwatchEventTarget).(map[string]interface{})
		}
	}
	v := map[string]Resource{
		"resource": r,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func convertMapToKeyString(v interface{}) interface{} {
	tv := reflect.TypeOf(v)
	if tv.Kind() != reflect.Map {
		return v
	}
	vv := reflect.ValueOf(v)
	rv := make(map[string]interface{}, vv.Len())
	for _, mk := range vv.MapKeys() {
		rv[fmt.Sprintf("%s", mk.Interface())] = convertMapToKeyString(vv.MapIndex(mk).Interface())
	}
	return rv
}
