package calender

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func TimeToString(tm time.Time) string {
	if tm.IsZero() {
		return ""
	} else {
		D := tm
		y, n, d := D.Date()
		m := int(n)
		d1 := fmt.Sprintf("%02d", d)
		m1 := fmt.Sprintf("%02d", m)
		return strconv.Itoa(y) + "-" + m1 + "-" + d1
	}
}

func StringToTime(str string) time.Time {
	if str == "" {
		return time.Time{}
	} else {
		sstr := strings.Split(str, "-")
		y, _ := strconv.Atoi(sstr[0])
		m, _ := strconv.Atoi(sstr[1])
		d, _ := strconv.Atoi(sstr[2])
		return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
	}
}
