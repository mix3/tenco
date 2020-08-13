package tenco

import (
	"fmt"
	"strconv"
	"strings"
)

type CronExpr struct {
	Minutes   string
	Hours     [][]int
	DayOfWeek int
}

func (c CronExpr) String() string {
	m := c.Minutes
	if m == "" {
		m = "*"
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
