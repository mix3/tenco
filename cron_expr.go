package tenco

import (
	"fmt"
	"strconv"
	"strings"
)

type CWCronExpr struct {
	Minutes   [][]int
	Hours     [][]int
	DayOfWeek int
}

func (c CWCronExpr) String() string {
	var m string
	if len(c.Minutes) == 0 {
		m = "*"
	} else {
		var tmp []string
		for _, v := range c.Minutes {
			if len(v) == 60 {
				tmp = append(tmp, "*")
			} else if 1 < len(v) {
				tmp = append(tmp, fmt.Sprintf("%d-%d", v[0], v[len(v)-1]))
			} else {
				tmp = append(tmp, strconv.Itoa(v[0]))
			}
		}
		m = strings.Join(tmp, ",")
	}

	var h string
	if len(c.Hours) == 0 {
		h = "*"
	} else {
		var tmp []string
		for _, v := range c.Hours {
			if len(v) == 24 {
				tmp = append(tmp, "*")
			} else if 1 < len(v) {
				tmp = append(tmp, fmt.Sprintf("%d-%d", v[0], v[len(v)-1]))
			} else {
				tmp = append(tmp, strconv.Itoa(v[0]))
			}
		}
		h = strings.Join(tmp, ",")
	}

	var day string
	if c.DayOfWeek == 0 {
		day = "*"
	} else {
		day = "?"
	}

	month := "*"

	dow := strconv.Itoa(c.DayOfWeek)
	if c.DayOfWeek == 0 {
		dow = "?"
	}

	year := "*"

	return strings.Join([]string{m, h, day, month, dow, year}, " ")
}

type CronExpr struct {
	Minutes   [][]int
	Hours     [][]int
	DayOfWeek int
}

func (c CronExpr) String() string {
	var m string
	if len(c.Minutes) == 0 {
		m = "*"
	} else {
		var tmp []string
		for _, v := range c.Minutes {
			if len(v) == 60 {
				tmp = append(tmp, "*")
			} else if 1 < len(v) {
				tmp = append(tmp, fmt.Sprintf("%d-%d", v[0], v[len(v)-1]))
			} else {
				tmp = append(tmp, strconv.Itoa(v[0]))
			}
		}
		m = strings.Join(tmp, ",")
	}

	var h string
	if len(c.Hours) == 0 {
		h = "*"
	} else {
		var tmp []string
		for _, v := range c.Hours {
			if len(v) == 24 {
				tmp = append(tmp, "*")
			} else if 1 < len(v) {
				tmp = append(tmp, fmt.Sprintf("%d-%d", v[0], v[len(v)-1]))
			} else {
				tmp = append(tmp, strconv.Itoa(v[0]))
			}
		}
		h = strings.Join(tmp, ",")
	}

	day := "*"
	month := "*"

	dow := strconv.Itoa(c.DayOfWeek - 1)
	if c.DayOfWeek == 0 {
		dow = "*"
	}

	return strings.Join([]string{m, h, day, month, dow}, " ")
}
