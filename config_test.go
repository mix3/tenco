package tenco_test

import (
	"strings"
	"testing"

	"github.com/mix3/tenco"
	"gopkg.in/yaml.v2"
)

func TestEventConfig(t *testing.T) {
	cases := []struct {
		input string
		err   string
	}{
		{
			input: `
schedule:
  minutes: "0"
  hours:   "0"
`,
			err: "Validation failed. field `name` required",
		},
		{
			input: `
name: test
`,
			err: "Validation failed. field `schedule` required",
		},
		{
			input: `
name: test
schedule:
  minutes: "0"
  hours:   "0"
`,
		},
	}
	for i, c := range cases {
		var ec tenco.ConfigEvent
		err := yaml.Unmarshal([]byte(c.input), &ec)
		if err != nil {
			if c.err == "" {
				t.Errorf("[%d] test failed. %s", i, err)
			} else if !strings.Contains(err.Error(), c.err) {
				t.Errorf("[%d] test failed. %q does not contain %q", i, err.Error(), c.err)
			}
			continue
		}
		if c.err != "" {
			t.Errorf("[%d] test failed", i)
			continue
		}
	}
}
