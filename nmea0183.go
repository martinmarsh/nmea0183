package nmea0183

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type settings struct {
	realTime			 bool
	autoClearPeriod		 int64               // in milliseconds
}


type sentences struct {
	formats, variables map[string][]string
}

type Handle struct {
	data                 map[string]string
	history				 map[string]int64
	messageDate			 time.Time
	upDated				 time.Time
	settings			 settings
	sentences			 *sentences
}

func setUp() *Handle{
	var h Handle
	var set	settings
	
	set.realTime = true       // false for historic message processing (or No real time clock) and sentances include a date
	set.autoClearPeriod = 0  // Disabled

	h.data = make(map[string]string)
	h.history = make(map[string] int64) 
	h.messageDate = time.Date(0,1,1,0,0,0,0,time.UTC)
	h.upDated = time.Now().UTC()
	h.settings = set

	return &h
}

func Create(sentences *sentences) *Handle {
	h := setUp()
	h.sentences = sentences
	return h
}

func Load(setting ...string) (*Handle, error) {
	h := setUp()
	var sent sentences
	configSet := []string{".", "nmea_config", "yaml"}
	copy(configSet, setting)

	viper.SetConfigName(configSet[1]) // name of config file (without extension)
	viper.SetConfigType(configSet[2]) // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(configSet[0]) // optionally look for config in the working directory

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err = fmt.Errorf("sentence file was not found. Use Create or download nmea_config.yaml: %w", err)
			return h, err
		} else {
			// Handle file was found but another error was produced
			err = fmt.Errorf("fatal error in config file: %w", err)
			return h, err
		}
	}

	sent.formats = viper.GetStringMapStringSlice("formats")
	sent.variables = viper.GetStringMapStringSlice("variables")

	h.sentences = &sent
	
	h.data = make(map[string]string)
	return h, err
}

func SaveConfig(setting ...string){
	configSet := []string{".", "nmea_config", "yaml"}
	copy(configSet, setting)

	viper.SetConfigName(configSet[1]) // name of config file (without extension)
	viper.SetConfigType(configSet[2]) // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(configSet[0]) // optionally look for config in the working directory

	viper.SetDefault("formats", GetDefaultFormats())
	viper.SetDefault("variables", GetDefaultVars())
	err := viper.ReadInConfig() // Find and read the config file
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		viper.SafeWriteConfig()
	}
}

func (h *Handle) GetMap() map[string]string{
	return h.data
}

func (h *Handle) Get(key string) string {
	if val, ok := h.data[key]; ok {
		return val
	 } else {
		 return ""
	 }
}

func (h *Handle) DateMap() map[string]time.Time{
	dateMap := make(map[string] time.Time)
	for k, v := range(h.history){
		dateMap[k] = time.UnixMilli(v)
	}
	return dateMap
}

func (h *Handle) Date(key string) time.Time {
	if val, ok := h.history[key]; ok {
		return time.UnixMilli(val)
	 } else {
		 return time.UnixMilli(0)
	 }
}

func (h *Handle) DeleteBefore(timeMS int64){
	/*
	Deletes variables in data which have a millisecond time stamp less than timeMS
	*/
	var timeNow int64

	if h.settings.realTime {
		timeNow = time.Now().UTC().UnixMilli()
	}else{
		timeNow = h.messageDate.UnixMilli()
	}
	timeBefore := timeNow-timeMS

	for i, v := range h.history{
		if v < timeBefore{
				delete(h.data, i)
		}
	}
}

func (h *Handle) Merge(results map[string]string) {
	var timeStamp int64

	if h.settings.autoClearPeriod > 0 { 
		h.DeleteBefore(h.settings.autoClearPeriod)
	}
	h.upDated = time.Now().UTC()
	
	if h.settings.realTime {
		timeStamp = h.upDated.UnixMilli()
	}else{
		timeStamp = h.messageDate.UnixMilli()
	}

	for n, v := range results {
		h.data[n] = v
		h.history[n] = timeStamp
	}
}

func (h *Handle) Preferences(clear int64, realTime bool) {
	if clear < 1 {
		h.settings.autoClearPeriod = 0
	}else{
		h.settings.autoClearPeriod = clear * 1000
	}
	h.settings.realTime = realTime
}

var dateTypeTemplates = []string {"day", "month", "year", "date", "time", "zone"}

func (h *Handle) Parse(nmea string) (string, string, error){
	data, preFix, postFix, error := h.ParseToMap(nmea)
	if error == nil{
		h.Merge(data)
	}
	return preFix, postFix, error
}


func (h *Handle) ParseToMap(nmea string)  (map[string]string, string, string, error){
	end_byte := len(nmea)
	var err error
	if nmea[end_byte-3] == '*' {
		check_code := checksum(nmea[:end_byte-3])
		end_byte -= 2
		if check_code != nmea[end_byte:] {
			err_mess := fmt.Sprintf("error: %s != %s", check_code, nmea[end_byte:])
			err = fmt.Errorf("check sum error: %s", err_mess)
		}
		end_byte--
	}

	parts := strings.Split(nmea[1:end_byte], ",")
	preFix := parts[0][:2]
	sentenceType := strings.ToLower(parts[0][2:])
	key, varList := findInMap(sentenceType, h.sentences.formats)
	results := make(map[string]string)
	date := ""
	dateTypes := make(map[string]string)

	if len(key) > 0 {
		fieldPointer := 1
		var typeStr string
		for varPointer := 0; varPointer < len(varList); varPointer++ {
			varName, templateList := findInMap(varList[varPointer], h.sentences.variables)
			conVar := ""
			typeStr = ""
			if len(varName) > 0 && fieldPointer < len(parts){
				for _, template := range templateList {
					conVar, typeStr = convert(parts[fieldPointer], template, conVar)
					fieldPointer++
				}
				for _, v := range dateTypeTemplates {
					if v == typeStr {
						dateTypes[v] = conVar
					}
				}
				results[varName] = conVar
			} else {
				fieldPointer++
			}
		}
		if len(dateTypes) > 1 {
			if date_found, ok := dateTypes["date"]; ok {
				date = date_found 
			}
			if day, ok := dateTypes["day"]; ok {
				if month, ok := dateTypes["month"]; ok 	{
					if year, ok := dateTypes["year"]; ok {
						d, _ := strconv.ParseInt(day, 10, 32)
						m, _ := strconv.ParseInt(month, 10, 32)
						y, _ := strconv.ParseInt(year, 10, 32)
						if y < 60 {
							y += 2000
						} else {
							y += 1900
						}
						date = fmt.Sprintf("%d-%02d-%02d", y, m, d)
					}
				}
			}
			if len(date) > 0 {
				if timeT, ok := dateTypes["time"]; ok {
					dateTime := date + "T" + timeT
					zone := "00:00"
					if z, ok := dateTypes["zone"]; ok {
						zone = z
					}
					rcDate := dateTime + "+" + zone	
					messageDate, err := time.Parse(time.RFC3339, rcDate)
					if err == nil{
						h.messageDate = messageDate
						results["datetime"] = rcDate

					}
				}	
			}
		}
	}
	return results,  preFix, sentenceType, err
}

func timeConv(data string) string{
	h, e := strconv.ParseInt(data[:2], 10, 16)
	m, e1 := strconv.ParseInt(data[2:4], 10, 16)
	s, e2 := strconv.ParseFloat(data[4:], 32)
	if e == nil && e1 == nil && e2 == nil {
		return fmt.Sprintf("%02d:%02d:%05.2f", h, m, s)
	}
	return ""
}

func convert(data string, template string, conVar string) (string, string) {
	switch template {
		case "hhmmss.ss":
			t := timeConv(data)
			if len(t) > 0{
				return conVar + t, "time"
			}
			return "", ""

		case "plan_hhmmss.ss":
			t := timeConv(data)
			if len(t) > 0{
				return conVar + t, "plan_time"
			}
			return "", ""
		case "DD_day":
			return data, "day"
		case "DD_month":
			return data, "month"
		case "DD_year":
			return data, "plan_year"
		case "DD_day_plan":
			return data, "plan_day"
		case "DD_month_plan":
			return data, "plan_month"
		case "DD_year_plan":
			return data, "plan_year"
		case "x.x", "-x.x", "Lx.xN", "x.xT":
			return data, "float"
		case "x":
			return data, "int"
		case "tz_h":
			return data, "tzh"
		
		case "tz:m":
			return conVar + ":" + data, "zone"
		case "plan_tz:m":
			return conVar + ":" + data, "plan_zone"
		
		case "A":
			return data, "A"

		case "str":
			return data, "str"

		case "T":
			return conVar + "°"+ data, "T"

		case "w":
			if data == "W" || data == "w" {
				return "-" + conVar, "float"
			}
			return conVar, "west"

		case "s":
			if data == "S" || data == "s" {
				return "-" + conVar, "float"
			}
			return conVar, "south"

		case "R":
			return data + conVar, ""

		case "N":
			return conVar + data, "xte"

		case "lat":
			d, _ := strconv.ParseInt(data[:2], 10, 32)
			m, _ := strconv.ParseFloat(data[2:], 32)
			return fmt.Sprintf("%02d° %07.4f'", d, m), ""

		case "long":
			d, _ := strconv.ParseInt(data[:3], 10, 32)
			m, _ := strconv.ParseFloat(data[3:], 32)
			return fmt.Sprintf("%03d° %07.4f'", d, m), ""
		
		case "pos_long":
			d, _ := strconv.ParseInt(data[:3], 10, 32)
			m, _ := strconv.ParseFloat(data[3:], 32)
			return conVar + ", " + fmt.Sprintf("%03d° %07.4f'", d, m), ""

		case "long_WE":
			return conVar + data, "long"
		
		case "pos_WE":
			return conVar + data, "position"

		case "lat_NS":
			return conVar + data, "lat"

		case "ddmmyy":
			date, err := DateStrFromStrs(data[:2], data[2:4], data[4:])
			if err == nil {
				return date, "date"
			}
			return "", ""

		case "plan_ddmmyy":
			date, err := DateStrFromStrs(data[:2], data[2:4], data[4:])
			if err == nil {
				return date, "plan_date"
			}
			return "", ""
	}
	return "", ""
}

func DateStrFromStrs(day, month, year string) (string, error){
	var err error
	err = nil
	d, e1 := strconv.ParseInt(day, 10, 32)
	m, e2 := strconv.ParseInt(month, 10, 32)
	y, e3 := strconv.ParseInt(year, 10, 32)
	if e1 != nil || e2 != nil || e3 != nil{
		err = fmt.Errorf("conversion error")
	}
	if y < 60 {
		y += 2000
	} else {
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

func findInMap(k string, m map[string][]string) (string, []string) {

	for i, v := range m {
		if i == k {
			return i, v
		}
	}
	return "", []string{""}
}

func (h *Handle) LatLongToFloat(params ...string) (float64, float64, error) {
	/*
	As a method on the handler structure the string parameters refer to variable names in
	c.data

	Give one parameter for a variable name that holds a position string ie in form lat, Long or
	give the names of 2 variables  which hold lat and long respectively.

	Returns a 2 floats lat and long and an error.  Minus values for South and West
	*/

	if len(params) == 2 {
		lat:= h.data[params[0]]
		long := h.data[params[1]]
		return LatLongToFloat(lat, long)
	}
	if len(params) == 1{
		return LatLongToFloat(h.data[params[0]])
	}

	return 0, 0, fmt.Errorf("illegal number of parmeters given to latlongtofloat")
}

func LatLongToFloat(params ...string) (float64, float64, error) {
	/*
	Give one parameter for a position string with lat, Long or
	2 parameters to give separate lat and long  strings.

	Returns a 2 floats lat and long and an error.  Minus values for South and West
	*/
	if len(params)<1 || len(params) > 2 {
		return 0,0, fmt.Errorf("illegal number of parmeters given to latlongtofloat")
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
		lenL --
		symbol := lat[lenL]
		lenL --
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
		lenL -- 
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


func (h *Handle) LatLongToString(latFloat, longFloat float64, params ...string) error {
	/*
	As a method on the data structure the string parameters refer to variable names in
	c.data wich will be set by converting the lat and long float parameters

	One variable name assumes a position string containing lat, long will be set
	two variable names assumes that 2 variables one for lat and on for long will bes set

	Returns an error.  Assumes Minus values are for South and West
	*/

	if len(params) < 1 || len(params) > 2 {
		return fmt.Errorf("illegal number of parmeters given to latlongtostring")
	}

	latStr, longStr, _ := LatLongToString(latFloat, longFloat)
    var timeNow int64

	if h.settings.realTime {
		timeNow = time.Now().UTC().UnixMilli()
	} else {
		timeNow = h.messageDate.UnixMilli()
	}

	if len(params) == 1 {
		h.data[params[0]] = latStr + ", " + longStr
		h.history[params[0]] = timeNow
	} else {
		h.data[params[0]] = latStr
		h.data[params[1]] = longStr
		h.history[params[0]] = timeNow
		h.history[params[1]] = timeNow
	}
	
	return nil
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


func (h *Handle) WriteSentence(manCode string, sentenceName string) (string, error) {
	sentenceType := strings.ToLower(sentenceName)
	key, varList := findInMap(sentenceType, h.sentences.formats)
	madeSentence := strings.ToUpper(manCode + sentenceName)

	if len(key) > 1 {
		for _, v := range varList {
			_, vFormats := findInMap(v, h.sentences.variables)
			value, ok := h.data[v]
			if ok && len(value)>0 {
				if len(v) == 0 || v == "n/a" {
					madeSentence += ","
				}else{
					m, err := getSentencePart(value, vFormats)
					if err != nil{
						return "", fmt.Errorf("field error definition %w", err)
					}
					madeSentence += m 
				}
			}else{
				for i:=0; i<len(vFormats); i++{
					madeSentence += ","
				}
			}
		}
		madeSentence = "$" + madeSentence
		checksum := checksum(madeSentence)
		madeSentence += "*" + checksum
		return madeSentence, nil
	}
	return "", fmt.Errorf("no matching sentence definition")
}

func getSentencePart(data string, varItems []string) (string, error){
	subString := ""
	lookForward := ""
	for _, template := range varItems{
		value, _ :=convertTo183(data, template, lookForward)
		subString += "," + value
	}

	return subString, nil
}


func convertTo183(data, template string, forward string) (string, string) {
	switch template {
		case "hhmmss.ss", "plan_hhmmss.ss":
			return data[:2]+data[3:5]+data[6:], ""
	
		case "DD_day":
			return data, "day"
		case "DD_month":
			return data, "month"
		case "DD_year":
			return data, "plan_year"
		case "DD_day_plan":
			return data, "plan_day"
		case "DD_month_plan":
			return data, "plan_month"
		case "DD_year_plan":
			return data, "plan_year"
		case "x.x":	
			if data[0] >= '0' && data[0] <= '9'{
				return data, "float"
			}
			return data[1:], "float"
		case "-x.x":
			return data, "float"
		case "Lx.xN":
			l := len(data)
			if l > 2 {
				return data[1:l-1], "LN"
			}
			return "", ""
		case "x.xT":
			l := len(data)
			if l > 2 {
				return data[:l-3], "xT"
			}
			return "", ""

		case "x":
			return data, "int"
		case "tz_h":
			return data[:2], "tzh"

		case "tz:m":
			return data[3:], "zone"

		case "plan_tz:m":
			return data[3:], "plan_zone"
		
		case "A":
			return data, "A"
		case "w":
			if data[0] == '-' {
				return "W", "west"
			}
			return "E", "west"
		case "s":
			if data[0] == '-' {
				return "S", "south"
			}
			return "N", "south"
		
		case "R":
			if data[0] == 'R' {
				return "R", "LR"
			}
			return "L", "LR"

		case "N":
			l := len(data)-1
			return string(data[l]), "N"
			
		case "str":
			return data, "str"

		case "T":
			l := len(data)-1
			return string(data[l]), "N"

		case "lat":
			split :=  strings.SplitN(data[4:], ",", 2)
			l := len(split[0]) - 2
			return data[:2] + split[0][1:l], ""

		case "long":
			l := len(data) - 2
			return data[:3] + data[5:l], ""

		case "pos_long":
			split :=  strings.SplitN(data[5:], ",", 2)
			l := len(split[1]) - 2
			return split[1][1:4] + split[1][7:l], ""

		case "long_WE":
			l := len(data)-1
			return string(data[l]), ""
	
		case "pos_WE":
			split :=  strings.SplitN(data[5:], ",", 2)
			l := len(split[1])-1
			return string(split[1][l]), ""

		case "lat_NS":
			split :=  strings.SplitN(data[4:], ",", 2)
			l := len(split[0])-1
			return string(split[0][l]), ""
			
		case "ddmmyy", "plan_ddmmyy":	
			return data[8:] + data[5:7] + data[2:4], ""
	}
	return "", ""
}
