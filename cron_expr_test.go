package tenco_test

import (
	"testing"

	"github.com/mix3/tenco"
)

func TestCWCronExpr(t *testing.T) {
	cases := []struct {
		cronExpr tenco.CWCronExpr
		expect   string
	}{
		{
			cronExpr: tenco.CWCronExpr{
				Minutes:   [][]int{},
				Hours:     [][]int{},
				DayOfWeek: 0,
			},
			expect: "* * * * ? *",
		},
		{
			cronExpr: tenco.CWCronExpr{
				Minutes:   [][]int{{0}},
				Hours:     [][]int{},
				DayOfWeek: 0,
			},
			expect: "0 * * * ? *",
		},
		{
			cronExpr: tenco.CWCronExpr{
				Minutes:   [][]int{{0}},
				Hours:     [][]int{{0}},
				DayOfWeek: 0,
			},
			expect: "0 0 * * ? *",
		},
		{
			cronExpr: tenco.CWCronExpr{
				Minutes:   [][]int{{0, 1, 2}, {4, 5, 6}},
				Hours:     [][]int{{0, 1, 2}, {4, 5, 6}},
				DayOfWeek: 0,
			},
			expect: "0-2,4-6 0-2,4-6 * * ? *",
		},
		{
			cronExpr: tenco.CWCronExpr{
				Minutes:   [][]int{{0, 1, 2}, {4, 5, 6}},
				Hours:     [][]int{{0, 1, 2}, {4, 5, 6}},
				DayOfWeek: 1,
			},
			expect: "0-2,4-6 0-2,4-6 ? * 1 *",
		},
	}
	for i, cc := range cases {
		if got, want := cc.cronExpr.String(), cc.expect; got != want {
			t.Errorf("[%d] test failed. want %q, but got %q", i, want, got)
		}
	}
}

func TestCronExpr(t *testing.T) {
	cases := []struct {
		cronExpr tenco.CronExpr
		expect   string
	}{
		{
			cronExpr: tenco.CronExpr{
				Minutes:   [][]int{},
				Hours:     [][]int{},
				DayOfWeek: 0,
			},
			expect: "* * * * *",
		},
		{
			cronExpr: tenco.CronExpr{
				Minutes:   [][]int{{0}},
				Hours:     [][]int{},
				DayOfWeek: 0,
			},
			expect: "0 * * * *",
		},
		{
			cronExpr: tenco.CronExpr{
				Minutes:   [][]int{{0}},
				Hours:     [][]int{{0}},
				DayOfWeek: 0,
			},
			expect: "0 0 * * *",
		},
		{
			cronExpr: tenco.CronExpr{
				Minutes:   [][]int{{0, 1, 2}, {4, 5, 6}},
				Hours:     [][]int{{0, 1, 2}, {4, 5, 6}},
				DayOfWeek: 0,
			},
			expect: "0-2,4-6 0-2,4-6 * * *",
		},
		{
			cronExpr: tenco.CronExpr{
				Minutes:   [][]int{{0, 1, 2}, {4, 5, 6}},
				Hours:     [][]int{{0, 1, 2}, {4, 5, 6}},
				DayOfWeek: 1,
			},
			expect: "0-2,4-6 0-2,4-6 * * 0",
		},
	}
	for i, cc := range cases {
		if got, want := cc.cronExpr.String(), cc.expect; got != want {
			t.Errorf("[%d] test failed. want %q, but got %q", i, want, got)
		}
	}
}
