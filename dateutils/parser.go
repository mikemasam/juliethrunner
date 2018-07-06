package dateutils

import (
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
)

func GetUnixTimeFromZero(tm string) (int64, bool) {
	t, err := dateparse.ParseStrict(tm)
	if err == nil {
		return t.Unix(), false
	}

	//tm = strings.Replace(tm, " ", "", -1)
	var rtn int64 = 0
	switch {
	case strings.Contains(tm, "day"):
		{
			h := removeTimeDateNames(tm)
			rtn = int64(toInt(h) * 86400)
			break
		}
	case strings.Contains(tm, "hour"):
		{
			h := removeTimeDateNames(tm)
			rtn = int64(toInt(h) * 3600)
			break
		}
	case strings.Contains(tm, "minute"):
		{
			h := removeTimeDateNames(tm)
			rtn = int64(toInt(h) * 60)
			break
		}
	case strings.Contains(tm, "second"):
		{
			h := removeTimeDateNames(tm)
			rtn = int64(toInt(h))
			break
		}
	case strings.Contains(tm, "time"):
		{
			h := removeTimeDateNames(tm)
			_t, _err := dateparse.ParseStrict("1970-01-01 " + h)
			if _err == nil {
				rtn = _t.Unix()
			}
			break
		}
	}
	return rtn, true
}

func removeTimeDateNames(value string) string {
	r := strings.NewReplacer(
		"seconds", "",
		"second", "",
		"minutes", "",
		"minute", "",
		"hours", "",
		"hour", "",
		"days", "",
		"day", "",
		"years", "",
		"year", "",
		"time", "")
	result := r.Replace(value)
	return result
}
func toInt(result string) int {
	result = strings.Replace(result, " ", "", -1)
	i, err := strconv.Atoi(result)
	if err != nil {
		return 0
	}
	return i
}

func currentEpochTime() int64 {
	return time.Now().Unix()
}
