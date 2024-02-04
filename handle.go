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



// The Handle structure contains private data used to define sentences, configuarations, and parsed data.
//Methods on the struct allow parsing, updating of data and writing of sentences
type Handle struct {
	data                 map[string]string
	history				 map[string]int64
	messageDate			 time.Time
	upDated				 time.Time
	settings			 settings
	sentences			 *Sentences
}

// Returns a copy of the current data set or results of merged parsed sentences
func (h *Handle) GetMap() map[string]string{
	return h.data
}

// Returns a copy of a parsed variable in string format or null if not present
func (h *Handle) Get(key string) string {
	if val, ok := h.data[key]; ok {
		return val
	 } else {
		 return ""
	 }
}

// Returns a map of each data variable and the date and time it was updated 
func (h *Handle) DateMap() map[string]time.Time{
	dateMap := make(map[string] time.Time)
	for k, v := range(h.history){
		dateMap[k] = time.UnixMilli(v)
	}
	return dateMap
}

// Returns a copy of the date a variable was updated
func (h *Handle) Date(key string) time.Time {
	if val, ok := h.history[key]; ok {
		return time.UnixMilli(val)
	 } else {
		 return time.UnixMilli(0)
	 }
}


// Deletes variables in data which have a millisecond time stamp less than timeMS
func (h *Handle) DeleteBefore(timeMS int64){
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



// Adds the results of a parsed sentence to the handlers data set.
// This is normally done automatically when Parse is used.
// The update map is copied over any previous variables in the data assuming that
// the data being updated is most recent. Existing variables are not deleted so it can
// be used to "merge" results from different sentences.
// Caution: do not write directly to the results map unless you understand the string format
// associated with the variable in the sentence definitions
func (h *Handle) Update(results map[string]string) {
	
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

	vtypes :=  make(map[string]string)
	for n, v := range results {
		if template, found := h.sentences.variables[n]; found {
			tType, _:= getConversion(template)
			vtypes[tType] = v
		}
		h.data[n] = v
		h.history[n] = timeStamp
	}
	rcDate := ""
	if v, found := vtypes["datetime"]; found {
		rcDate = v
	} else if t, found := vtypes["time"]; found {
		date := ""
		if d, found := vtypes["date"]; found {
			date = d
		} else {
			day, okd := vtypes["day"]
			month, okm := vtypes["month"] 	
			year, oky := vtypes["year"]
			if okd && okm && oky {
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
		if len(date) > 0 {
			dateTime := date + "T" + t
			zone := "00:00"
			if z, ok := vtypes["zone"]; ok {
				zone = z
			}
			rcDate = dateTime + zone	
		}
	}
	if len(rcDate) > 0 {
		messageDate, err := time.Parse(time.RFC3339, rcDate)
		if err == nil{
			h.messageDate = messageDate
		}
	}
}

// Set handle settings preferences:
// autoClearPeriod = 0 for no automatic deletion of data or the value in seconds to keep data for
// realTime = True to use the processor clock, false to take the time from the sentences being parsed
func (h *Handle) Preferences(clear int64, realTime bool) {
	
	if clear < 1 {
		h.settings.autoClearPeriod = 0
	}else{
		h.settings.autoClearPeriod = clear * 1000
	}
	h.settings.realTime = realTime
}

// Parses a given sentence into handles data set using Update ie
// assumes the sentence parsed is more recent.
// Returns the GPS sentence device code and sentence name and any errors
// If check sum is fails the sentence will not update the data set and will be
// discarded.  The error returned should be checked to report the error.
// prefix returned is manufacturing/device code eg HC in $GPRMS,....
// sentence type is sentence code define in the config eg RMS in $GPRMS,...
//
// preFix, sentenceType, err = Parse(nmea_sentence)

func (h *Handle) Parse(nmea string) (string, string, error){
	data, preFix, postFix, error := h.ParseToMap(nmea)
	if error == nil{
		h.Update(data)
	}
	return preFix, postFix, error
}

// As Parse but writes to the data variables which have supplied your own prefix
// this would typically be used to distinguish in the data between variables from
// different sources or devices which produce identical sentences. 
//  
// results data map,  preFix, sentenceType, err = ParseToMap(nmea_sentence, variable_prefix)
// or if the sentence prefix can be used to distinguish:
// results data map,  preFix, sentenceType, err = ParseToMap(nmea_sentence, nmea_sentence[1:2])
//
func (h *Handle) ParsePrefixVar(nmea string, preFixVar string) (string, string, error){
	data, preFix, postFix, error := h.ParseToMap(nmea, preFixVar)
	if error == nil{
		h.Update(data)
	}
	return preFix, postFix, error
}

 
// Similar to Parse and ParsePrefixVar but does not update the data set but returns a map of the
// variables obtained from the sentence.
// This allows checking and filtering and optional updating of the data set
// using Update.  Be careful not to corrupt the string format returned if you are going to
// use in built functions to read the data or if you intend to update the data set.
//
// results data map,  preFix, sentenceType, err = ParseToMap(nmea_sentence)
// or
// results data map,  preFix, sentenceType, err = ParseToMap(nmea_sentence, variable_prefix)
//
func (h *Handle) ParseToMap(params ...string)  (map[string]string, string, string, error){
	l := len(params)
	nmea := ""
	var_prefix := ""
	if l >= 1 {
		nmea = strings.TrimSpace(params[0])
	}
	if l >= 2 {
		var_prefix = params[1]
	}
	if len(nmea)<5 || len(nmea)>89 || nmea[0] != '$'{
		return nil,"","",fmt.Errorf("%s", "Sentence must be between 5 and 89 and start with a $")
	}
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

	if varList, found := h.sentences.formats[sentenceType]; found{
		fieldPointer := 1
		for varPointer := 0; varPointer < len(varList); varPointer++ {
			conVar := ""
			varName := varList[varPointer]
			if template, found := h.sentences.variables[varName]; found && fieldPointer < len(parts){
				_, cv := getConversion(template)
				if fieldPointer + cv.fCount <= len(parts){
						conVar = cv.from(fieldPointer, &parts)
						// = varName
				} else {
						conVar = ""
				}
				fieldPointer += cv.fCount
				results[var_prefix + varName] = conVar
			} else {
				fieldPointer++
			}
		}
	}
	return results,  preFix, sentenceType, err
}


// As a method on the handler structure the string parameters refer to variable names in data
// Give one parameter for a variable name that holds a position string ie in form lat, Long or
// give the names of 2 variables  which hold lat and long respectively.
// Returns a 2 floats lat and long and an error.  Minus values for South and West
func (h *Handle) LatLongToFloat(params ...string) (float64, float64, error) {
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

// As a method on the data structure the string parameters refer to variable names in
// data which will be set by converting the lat and long float parameters
// One variable name assumes a position string containing lat, long will be set
// two variable names assumes that 2 variables one for lat and on for long will bes set
// Returns an error.  Assumes Minus values are for South and West
func (h *Handle) LatLongToString(latFloat, longFloat float64, params ...string) error {
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


// Writes a sentence using the handlers data and sentence definitions. 
// The sentence prefix is parsed in the first parameter followed by a string which 
// matches the sentence definitions. The prefix is added in the resulting string after the $
// and is included in the checksum. The prefix can be blank.
func (h *Handle) WriteSentence(manCode string, sentenceName string) (string, error) {
	return h.WriteSentencePrefixVar(manCode, sentenceName, "")
}

// This form of write sentence allows variables to be found in the data set which have been prefixed
// see ParsePrefixVar for how to create them

func (h *Handle) WriteSentencePrefixVar(manCode string, sentenceName string, prefixVar string) (string, error) {
	sentenceType := strings.ToLower(sentenceName)
	madeSentence := strings.ToUpper(manCode + sentenceName)

	if varList, found := h.sentences.formats[sentenceType]; found {
		for _, v := range varList {
			if vFormat, foundVar := h.sentences.variables[v]; foundVar {
				_, cv := getConversion(vFormat)
				if value, ok := h.data[prefixVar + v]; !ok || len(v) == 0 || v == "n/a" || len(value) == 0{
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
