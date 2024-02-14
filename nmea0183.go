package nmea0183

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

type varTypeStruct struct {
	fType string
	fConv varFormatStruct
}

type varFormatStruct struct {
	fCount int
	from   func(int, *[]string) string
	to     func(string) string
}

func getConversion(format string) (string, varFormatStruct) {
	// given a format string returns the type and conversion structure
	conv := *makeConv()
	if varConv, ok := conv[format]; ok {
		return varConv.fType, varConv.fConv
	} else {
		return "", varFormatStruct{
			fCount: 1, from: func(pos int, parts *[]string) string { return "" },
			to: func(data string) string { return "" },
		}
	}
}

func makeConv() *map[string]varTypeStruct {
	timeConv := varFormatStruct{
		fCount: 1,
		from:   func(pos int, parts *[]string) string { return timeFormat((*parts)[pos]) },
		to: func(data string) string {
			return data[:2] + data[3:5] + data[6:]
		},
	}

	copyField := varFormatStruct{
		fCount: 1,
		from:   func(pos int, parts *[]string) string { return (*parts)[pos] },
		to: func(data string) string {
			return data
		},
	}

	compass := varFormatStruct{
		fCount: 2,
		from: func(pos int, parts *[]string) string {
			return (*parts)[pos] + "°" + (*parts)[pos+1]
		},
		to: func(data string) string {
			l := len(data)
			if l > 2 {
				return data[:l-3] + "," + string(data[l-1])
			}
			return ","
		},
	}

	hrsMins := varFormatStruct{
		fCount: 2,
		from: func(pos int, parts *[]string) string {
			hr := (*parts)[pos]
			mins := (*parts)[pos+1]
			sign := "+"
			if len(hr) > 1 && hr[0] == '-' {
				sign = "-"
				hr = hr[1:]
			}
			if len(mins) > 1 && mins[0] == '-' {
				sign = "-"
				mins = mins[1:]
			}
			return sign + hr + ":" + mins
		},
		to: func(data string) string {
			l := len(data)
			if l > 4 {
				if data[0] == '+' {
					return data[1:3] + "," + data[4:]
				}
				mins := data[4:]
				if mins != "00" {
					mins = "-" + mins
				}
				return data[:3] + "," + mins
			}
			return ","
		},
	}

	deviation := varFormatStruct{
		fCount: 2,
		from: func(pos int, parts *[]string) string {
			data := (*parts)[pos]
			data2 := (*parts)[pos+1]
			if data2 == "W" || data2 == "w" {
				return "-" + data
			}
			return data
		},
		to: func(data string) string {
			if data[0] == '-' {
				return data[1:] + ",W"
			}
			return data + ",E"
		},
	}

	xte := varFormatStruct{
		fCount: 3,
		from: func(pos int, parts *[]string) string {
			return (*parts)[pos+1] + (*parts)[pos] + (*parts)[pos+2]
		},
		to: func(data string) string {
			l := len(data)
			if l > 2 {
				val := data[1 : l-1]
				return val + "," + string(data[0]) + "," + string(data[l-1])
			}
			return ",,"
		},
	}

	lat := varFormatStruct{
		fCount: 2,
		from: func(pos int, parts *[]string) string {
			return latStr((*parts)[pos]) + (*parts)[pos+1]
		},
		to: func(data string) string {
			return latVar(data) + "," + latNVar(data)
		},
	}

	long := varFormatStruct{
		fCount: 2,
		from: func(pos int, parts *[]string) string {
			return longStr((*parts)[pos]) + (*parts)[pos+1]
		},
		to: func(data string) string {
			return longVar(data) + "," + longWeVar(data)
		},
	}

	position := varFormatStruct{
		fCount: 4,
		from: func(pos int, parts *[]string) string {
			return latStr((*parts)[pos]) + (*parts)[pos+1] + ", " + longStr((*parts)[pos+2]) + (*parts)[pos+3]
		},
		to: func(data string) string {
			return latVar(data) + "," + latNVar(data) + "," + posLongVar(data) + "," + posWeVar(data)
		},
	}

	date := varFormatStruct{
		fCount: 1,
		from: func(pos int, parts *[]string) string {
			data := (*parts)[pos]
			date, err := DateStrFromStrs(data[:2], data[2:4], data[4:])
			if err == nil {
				return date
			}
			return ""
		},
		to: func(data string) string {
			return data[8:] + data[5:7] + data[2:4]
		},
	}

	dateTime := varFormatStruct{
		fCount: 6,
		from: func(pos int, parts *[]string) string {
			timeofday := timeFormat((*parts)[pos])
			day := (*parts)[pos+1]
			month := (*parts)[pos+2]
			year := (*parts)[pos+3]
			tz := hrsMins.from(pos+4, parts)
			date, errd := DateStrFromStrs(day, month, year)
			rcDate := date + "T" + timeofday + tz

			if errd == nil {
				return rcDate
			}
			return ""
		},
		to: func(data string) string {
			return dateTimeToCSV(data)
		},
	}

	varConv := map[string]varTypeStruct{
		"hhmmss.ss":                     {fType: "time", fConv: timeConv},
		"plan_hhmmss.ss":                {fType: "plan time", fConv: timeConv},
		"A":                             {fType: "status", fConv: copyField},
		"c--c":                          {fType: "string", fConv: copyField},
		"x.x":                           {fType: "float", fConv: copyField},
		"-x.x":                          {fType: "signed float", fConv: copyField},
		"x":                             {fType: "integer", fConv: copyField},
		"-x":                            {fType: "signed integer", fConv: copyField},
		"xxx,T":                         {fType: "compass", fConv: compass},
		"x.x,T":                         {fType: "compass", fConv: compass},
		"T":                             {fType: "magnetic", fConv: copyField},
		"x.x,R,N":                       {fType: "cross track error", fConv: xte},
		"lat,NS":                        {fType: "lat", fConv: lat},
		"long,WE":                       {fType: "long", fConv: long},
		"lat,NS,long,WE":                {fType: "position", fConv: position},
		"ddmmyy":                        {fType: "date", fConv: date},
		"plan_ddmmyy":                   {fType: "plan date", fConv: date},
		"x.x,w":                         {fType: "deviation", fConv: deviation},
		"DD_day":                        {fType: "day", fConv: copyField},
		"DD_month":                      {fType: "month", fConv: copyField},
		"DD_year":                       {fType: "year", fConv: copyField},
		"tz_h,tz_m":                     {fType: "zone", fConv: hrsMins},
		"plan_DD_day":                   {fType: "plan_day", fConv: copyField},
		"plan_DD_month":                 {fType: "plan_month", fConv: copyField},
		"plan_DD_year":                  {fType: "plan_year", fConv: copyField},
		"plan_tz_h,tz_m":                {fType: "plan_zone", fConv: hrsMins},
		"hhmmss,day,month,year,tz":      {fType: "datetime", fConv: dateTime},
		"plan_hhmmss,day,month,year,tz": {fType: "plan_datetime", fConv: dateTime},
	}

	return &varConv
}

//var dateTypeTemplates = []string {"day", "month", "year", "date", "time", "zone"}

func latStr(data string) string {
	d, _ := strconv.ParseInt(data[:2], 10, 32)
	m, _ := strconv.ParseFloat(data[2:], 32)
	return fmt.Sprintf("%02d° %07.4f'", d, m)
}

func longStr(data string) string {
	d, _ := strconv.ParseInt(data[:3], 10, 32)
	m, _ := strconv.ParseFloat(data[3:], 32)
	return fmt.Sprintf("%03d° %07.4f'", d, m)
}

func latVar(data string) string {
	split := strings.SplitN(data[4:], ",", 2)
	l := len(split[0]) - 2
	return data[:2] + split[0][1:l]
}

func longVar(data string) string {
	l := len(data) - 2
	return data[:3] + data[5:l]
}

func posLongVar(data string) string {
	split := strings.SplitN(data[5:], ",", 2)
	l := len(split[1]) - 2
	return split[1][1:4] + split[1][7:l]
}

func longWeVar(data string) string {
	l := len(data) - 1
	return string(data[l])
}

func posWeVar(data string) string {
	split := strings.SplitN(data[5:], ",", 2)
	l := len(split[1]) - 1
	return string(split[1][l])
}

func latNVar(data string) string {
	split := strings.SplitN(data[4:], ",", 2)
	l := len(split[0]) - 1
	return string(split[0][l])
}

func dateTimeToCSV(data string) string {
	// give RFC3339 formated string
	// return time, day, month, year, zonehrs, zonemins as comma separated strings
	timeParts := strings.Split(data, "T")
	d := timeParts[0]
	timeZ := strings.Split(timeParts[1], "+")
	sign := ""
	signMins := ""
	if len(timeZ) == 1 {
		timeZ = strings.Split(timeParts[1], "-")
		sign = "-"
	}
	t := timeZ[0]
	zone := strings.Split(timeZ[1], ":")
	signMins = sign
	if zone[0] == "00" {
		signMins = ""
	}
	return t[:2] + t[3:5] + t[6:] + "," + d[8:] + "," + d[5:7] + "," + d[:4] + "," + sign + zone[0] + "," + signMins + zone[1]

}

func setUp() *Handle {
	var h Handle
	var set settings

	set.realTime = true     // false for historic message processing (or No real time clock) and sentences include a date
	set.autoClearPeriod = 0 // Disabled

	h.data = make(map[string]string)
	h.history = make(map[string]int64)
	h.messageDate = time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)
	h.upDated = time.Now().UTC()
	h.settings = set

	return &h
}

func timeFormat(data string) string {
	h, e := strconv.ParseInt(data[:2], 10, 16)
	m, e1 := strconv.ParseInt(data[2:4], 10, 16)
	s, e2 := strconv.ParseFloat(data[4:], 32)
	if e == nil && e1 == nil && e2 == nil {
		return fmt.Sprintf("%02d:%02d:%05.2f", h, m, s)
	}
	return ""
}

func DateStrFromStrs(day, month, year string) (string, error) {
	var err error
	err = nil
	d, e1 := strconv.ParseInt(day, 10, 32)
	m, e2 := strconv.ParseInt(month, 10, 32)
	y, e3 := strconv.ParseInt(year, 10, 32)
	if e1 != nil || e2 != nil || e3 != nil {
		err = fmt.Errorf("conversion error")
	}
	if y < 60 {
		y += 2000
	} else if y < 100 {
		y += 1900
	}
	return fmt.Sprintf("%d-%02d-%02d", y, m, d), err
}

func checksum(s string) string {
	check_sum := 0

	nmea_data := []byte(s)

	for i := 1; i < len(s); i++ {
		check_sum ^= (int)(nmea_data[i])
	}

	return fmt.Sprintf("%02X", check_sum)
}

func LatLongToFloat(params ...string) (float64, float64, error) {
	/*
		Give one parameter for a position string with lat, Long or
		2 parameters to give separate lat and long  strings.

		Returns a 2 floats lat and long and an error.  Minus values for South and West
	*/
	if len(params) < 1 || len(params) > 2 {
		return 0, 0, fmt.Errorf("illegal number of parmeters given to latlongtofloat")
	}
	var lat, long string
	var symbol byte
	var retLat, retLong float64 = 0, 0

	if len(params) == 1 {
		params = strings.SplitN(params[0], ", ", 2)
	}
	if len(params) == 2 {
		lat = params[0]
		long = params[1]
	}

	lenL := len(lat)
	if lenL > 8 {
		lenL--
		symbol := lat[lenL]
		lenL--
		dlat, _ := strconv.ParseFloat(lat[:2], 64)
		mlat, _ := strconv.ParseFloat(lat[5:lenL], 64)
		dlat += mlat / 60
		if symbol == 'S' {
			dlat = -dlat
		}
		retLat = dlat
	}
	lenL = len(long)
	if lenL > 8 {
		lenL--
		symbol = long[lenL]
		lenL--
		dlong, _ := strconv.ParseFloat(long[:3], 64)
		mlong, _ := strconv.ParseFloat(long[6:lenL], 64)
		dlong += mlong / 60
		if symbol == 'W' {
			dlong = -dlong
		}
		retLong = dlong
	}

	return retLat, retLong, nil
}

func LatLongToString(latFloat, longFloat float64) (string, string, error) {
	/*
		Give  2 variables lat and long respectively. Minus values given denote South and West
		Returns a 2 fromatted strings as lat and long and an error.
	*/

	latAbs := math.Abs(latFloat)
	latInt := int(latAbs)
	latMins := (latAbs - float64(latInt)) * 60
	symbol := "N"
	if latFloat < 0 {
		symbol = "S"
	}

	lat := fmt.Sprintf("%02d° %07.4f'%s", latInt, latMins, symbol)

	longAbs := math.Abs(longFloat)
	longInt := int(longAbs)
	longMins := (longAbs - float64(longInt)) * 60
	symbol = "E"
	if longFloat < 0 {
		symbol = "W"
	}

	long := fmt.Sprintf("%03d° %07.4f'%s", longInt, longMins, symbol)

	return lat, long, nil
}
