package nmea0183

import (
	"fmt"
	"github.com/spf13/viper"
	"math"
	"time"
	"strconv"
	"strings"
)

type Config struct {
	Sentences, Variables map[string][]string
	Data                 map[string]string
	History				 map[string]int64
	MessageDate			 time.Time
	UpDated				 time.Time
	realTime			 bool
	autoClearPeriod		 int64               // in milliseconds
}

func setUp() *Config{
	var c Config
	c.Data = make(map[string]string)
	c.History = make(map[string] int64) 
	c.MessageDate = time.Date(0,1,1,0,0,0,0,time.UTC)
	c.UpDated = time.Now().UTC()
	c.realTime = true       // false for historic message processing (or No real time clock) and sentances include a date
	c.autoClearPeriod = 0  // Disabled
	return &c
}

func Create(params ...map[string][]string) *Config {
	c := setUp()
	if len(params) == 1 {
		c.Sentences = params[0]
		c.Variables = *GetDefaultVars()
	}
	if len(params) == 2 {
		c.Sentences = params[0]
		c.Variables = params[1]
	}else{
		c.Sentences = *GetDefaultSentences()
		c.Variables = *GetDefaultVars()
	}
	return c
}

func Load(setting ...string) (*Config, error) {
	c := setUp()
	configSet := []string{".", "nmea_config", "yaml"}
	copy(configSet, setting)

	viper.SetConfigName(configSet[1]) // name of config file (without extension)
	viper.SetConfigType(configSet[2]) // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(configSet[0]) // optionally look for config in the working directory

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err = fmt.Errorf("config file was not found. Use Create or download nmea_config.yaml: %w", err)
			return c, err
		} else {
			// Config file was found but another error was produced
			err = fmt.Errorf("fatal error in config file: %w", err)
			return c, err
		}
	}

	c.Sentences = viper.GetStringMapStringSlice("sentences")
	c.Variables = viper.GetStringMapStringSlice("variables")
	c.Data = make(map[string]string)
	return c, err
}

func SaveConfig(setting ...string){
	configSet := []string{".", "nmea_config", "yaml"}
	copy(configSet, setting)

	viper.SetConfigName(configSet[1]) // name of config file (without extension)
	viper.SetConfigType(configSet[2]) // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(configSet[0]) // optionally look for config in the working directory

	viper.SetDefault("sentences", GetDefaultSentences())
	viper.SetDefault("variables", GetDefaultVars())
	err := viper.ReadInConfig() // Find and read the config file
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		viper.SafeWriteConfig()
	}
}

func (c *Config) DeleteBefore(timeMS int64){
	/*
	Deletes variables in data which have a millisecond time stamp less than timeMS
	*/
	var timeNow int64

	if c.realTime {
		timeNow = time.Now().UTC().UnixMilli()
	}else{
		timeNow = c.MessageDate.UnixMilli()
	}
	timeBefore := timeNow-timeMS

	for i, v := range c.History{
		if v < timeBefore{
				delete(c.Data, i)
		}
	}
}

func (c *Config) Merge(results map[string]string) {
	var timeStamp int64

	if c.autoClearPeriod > 0 { 
		c.DeleteBefore(c.autoClearPeriod)
	}
	c.UpDated = time.Now().UTC()
	
	if c.realTime {
		timeStamp = c.UpDated.UnixMilli()
	}else{
		timeStamp = c.MessageDate.UnixMilli()
	}

	for n, v := range results {
		c.Data[n] = v
		c.History[n] = timeStamp
	}
}

func (c *Config) Preferences(clear int64, realTime bool) {
	if clear < 1 {
		c.autoClearPeriod = 0
	}else{
		c.autoClearPeriod = clear * 1000
	}
	c.realTime = realTime
}

var dateTypeTemplates = []string {"day", "month", "year", "date", "time", "zone"}

func (c *Config) Parse(nmea string) (string, string, error){
	data, preFix, postFix, error := c.ParseToMap(nmea)
	if error == nil{
		c.Merge(data)
	}
	return preFix, postFix, error
}


func (c *Config) ParseToMap(nmea string)  (map[string]string, string, string, error){
	end_byte := len(nmea)
	var err error
	if nmea[end_byte-3] == '*' {
		check_code := checksum(nmea)
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
	key, varList := findInMap(sentenceType, c.Sentences)
	results := make(map[string]string)
	date := ""
	dateTypes := make(map[string]string)

	if len(key) > 0 {
		fieldPointer := 1
		var typeStr string
		for varPointer := 0; varPointer < len(varList); varPointer++ {
			varName, templateList := findInMap(varList[varPointer], c.Variables)
			conVar := ""
			typeStr = ""
			if len(varName) > 0 {
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
						c.MessageDate = messageDate
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
		return fmt.Sprintf("%02d:%02d:%02.2f", h, m, s)
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
	case "x.x":
		return data, "float"
	case "x":
		return data, "int"
	case "tz_h":
		return data, "tzh"
	case "tz_m":
		h, e := strconv.ParseFloat(conVar, 32)
		m, e1 := strconv.ParseFloat(data, 32)
		if e == nil && e1 == nil {
			h += m / 60
			return fmt.Sprintf("%02.2f", h), "tzfloat"
		}
	case "tz:m":
		return conVar + ":" + data, "zone"
	case "plan_tz:m":
		return conVar + ":" + data, "plan_zone"
	
	case "A":
		return data, "A"

	case "str":
		return data, "str"

	case "T":
		return data, "T"

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
		return data + conVar, "xte"

	case "lat":
		d, _ := strconv.ParseInt(data[:2], 10, 32)
		m, _ := strconv.ParseFloat(data[2:], 32)
		return fmt.Sprintf("%02d° %02.4f'", d, m), ""

	case "long":
		d, _ := strconv.ParseInt(data[:3], 10, 32)
		m, _ := strconv.ParseFloat(data[3:], 32)
		return fmt.Sprintf("%03d° %02.4f'", d, m), ""
	
	case "pos_long":
		d, _ := strconv.ParseInt(data[:3], 10, 32)
		m, _ := strconv.ParseFloat(data[3:], 32)
		return conVar + ", " + fmt.Sprintf("%03d° %02.4f'", d, m), ""

	case "lat_WE":
		return conVar + data, "long"
	
	case "pos_WE":
		return conVar + data, "position"

	case "lat_NS":
		return conVar + data, "lat"

	case "ddmmyy":
		date, err := DateSteFromStrs(data[:2], data[2:4], data[4:])
		if err == nil {
			return date, "date"
		}
		return "", ""

	case "plan_ddmmyy":
		date, err := DateSteFromStrs(data[:2], data[2:4], data[4:])
		if err == nil {
			return date, "plan_date"
		}
		return "", ""
}
	return "", ""
}

func DateSteFromStrs(day, month, year string) (string, error){
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

	for i := 1; i < len(s)-3; i++ {
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

func (c *Config) LatLongToFloat(params ...string) (float64, float64, error) {
	/*
	As a method on the handler structure the string parameters refer to variable names in
	c.Data

	Give one parameter for a variable name that holds a position string ie in form lat, Long or
	give the names of 2 variables  which hold lat and long respectively.

	Returns a 2 floats lat and long and an error.  Minus values for South and West
	*/

	if len(params) == 2 {
		lat:= c.Data[params[0]]
		long := c.Data[params[1]]
		return LatLongToFloat(lat, long)
	}
	if len(params) == 1{
		return LatLongToFloat(c.Data[params[0]])
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


func (c *Config) LatLongToString(latFloat, longFloat float64, params ...string) error {
	/*
	As a method on the data structure the string parameters refer to variable names in
	c.Data wich will be set by converting the lat and long float parameters

	One variable name assumes a position string containing lat, long will be set
	two variable names assumes that 2 variables one for lat and on for long will bes set

	Returns an error.  Assumes Minus values are for South and West
	*/

	if len(params) < 1 || len(params) > 2 {
		return fmt.Errorf("illegal number of parmeters given to latlongtostring")
	}

	latStr, longStr, _ := LatLongToString(latFloat, longFloat)
    var timeNow int64

	if c.realTime {
		timeNow = time.Now().UTC().UnixMicro()
	} else {
		timeNow = c.MessageDate.UnixMilli()
	}

	if len(params) == 1 {
		c.Data[params[0]] = latStr + ", " + longStr
		c.History[params[0]] = timeNow
	} else {
		c.Data[params[0]] = latStr
		c.Data[params[1]] = longStr
		c.History[params[0]] = timeNow
		c.History[params[1]] = timeNow
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

	lat := fmt.Sprintf("%02d° %02.4f'%s", latInt, latMins, symbol)

	longAbs := math.Abs(longFloat)
	longInt := int(longAbs)
	longMins := (longAbs - float64(longInt)) * 60
	symbol = "E"
	if longFloat < 0 {
		symbol = "W"
	}

	long := fmt.Sprintf("%03d° %02.4f'%s", longInt, longMins, symbol)

	return lat, long, nil
}
