package tenco_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/mix3/tenco"
	"gopkg.in/yaml.v2"
)

func TestMinutes(t *testing.T) {
	cases := []struct {
		input  string
		err    bool
		expect tenco.Minutes
	}{
		{
			input: "-1",
			err:   true,
		},
		{
			input: "60",
			err:   true,
		},
		{
			input: "-1-4",
			err:   true,
		},
		{
			input: "-1/4",
			err:   true,
		},
		{
			input:  "",
			expect: tenco.Minutes(""),
		},
		{
			input:  "*",
			expect: tenco.Minutes("*"),
		},
		{
			input:  "1/4",
			expect: tenco.Minutes("1/4"),
		},
		{
			input:  "1-4",
			expect: tenco.Minutes("1-4"),
		},
		{
			input:  "*-4",
			expect: tenco.Minutes("*-4"),
		},
		{
			input:  "*,*-4,*/9",
			expect: tenco.Minutes("*,*-4,*/9"),
		},
	}
	for i, c := range cases {
		var s tenco.Schedule
		err := yaml.Unmarshal([]byte(fmt.Sprintf(`minutes: "%s"`, c.input)), &s)
		if err != nil {
			if !c.err {
				t.Errorf("[%d] test failed. %s", i, err)
			}
			continue
		}
		if c.err {
			if err == nil {
				t.Errorf("[%d] test failed", i)
			}
			continue
		}
		if g, w := s.Minutes, c.expect; !reflect.DeepEqual(g, w) {
			t.Errorf("[%d] test faield. want %+v, but got %+v", i, w, g)
		}
	}
}

func TestHours(t *testing.T) {
	cases := []struct {
		input  string
		err    bool
		expect tenco.Hours
	}{
		{
			input: "-1",
			err:   true,
		},
		{
			input: "24",
			err:   true,
		},
		{
			input: "-1-4",
			err:   true,
		},
		{
			input: "-1/4",
			err:   true,
		},
		{
			input: "*-4",
			err:   true,
		},
		{
			input: "*/4",
			err:   true,
		},
		{
			input:  "",
			expect: tenco.Hours{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
		},
		{
			input:  "*",
			expect: tenco.Hours{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
		},
		{
			input:  "1,5,9",
			expect: tenco.Hours{1, 5, 9},
		},
		{
			input:  "1-4",
			expect: tenco.Hours{1, 2, 3, 4},
		},
		{
			input:  "23-2",
			expect: tenco.Hours{23, 0, 1, 2},
		},
		{
			input:  "19/2",
			expect: tenco.Hours{19, 21, 23},
		},
		{
			input:  "1,14-16,19/2",
			expect: tenco.Hours{1, 14, 15, 16, 19, 21, 23},
		},
		{
			input:  "1,14-16,19/2,23-0",
			expect: tenco.Hours{1, 14, 15, 16, 19, 21, 23, 23, 0},
		},
	}
	for i, c := range cases {
		var s tenco.Schedule
		err := yaml.Unmarshal([]byte(fmt.Sprintf(`hours: "%s"`, c.input)), &s)
		if err != nil {
			if !c.err {
				t.Errorf("[%d] test failed. %s", i, err)
			}
			continue
		}
		if c.err {
			if err == nil {
				t.Errorf("[%d] test failed", i)
			}
			continue
		}
		if g, w := s.Hours, c.expect; !reflect.DeepEqual(g, w) {
			t.Errorf("[%d] test faield. want %+v, but got %+v", i, w, g)
		}
	}
}

func TestDayOfWeeks(t *testing.T) {
	cases := []struct {
		input  string
		err    bool
		expect tenco.DayOfWeeks
	}{
		{
			input: "-1",
			err:   true,
		},
		{
			input: "XXX",
			err:   true,
		},
		{
			input: "MON-XXX",
			err:   true,
		},
		{
			input: "1-MON",
			err:   true,
		},
		{
			input:  "MON-FRI",
			expect: tenco.DayOfWeeks{2, 3, 4, 5, 6},
		},
		{
			input:  "FRI-MON",
			expect: tenco.DayOfWeeks{6, 7, 1, 2},
		},
		{
			input:  "FRI-MON,SUN-TUE",
			expect: tenco.DayOfWeeks{6, 7, 1, 2, 3},
		},
		{
			input:  "1,MON,3,WED,THU-SAT",
			expect: tenco.DayOfWeeks{1, 2, 3, 4, 5, 6, 7},
		},
	}
	for i, c := range cases {
		var s tenco.Schedule
		err := yaml.Unmarshal([]byte(fmt.Sprintf("day_of_weeks: %s", c.input)), &s)
		if err != nil {
			if !c.err {
				t.Errorf("[%d] test failed. %s", i, err)
			}
			continue
		}
		if c.err {
			if err == nil {
				t.Errorf("[%d] test failed", i)
			}
			continue
		}
		if g, w := s.DayOfWeeks, c.expect; !reflect.DeepEqual(g, w) {
			t.Errorf("[%d] test faield. want %+v, but got %+v", i, w, g)
		}
	}
}

func TestCronExprs(t *testing.T) {
	cases := []struct {
		minutes    string
		hours      string
		dayOfWeeks string
		offset     int
		err        bool
		expect     []tenco.CronExpr
	}{
		{
			minutes:    "",
			hours:      "",
			dayOfWeeks: "",
			offset:     0,
			expect: []tenco.CronExpr{
				{
					Minutes: "",
					Hours: [][]int{
						{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
					},
					DayOfWeek: 0,
				},
			},
		},
		{
			minutes:    "",
			hours:      "0",
			dayOfWeeks: "",
			offset:     0,
			expect: []tenco.CronExpr{
				{
					Minutes: "",
					Hours: [][]int{
						{0},
					},
					DayOfWeek: 0,
				},
			},
		},
		{
			minutes:    "0",
			hours:      "0",
			dayOfWeeks: "SUN",
			offset:     0,
			expect: []tenco.CronExpr{
				{
					Minutes: "0",
					Hours: [][]int{
						{0},
					},
					DayOfWeek: 1,
				},
			},
		},
		{
			minutes:    "",
			hours:      "0",
			dayOfWeeks: "",
			offset:     -9,
			expect: []tenco.CronExpr{
				{
					Minutes: "",
					Hours: [][]int{
						{15},
					},
					DayOfWeek: 0,
				},
			},
		},
		{
			minutes:    "0",
			hours:      "0",
			dayOfWeeks: "SUN",
			offset:     -9,
			expect: []tenco.CronExpr{
				{
					Minutes: "0",
					Hours: [][]int{
						{15},
					},
					DayOfWeek: 7,
				},
			},
		},
		{
			minutes:    "0",
			hours:      "8,9",
			dayOfWeeks: "MON",
			offset:     -9,
			expect: []tenco.CronExpr{
				{
					Minutes: "0",
					Hours: [][]int{
						{23},
					},
					DayOfWeek: 1,
				},
				{
					Minutes: "0",
					Hours: [][]int{
						{0},
					},
					DayOfWeek: 2,
				},
			},
		},
		{
			minutes:    "0",
			hours:      "*",
			dayOfWeeks: "MON-FRI",
			offset:     -9,
			expect: []tenco.CronExpr{
				{
					Minutes: "0",
					Hours: [][]int{
						{15, 16, 17, 18, 19, 20, 21, 22, 23},
					},
					DayOfWeek: 1,
				},
				{
					Minutes: "0",
					Hours: [][]int{
						{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
					},
					DayOfWeek: 2,
				},
				{
					Minutes: "0",
					Hours: [][]int{
						{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
					},
					DayOfWeek: 3,
				},
				{
					Minutes: "0",
					Hours: [][]int{
						{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
					},
					DayOfWeek: 4,
				},
				{
					Minutes: "0",
					Hours: [][]int{
						{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23},
					},
					DayOfWeek: 5,
				},
				{
					Minutes: "0",
					Hours: [][]int{
						{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
					},
					DayOfWeek: 6,
				},
			},
		},
	}
	for i, c := range cases {
		var s tenco.Schedule
		input := fmt.Sprintf(`
minutes:      "%s"
hours:        "%s"
day_of_weeks: "%s"
`, c.minutes, c.hours, c.dayOfWeeks)
		err := yaml.Unmarshal([]byte(input), &s)
		if err != nil {
			if !c.err {
				t.Errorf("[%d] test failed. %s", i, err)
			}
			continue
		}
		if c.err {
			if err == nil {
				t.Errorf("[%d] test failed", i)
			}
			continue
		}
		if g, w := s.CronExprs(c.offset), c.expect; !reflect.DeepEqual(g, w) {
			t.Errorf("[%d] test faield. want %+v, but got %+v", i, w, g)
		}
	}
}
