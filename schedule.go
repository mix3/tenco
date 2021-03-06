package tenco

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Schedule struct {
	Minutes      Minutes
	Hours        Hours
	DayOfWeeks   DayOfWeeks
	OrigSchedule OrigSchedule
}

type schedule struct {
	Minutes    Minutes    `yaml:"minutes"`
	Hours      Hours      `yaml:"hours"`
	DayOfWeeks DayOfWeeks `yaml:"day_of_weeks"`
}

type OrigSchedule struct {
	Minutes    string `yaml:"minutes"`
	Hours      string `yaml:"hours"`
	DayOfWeeks string `yaml:"day_of_weeks"`
}

func (o OrigSchedule) String() string {
	m := o.Minutes
	if m == "" {
		m = "*"
	}

	h := o.Hours
	if h == "" {
		h = "*"
	}

	var day string
	if o.DayOfWeeks == "" {
		day = "*"
	} else {
		day = "?"
	}

	month := "*"

	dow := o.DayOfWeeks
	if o.DayOfWeeks == "" {
		dow = "?"
	}

	year := "*"
	return strings.Join([]string{m, h, day, month, dow, year}, " ")
}

func (s *Schedule) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var ret schedule
	if err := unmarshal(&ret); err != nil {
		return err
	}
	var orig OrigSchedule
	if err := unmarshal(&orig); err != nil {
		return err
	}
	*s = Schedule{
		Minutes:      ret.Minutes,
		Hours:        ret.Hours,
		DayOfWeeks:   ret.DayOfWeeks,
		OrigSchedule: orig,
	}
	return nil
}

func (s *Schedule) CronExprs(offset int) []CronExpr {
	if len(s.DayOfWeeks) == 0 {
		return []CronExpr{
			{
				Minutes: string(s.Minutes),
				Hours:   s.Hours.offset(offset).unique().merge(),
			},
		}
	}

	merge := make(map[DayOfWeek]Hours, 7)
	vecToHours := s.Hours.offsetWithVec(offset)
	for _, dow := range s.DayOfWeeks {
		for vec, hours := range vecToHours {
			w := dow.shift(vec)
			merge[w] = append(merge[w], hours...)
		}
	}

	ret := make([]CronExpr, 0, 7)
	for i := 1; i <= 7; i++ {
		w := DayOfWeek(i)
		if v, ok := merge[w]; ok {
			ret = append(ret, CronExpr{
				Minutes:   string(s.Minutes),
				Hours:     v.unique().merge(),
				DayOfWeek: i,
			})
		}
	}
	return ret
}

type Minutes string

func (ms *Minutes) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}

	// validate
	for _, v := range strings.Split(str, ",") {
		switch {
		// n-m, *-m
		case dashReg.MatchString(v):
			r := strings.Split(v, "-")
			if r[0] != "*" {
				if _, err := strconv.Atoi(r[0]); err != nil {
					return fmt.Errorf("parse start minutes failed. %q:%q %w", str, v, err)
				}
			}
			if _, err := strconv.Atoi(r[1]); err != nil {
				return fmt.Errorf("parse end minutes failed. %q:%q %w", str, v, err)
			}
		// n/m, */m
		case slashReg.MatchString(v):
			r := strings.Split(v, "/")
			if r[0] != "*" {
				if _, err := strconv.Atoi(r[0]); err != nil {
					return fmt.Errorf("parse start minutes failed. %q:%q %w", str, v, err)
				}
			}
			if _, err := strconv.Atoi(r[1]); err != nil {
				return fmt.Errorf("parse increment minutes failed. %q:%q %w", str, v, err)
			}
		// n, *
		default:
			if v != "" && v != "*" {
				i, err := strconv.Atoi(v)
				if err != nil {
					return fmt.Errorf("parse minutes failed. %q:%q %w", str, v, err)
				}
				if i < 0 || 60 <= i {
					return fmt.Errorf("hours allow range [0-59]")
				}
			}
		}
	}

	*ms = Minutes(str)

	return nil
}

type Hour int

func (h Hour) offset(offset int) (Hour, int) {
	var vec int
	ret := h + Hour(offset)
	switch {
	case ret < 0:
		ret += 24
		vec--
	case 24 <= ret:
		ret %= 24
		vec++
	}
	return ret, vec
}

type Hours []Hour

func (hs Hours) raw() []int {
	ret := make([]int, 0, len(hs))
	for _, h := range hs {
		ret = append(ret, int(h))
	}
	return ret
}

var (
	dashReg  = regexp.MustCompile(`^(\*|\w+)-(\w+)$`)
	slashReg = regexp.MustCompile(`^(\*|\w+)/(\w+)$`)
)

func (hs *Hours) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}

	// "" "*" -> 0,1,2,...
	if str == "" || str == "*" {
		for i := 0; i < 24; i++ {
			hs.add(i)
		}
		return nil
	}

	for _, v := range strings.Split(str, ",") {
		switch {
		// n-m 1-5  -> 1,2,3,4,5
		//     23-1 -> 23,0,1
		case dashReg.MatchString(v):
			r := strings.Split(v, "-")
			start, err := strconv.Atoi(r[0])
			if err != nil {
				return fmt.Errorf("parse start hours failed. %q:%q %w", str, v, err)
			}
			end, err := strconv.Atoi(r[1])
			if err != nil {
				return fmt.Errorf("parse end hours failed. %q:%q %w", str, v, err)
			}
			for _, v := range loopRange(start, end, 0, 23) {
				if err := hs.add(v); err != nil {
					return fmt.Errorf("parse hours failed. %q:%q %w", str, v, err)
				}
			}
		// n/m 19/2 -> 19,21,23
		case slashReg.MatchString(v):
			r := strings.Split(v, "/")
			start, err := strconv.Atoi(r[0])
			if err != nil {
				return fmt.Errorf("parse start hours failed. %q:%q %w", str, v, err)
			}
			incr, err := strconv.Atoi(r[1])
			if err != nil {
				return fmt.Errorf("parse increment hours failed. %q:%q %w", str, v, err)
			}
			for i := start; i < 24; i += incr {
				if err := hs.add(i); err != nil {
					return fmt.Errorf("parse hours failed. %q:%q %w", str, v, err)
				}
			}
		// n 1 -> 1
		default:
			i, err := strconv.Atoi(v)
			if err != nil {
				return fmt.Errorf("parse hours failed. %q:%q %w", str, v, err)
			}
			if err := hs.add(i); err != nil {
				return fmt.Errorf("parse hours failed. %q:%q %w", str, v, err)
			}
		}
	}

	return nil
}

func (hs *Hours) add(i int) error {
	if i < 0 || 24 <= i {
		return fmt.Errorf("hours allow range [0-23]")
	}
	*hs = append(*hs, Hour(i))
	return nil
}

func (hs Hours) offset(offset int) Hours {
	ret := make(Hours, 0, len(hs))
	for _, h := range hs {
		v, _ := h.offset(offset)
		ret = append(ret, v)
	}
	return ret
}

func (hs Hours) offsetWithVec(offset int) map[int]Hours {
	ret := make(map[int]Hours, 2)
	for _, v := range hs {
		v, vec := v.offset(offset)
		ret[vec] = append(ret[vec], v)
	}
	return ret
}

func (hs Hours) unique() Hours {
	ret := make(Hours, 0, len(hs))
	dup := make(map[Hour]struct{}, len(hs))
	for _, h := range hs {
		if _, ok := dup[h]; !ok {
			ret = append(ret, h)
			dup[h] = struct{}{}
		}
	}
	return ret
}

func (hs Hours) merge() [][]int {
	hours := hs.raw()
	sort.Ints(hours)
	ret := make([][]int, 0, len(hours))
	tmp := make([]int, 0, 24)
	for _, h := range hours {
		if len(tmp) == 0 || tmp[len(tmp)-1]+1 == int(h) {
			tmp = append(tmp, int(h))
		} else {
			ret = append(ret, tmp)
			tmp = make([]int, 0, 24)
			tmp = append(tmp, int(h))
		}
	}
	ret = append(ret, tmp)
	return ret
}

type DayOfWeek int

func (w DayOfWeek) shift(vec int) DayOfWeek {
	ret := w + DayOfWeek(vec)
	switch {
	case ret < 1:
		ret += 7
	case 7 <= ret:
		ret %= 7
	}
	return ret
}

type DayOfWeeks []DayOfWeek

var dayOfWeekNumMap = map[string]int{
	"SUN": 1,
	"MON": 2,
	"TUE": 3,
	"WED": 4,
	"THU": 5,
	"FRI": 6,
	"SAT": 7,
	"1":   1,
	"2":   2,
	"3":   3,
	"4":   4,
	"5":   5,
	"6":   6,
	"7":   7,
}

func createRange(s, e int) []int {
	ret := make([]int, 0, e-s+1)
	for n := s; n <= e; n++ {
		ret = append(ret, n)
	}
	return ret
}

// loopRange(5, 2, 1, 7) -> [5, 6, 7, 1, 2]
func loopRange(s, e, rs, re int) []int {
	r := createRange(rs, re)
	var i int
	for n := rs; n <= re; n++ {
		r = append(r, n)
		if n < s {
			i++
		}
	}
	r = append(r[i:], r[:i]...)
	ret := make([]int, 0, len(r))
	for _, v := range r {
		ret = append(ret, v)
		if v == e {
			break
		}
	}
	return ret
}

func isStr(str string) bool {
	switch str {
	case "SUN", "MON", "TUE", "WED", "THU", "FRI", "SAT":
		return true
	}
	return false
}

func (ws *DayOfWeeks) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}

	if str == "" {
		return nil
	}

	for _, v := range strings.Split(str, ",") {
		switch {
		// n-m 1-5     -> 1,2,3,4,5
		//     5-2     -> 5,6,7,1,2
		//     WED-MON -> 4,5,6,7,1,2
		case dashReg.MatchString(v):
			r := strings.Split(v, "-")
			start, ok := dayOfWeekNumMap[r[0]]
			if !ok {
				return fmt.Errorf("parse start day_of_weeks failed. %q:%q", str, v)
			}
			end, ok := dayOfWeekNumMap[r[1]]
			if !ok {
				return fmt.Errorf("parse end day_of_weeks failed. %q:%q", str, v)
			}
			if st, et := isStr(r[0]), isStr(r[1]); st != et {
				return fmt.Errorf("parse day_of_weeks failed. %q:%q", str, v)
			}
			for _, n := range loopRange(start, end, 1, 7) {
				*ws = append(*ws, DayOfWeek(n))
			}
		// n MON -> 2
		//   7   -> 7
		default:
			n, ok := dayOfWeekNumMap[v]
			if !ok {
				return fmt.Errorf("parse day_of_weeks failed. %q:%q", str, v)
			}
			*ws = append(*ws, DayOfWeek(n))
		}
	}

	*ws = ws.unique()

	return nil
}

func (ws DayOfWeeks) unique() DayOfWeeks {
	ret := make(DayOfWeeks, 0, len(ws))
	dup := make(map[DayOfWeek]struct{}, len(ws))
	for _, w := range ws {
		if _, ok := dup[w]; !ok {
			ret = append(ret, w)
			dup[w] = struct{}{}
		}
	}
	return ret
}
