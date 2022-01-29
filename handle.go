package nmea0183

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)


type settings struct {
	realTime			 bool
	autoClearPeriod		 int64               // in milliseconds
}


type Handle struct {
	data                 map[string]string
	history				 map[string]int64
	messageDate			 time.Time
	upDated				 time.Time
	settings			 settings
	sentences			 *Sentences
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

	results := make(map[string]string)
	date := ""
	dateTypes := make(map[string]string)

	if varList, found := h.sentences.formats[sentenceType]; found{
		fieldPointer := 1
		for varPointer := 0; varPointer < len(varList); varPointer++ {
			conVar := ""
			varName := varList[varPointer]
			if template, found := h.sentences.variables[varName]; found && fieldPointer < len(parts){
				fType, cv := getConversion(template)
				if fieldPointer + cv.fCount <= len(parts){
						conVar = cv.from(fieldPointer, &parts)
						dateTypes[fType] = varName
				} else {
						conVar = ""
				}
				fieldPointer += cv.fCount
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


					}
				}	
			}
		}
	}
	return results,  preFix, sentenceType, err
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


func (h *Handle) WriteSentence(manCode string, sentenceName string) (string, error) {
	sentenceType := strings.ToLower(sentenceName)
	madeSentence := strings.ToUpper(manCode + sentenceName)

	if varList, found := h.sentences.formats[sentenceType]; found {
		for _, v := range varList {
			if vFormat, foundVar := h.sentences.variables[v]; foundVar {
				_, cv := getConversion(vFormat)
				if value, ok := h.data[v]; !ok || len(v) == 0 || v == "n/a" || len(value) == 0{
					for i:=0; i < cv.fCount; i++{
						madeSentence += ","
					}
				}else{
					
					madeSentence += "," + cv.to(value) 
				}
			}else{
				madeSentence += ","
			}
		}
	
		madeSentence = "$" + madeSentence
		checksum := checksum(madeSentence)
		madeSentence += "*" + checksum
		return madeSentence, nil
	}
	return "", fmt.Errorf("no matching sentence definition")
}
