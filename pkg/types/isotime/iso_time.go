package isotime

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

const hhmmssFormat = "15:04:05"
const format = "2006-01-02T15:04:05.000000000Z"
const formatMilli = "2006-01-02T15:04:05.000Z"
const format2 = "2006-01-02T15:04:05Z"
const jsonFormat = `"` + format + `"`
const dateFormat = "2006-01-02"
const datetimeFormat = "2006-01-02 15:04:05"

const DateAbbrFormat = "20060102"

var layouts = []string{
	format,
	formatMilli,
	time.RFC3339,
	format2,
	dateFormat,
	datetimeFormat,
}

var (
	TestNow    ISOTime
	TestNowSet bool
	jsonNull   = []byte("null")
)

// ISOTime ISOTime
type ISOTime time.Time

// SetTime sets the internal time value
func (it *ISOTime) SetTime(t time.Time) {
	*it = ISOTime(t.Local())
}

// Now Now
func Now() ISOTime {
	if TestNowSet {
		return TestNow
	}

	return New(time.Now())
}

// New NewISOTime
func New(t time.Time) ISOTime {
	return ISOTime(t.Local())
}

// NewFromDate .
func NewFromDate(date Date) ISOTime {
	return New(time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location()))
}

// Null .
func Null() ISOTime {
	return New(time.Time{})
}

// Parse .
func Parse(value string) (ISOTime, error) {
	if value == "null" {
		return Null(), nil
	}
	var err error
	var t time.Time
	for _, layout := range layouts {
		t, err = time.Parse(layout, value)
		if err != nil {
			continue
		}

		return New(t), nil
	}

	return New(t), err
}

func Since(v ISOTime) time.Duration {
	return time.Since(v.ToTime())
}

// MarshalJSON MarshalJSON
func (it ISOTime) MarshalJSON() ([]byte, error) {
	if it.ToTime().IsZero() {
		return jsonNull, nil
	}

	return []byte(time.Time(it).UTC().Format(jsonFormat)), nil
}

// UnmarshalJSON UnmarshalJSON
func (it *ISOTime) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, jsonNull) {
		*it = Null()
		return nil
	}

	t := &time.Time{}
	if err := json.Unmarshal(data, t); err != nil {
		return err
	}
	*it = ISOTime(t.Local())

	return nil
}

// UnmarshalParam UnmarshalParam
func (it *ISOTime) UnmarshalParam(param string) error {
	isoTime, err := Parse(param)
	if err != nil {
		return err
	}
	*it = isoTime
	return nil
}

// Scan implements the Scanner interface.
func (it *ISOTime) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		t, err := Parse(string(v))
		if err != nil {
			return err
		}
		*it = t
	case string:
		t, err := Parse(v)
		if err != nil {
			return err
		}
		*it = t
	case time.Time:
		*it = ISOTime(v.Local())
	case nil:
		*it = ISOTime{}
	default:
		return fmt.Errorf("cannot sql.Scan() ISOTime from: %#v", v)
	}
	return nil
}

// Value implements the driver Valuer interface.
func (it ISOTime) Value() (driver.Value, error) {
	if it.ToTime().IsZero() {
		return nil, nil
	}
	return it.ToTime(), nil
}

// String String
func (it *ISOTime) String() string {
	return time.Time(*it).UTC().Format(format)
}

func (it *ISOTime) LocalString() string {
	return time.Time(*it).Format(datetimeFormat)
}

func (it *ISOTime) JsonString() string {
	return time.Time(*it).UTC().Format(jsonFormat)
}

func (it *ISOTime) RawString() string {
	return time.Time(*it).UTC().Format(format)
}

// DateString String
func (it *ISOTime) DateString() string {
	return time.Time(*it).UTC().Format(dateFormat)
}

// ToTime .
func (it ISOTime) ToTime() time.Time {
	return time.Time(it)
}

// IsZero .
func (it ISOTime) IsZero() bool {
	return it.ToTime().IsZero()
}

// Equal .
func (it ISOTime) Equal(v ISOTime) bool {
	return it.ToTime().Equal(v.ToTime())
}

// StartOfWeek 返回本周一的开始时间 (00:00:00)
func (it ISOTime) StartOfWeek() ISOTime {
	weekday := time.Time(it).Local().Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	return ISOTime(time.Time(it).Local().AddDate(0, 0, -int(weekday-1)).Truncate(24 * time.Hour))
}

// EndOfWeek 返回本周日的结束时间 (23:59:59)
func (it ISOTime) EndOfWeek() ISOTime {
	weekday := time.Time(it).Local().Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	endOfWeek := time.Time(it).Local().AddDate(0, 0, 7-int(weekday))
	return ISOTime(time.Date(
		endOfWeek.Year(),
		endOfWeek.Month(),
		endOfWeek.Day(),
		23, 59, 59, 999999999,
		endOfWeek.Location(),
	).Local())
}

// SetHHmmss 设置时分秒
func (it ISOTime) SetHHmmss(value string) (ISOTime, error) {
	parseTime, err := time.Parse(hhmmssFormat, value)
	if err != nil {
		return ISOTime{}, err
	}

	return New(time.Date(it.ToTime().Year(), it.ToTime().Month(), it.ToTime().Day(), parseTime.Hour(), parseTime.Minute(), parseTime.Second(), 0, it.ToTime().Location())), nil
}

// IsSameDay 判断两个时间是否是同一天
func (it ISOTime) IsSameDay(v ISOTime) bool {
	return it.ToTime().Year() == v.ToTime().Year() && it.ToTime().Month() == v.ToTime().Month() && it.ToTime().Day() == v.ToTime().Day()
}

// IsSameMonth 判断两个时间是否是同一个月
func (it ISOTime) IsSameMonth(v ISOTime) bool {
	t := it.ToTime().Local()
	vT := v.ToTime().Local()
	return t.Year() == vT.Year() && t.Month() == vT.Month()
}

// StartOfYear 年初时间（本地时区）
func (it ISOTime) StartOfYear() ISOTime {
	t := it.ToTime().Local()
	return ISOTime(time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location()))
}

// EndOfYear 年末时间（本地时区）
func (it ISOTime) EndOfYear() ISOTime {
	t := it.ToTime().Local()
	return ISOTime(time.Date(t.Year(), 12, 31, 23, 59, 59, 999999999, t.Location()))
}

// StartOfMonth 月初时间（本地时区）
func (it ISOTime) StartOfMonth() ISOTime {
	t := it.ToTime().Local()
	return ISOTime(time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()))
}

// EndOfMonth 月末时间（本地时区）
func (it ISOTime) EndOfMonth() ISOTime {
	t := it.ToTime().Local()
	return ISOTime(time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location()).Add(-time.Millisecond))
}

// Yesterday 昨天时间（本地时区）
func (it ISOTime) Yesterday() ISOTime {
	return it.Add(-time.Hour * 24)
}

// StartOfNextDay 次日0点时间（本地时区）
func (it ISOTime) StartOfNextDay() ISOTime {
	startOfDay := it.StartOfDay()
	startOfNextDay := startOfDay.Add(time.Hour * 24)
	return startOfNextDay
}

// StartOfDay 当日0点时间（本地时区）
func (it ISOTime) StartOfDay() ISOTime {
	t := it.ToTime().Local()
	return ISOTime(time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()))
}

// EndOfDay 当日23:59:59时间（本地时区）
func (it ISOTime) EndOfDay() ISOTime {
	t := it.ToTime().Local()
	return ISOTime(time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location()))
}

// DiffInDays 计算两个时间之间的天数差
func (it ISOTime) DiffInDays(v ISOTime) int64 {
	return int64(it.EndOfDay().ToTime().Sub(v.ToTime()).Hours() / 24)
}

// After .
func (it ISOTime) After(v ISOTime) bool {
	return it.ToTime().After(v.ToTime())
}

// Before .
func (it ISOTime) Before(v ISOTime) bool {
	return it.ToTime().Before(v.ToTime())
}

// Add .
func (it ISOTime) Add(d time.Duration) ISOTime {
	return New(it.ToTime().Add(d))
}

// AddDate 添加年月日
func (it ISOTime) AddDate(year int, month int, day int) ISOTime {
	return New(it.ToTime().AddDate(year, month, day))
}

// LocalDateString YYYY-MM-DD
func (it ISOTime) LocalDateString() string {
	return time.Time(it).Format(dateFormat)
}

// Sub 计算两个时间之间的差值
func (it ISOTime) Sub(v ISOTime) time.Duration {
	return it.ToTime().Sub(v.ToTime())
}

// ISOWeekday 获取 ISO 周几
func (it ISOTime) ISOWeekday() int {
	weekday := it.ToTime().Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	return int(weekday)
}

// SetISOWeekday 设置ISO周几
func (it ISOTime) SetISOWeekday(isoWeekday int) ISOTime {
	// 周日为7
	if isoWeekday == int(time.Sunday) {
		isoWeekday = 7
	}

	currentISOWeekday := it.ISOWeekday()
	diff := isoWeekday - currentISOWeekday

	return it.AddDate(0, 0, diff)
}

// ToDate 转换为ISO日期
func (it ISOTime) ToDate() Date {
	return NewDateFromISOTime(it)
}

// Day .
func (it ISOTime) Day() int {
	return it.ToTime().Day()
}

// SetDay 设置天数
func (it ISOTime) SetDay(day int) ISOTime {
	t := it.ToTime().Local()
	return ISOTime(time.Date(t.Year(), t.Month(), day, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location()))
}
