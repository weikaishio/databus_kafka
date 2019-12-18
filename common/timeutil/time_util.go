package timeutil

import (
	"strconv"
	"strings"
	"time"
)

const (
	DATE_FORMAT_YYYYMMDD  = "2006-01-02"
	DATE_FORMAT_YYYYMMDD2 = "20060102"
	DATE_FORMAT_YYYYMM    = "2006-01"
	DATE_FORMAT_YYYYMM2   = "200601"
	DATE_FORMAT_YYYY      = "2006"
	SEP_STR               = "-"
)

var (
	startTime    int64 = 0
	endTime      int64 = 0
	startTimeStr       = ""
	endTimeStr         = ""
)

//将时间戳按照时区转换成字符串
func TimestampToInt(timestamp int64, convType int) (int64, error) {
	//将十三位的时间戳时间截取到十位
	timeStr := strconv.FormatInt(timestamp, 10)
	newTimeStr := timeStr[0:10]
	newTimeStamp, _ := strconv.ParseInt(newTimeStr, 10, 64)
	tm := time.Unix(newTimeStamp, 0)
	switch convType {
	case 1: //年月日
		dateStr := tm.Format(DATE_FORMAT_YYYYMMDD)
		dateArr := strings.Split(dateStr, SEP_STR)
		dateInt := strings.Join(dateArr, "")
		return strconv.ParseInt(dateInt, 10, 64)
	case 2: //年和周
		year, week := tm.ISOWeek()
		dateInt := strconv.Itoa(year) + strconv.Itoa(week)
		return strconv.ParseInt(dateInt, 10, 64)
	case 3: //年和月
		dateStr := tm.Format(DATE_FORMAT_YYYYMM)
		dateArr := strings.Split(dateStr, SEP_STR)
		dateInt := strings.Join(dateArr, "")
		return strconv.ParseInt(dateInt, 10, 64)
	case 4: //年
		dateStr := tm.Format(DATE_FORMAT_YYYY)
		dateArr := strings.Split(dateStr, SEP_STR)
		dateInt := strings.Join(dateArr, "")
		return strconv.ParseInt(dateInt, 10, 64)
	default:
		return 0, nil
	}
	return 0, nil
}

//获取指定周的开始时间和结束时间
func GetWeekStartAndEnd(timestamp int64) (int64, int64) {
	//将十三位的时间戳时间截取到十位
	timeStr := strconv.FormatInt(timestamp, 10)
	newTimeStr := timeStr[0:10]
	newTimeStamp, _ := strconv.ParseInt(newTimeStr, 10, 64)
	tm := time.Unix(newTimeStamp, 0)
	whatDay := strings.ToLower(tm.Weekday().String())
	switch whatDay {
	case "sunday":
		startTimeStr = tm.Format(DATE_FORMAT_YYYYMMDD)
		endTimeStr = tm.AddDate(0, 0, +6).Format(DATE_FORMAT_YYYYMMDD)

	case "monday":
		startTimeStr = tm.AddDate(0, 0, -1).Format(DATE_FORMAT_YYYYMMDD)
		endTimeStr = tm.AddDate(0, 0, +5).Format(DATE_FORMAT_YYYYMMDD)

	case "tuesday":
		startTimeStr = tm.AddDate(0, 0, -2).Format(DATE_FORMAT_YYYYMMDD)
		endTimeStr = tm.AddDate(0, 0, +4).Format(DATE_FORMAT_YYYYMMDD)

	case "wednesday":
		startTimeStr = tm.AddDate(0, 0, -3).Format(DATE_FORMAT_YYYYMMDD)
		endTimeStr = tm.AddDate(0, 0, +3).Format(DATE_FORMAT_YYYYMMDD)

	case "thursday":
		startTimeStr = tm.AddDate(0, 0, -4).Format(DATE_FORMAT_YYYYMMDD)
		endTimeStr = tm.AddDate(0, 0, +2).Format(DATE_FORMAT_YYYYMMDD)

	case "friday":
		startTimeStr = tm.AddDate(0, 0, -5).Format(DATE_FORMAT_YYYYMMDD)
		endTimeStr = tm.AddDate(0, 0, +1).Format(DATE_FORMAT_YYYYMMDD)
	case "saturday":
		startTimeStr = tm.AddDate(0, 0, -6).Format(DATE_FORMAT_YYYYMMDD)
		endTimeStr = tm.Format(DATE_FORMAT_YYYYMMDD)

	default:

	}
	startDateArr := strings.Split(startTimeStr, SEP_STR)
	startDateInt := strings.Join(startDateArr, "")
	startTime, _ = strconv.ParseInt(startDateInt, 10, 64)

	endDateArr := strings.Split(endTimeStr, SEP_STR)
	endDateInt := strings.Join(endDateArr, "")
	endTime, _ = strconv.ParseInt(endDateInt, 10, 64)

	return startTime, endTime
}

//获取指定月的开始时间和结束时间
func GetMonthStartAndEnd(timestamp int64) (int64, int64) {
	//将十三位的时间戳时间截取到十位
	timeStr := strconv.FormatInt(timestamp, 10)
	newTimeStr := timeStr[0:10]
	newTimeStamp, _ := strconv.ParseInt(newTimeStr, 10, 64)
	year, month, _ := time.Unix(newTimeStamp, 0).Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	start := thisMonth.AddDate(0, 0, 0).Format(DATE_FORMAT_YYYYMMDD)
	end := thisMonth.AddDate(0, 1, -1).Format(DATE_FORMAT_YYYYMMDD)

	startDateArr := strings.Split(start, SEP_STR)
	startDateInt := strings.Join(startDateArr, "")
	startTime, _ = strconv.ParseInt(startDateInt, 10, 64)

	endDateArr := strings.Split(end, SEP_STR)
	endDateInt := strings.Join(endDateArr, "")
	endTime, _ = strconv.ParseInt(endDateInt, 10, 64)

	return startTime, endTime
}

//获取指定年的开始时间和结束时间
func GetYearStartAndEnd(timestamp int64) (int64, int64) {
	//将十三位的时间戳时间截取到十位
	timeStr := strconv.FormatInt(timestamp, 10)
	newTimeStr := timeStr[0:10]
	newTimeStamp, _ := strconv.ParseInt(newTimeStr, 10, 64)
	year, _, _ := time.Unix(newTimeStamp, 0).Date()
	thisMonth := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
	start := thisMonth.AddDate(0, 0, 0).Format(DATE_FORMAT_YYYYMMDD)
	end := thisMonth.AddDate(1, 0, -1).Format(DATE_FORMAT_YYYYMMDD)

	startDateArr := strings.Split(start, SEP_STR)
	startDateInt := strings.Join(startDateArr, "")
	startTime, _ = strconv.ParseInt(startDateInt, 10, 64)

	endDateArr := strings.Split(end, SEP_STR)
	endDateInt := strings.Join(endDateArr, "")
	endTime, _ = strconv.ParseInt(endDateInt, 10, 64)

	return startTime, endTime
}

func IntToStartAndEnd(currentTime int64, convType int) (int64, int64) {
	switch convType {
	case 1:
		//带天的数据直接使用
	case 2:
		return IntToWeekSartAndEnd(currentTime)
	case 3:
		currentTimeStr := strconv.FormatInt(currentTime, 10)
		startTime, _ := time.Parse(DATE_FORMAT_YYYYMM2, currentTimeStr)
		return GetMonthStartAndEnd(startTime.Unix())
	case 4:
		currentTimeStr := strconv.FormatInt(currentTime, 10)
		startTime, _ := time.Parse(DATE_FORMAT_YYYY, currentTimeStr)
		return GetYearStartAndEnd(startTime.Unix())
	default:

	}
	return startTime, endTime
}

//得到指定周的开始时间
func FirstDayOfISOWeek(year int, week int) (int64, int64) {
	startDate := time.Date(year, 12, 31, 0, 0, 0, 0, time.Local)
	isoYear, isoWeek := startDate.ISOWeek()
	//fmt.Printf("startDate:%v, isoYear:%d, isoWeek:%d\n", startDate, isoYear, isoWeek)

	for isoYear < year { // iterate forward to the first day of the first week
		startDate = startDate.AddDate(0, 0, 1)
		isoYear, isoWeek = startDate.ISOWeek()
		//fmt.Printf("222 isoYear:%d, isoWeek:%d\n", isoYear, isoWeek)
	}
	for week < isoWeek && isoYear == year { // iterate forward to the first day of the given week
		startDate = startDate.AddDate(0, 0, -1)
		isoYear, isoWeek = startDate.ISOWeek()
		//fmt.Printf("333 isoYear:%d, isoWeek:%d\n", isoYear, isoWeek)
	}
	for startDate.Weekday() != time.Monday { // iterate back to Sunday
		startDate = startDate.AddDate(0, 0, -1)
		isoYear, isoWeek = startDate.ISOWeek()
		//fmt.Printf("111 isoYear:%d, isoWeek:%d\n", isoYear, isoWeek)
	}
	endDate := startDate.AddDate(0, 0, 6)
	startTimeStr = startDate.Format(DATE_FORMAT_YYYYMMDD2)
	endTimeStr = endDate.Format(DATE_FORMAT_YYYYMMDD2)
	startTime, _ = strconv.ParseInt(startTimeStr, 10, 64)
	endTime, _ = strconv.ParseInt(endTimeStr, 10, 64)
	return startTime, endTime
}

func IntToWeekSartAndEnd(timeInt int64) (int64, int64) {
	//先将整型转成字符串,然后截取最后两位
	timeStr := strconv.FormatInt(timeInt, 10)
	yearStr := timeStr[0:4]
	weekStr := timeStr[4:]
	yearInt, _ := strconv.ParseInt(yearStr, 10, 64)
	weekInt, _ := strconv.ParseInt(weekStr, 10, 64)
	//fmt.Printf("y:%d,w:%d\n", yearInt, weekInt)
	return FirstDayOfISOWeek(int(yearInt), int(weekInt))
}
