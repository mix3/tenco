package tenco_test

import (
	"testing"

	"github.com/mix3/tenco"
)

func TestCWCronExpr(t *testing.T) {
	minutesAsterisc := []int{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
		10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
		20, 21, 22, 23, 24, 25, 26, 27, 28, 29,
		30, 31, 32, 33, 34, 35, 36, 37, 38, 39,
		40, 41, 42, 43, 44, 45, 46, 47, 48, 49,
		50, 51, 52, 53, 54, 55, 56, 57, 58, 59,
	}
	hoursAsterisc := []int{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23,
	}
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
				Minutes:   [][]int{minutesAsterisc},
				Hours:     [][]int{hoursAsterisc},
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
	minutesAsterisc := []int{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
		10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
		20, 21, 22, 23, 24, 25, 26, 27, 28, 29,
		30, 31, 32, 33, 34, 35, 36, 37, 38, 39,
		40, 41, 42, 43, 44, 45, 46, 47, 48, 49,
		50, 51, 52, 53, 54, 55, 56, 57, 58, 59,
	}
	hoursAsterisc := []int{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23,
	}
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
				Minutes:   [][]int{minutesAsterisc},
				Hours:     [][]int{hoursAsterisc},
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
