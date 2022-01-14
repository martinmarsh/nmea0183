package nmea0183

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"github.com/spf13/viper"
)

var Sentences = make(map[string][]string)
var Variables = make(map[string][]string)

/*
The following is a get started config. used to generate nmea_conf.yaml
if it is missing in root directory

*/
var yamlExample = []byte(`
sentences:
    RMC:
        - "time"
        - "status"
        - "lat"
        - "long"
        - "sog"
        - "tmg"
        - "date"
        - "mag_var"
    ZDA:
        - "time"
        - "day"
        - "month"
        - "year"
        - "tz"
    APB:
        - "status"
        - "n/a"
        - "xte"
        - "xte_units"
        - "acir"
        - "aper"
        - "bod"
        - "bod_true"
        - "did"
        - "bpd"
        - "bpd_true"
        - "hts"
        - "hts_true"
    HDG:
        - "n/a"
        - "n/a" 
        - "n/a"
        - "mag_var"
    HDM:
        - "hdm"
    DPT:
        - "dbt"
        - "toff"
    VHW:
        - "n/a"
        - "n/a"
        - "n/a"
        - "n/a"
        - "stw"
    VLW:
        - "n/a"
        - "n/a"
        - "wd"

variables:
    time: "hhmmss.ss"  # time of fix
    status: "A"   # status of fix A = ok ie 1 V = fail ie 0
    lat:
        - "llll.lll"
        - "NS"  # lat N / S postfix
    long:
        - "yyyyy.yyyy"
        - "WE"   # long float W/E postfix
    sog: "x.x"  # Speed Over Ground  float knots
    tmg: "x.x"  # Track Made Good
    date: "ddmmyy" # Date of fix may not be valid with some GPS
    mag_var:
        - "x.x"
        - "w"     # Mag Var E positive, W negative
    day: "x"
    month: "x"
    year: "x"
    tz:
        - "tz_h"
        - "tz_m"  # Datetime from ZDA if available - tz return hours and mins as a float
    xte:
        - "x.x"
        - "R"           # Cross Track Error turn R or L
        
    xte_units: "A"      # Units for XTE - N = Nm
    acir: "A"           # Arrived at way pt circle
    aper: "A"           # Perpendicular passing of way pt
    bod: "x.x"
    bod_true: "T"        # Bearing origin to destination True
    did: "str"             # Destination Waypoint ID as a str
    bpd: "x.x"
    bdp_true: "T"        # Bearing, present position to Destination True
    hts: "x.x"
    hts_true: "T"        # Heading to Steer True
    hdm: "x.x"          # Heading Magnetic
    dbt: "x.x"          # Depth below transducer
    toff: "x.x"         # Transducer offset -ve from transducer to keel +ve transducer to water line
    stw: "x.x"          # Speed Through Water float knots
    dw:  "x.x"             # Water distance since reset float knots

`)

/*type nmea interface {
    area() float64
    perim() float64
}*/


func Config(setting ...string) error{
	configSet := []string{".", "nmea_config", "yaml"}
	copy (configSet, setting)

	viper.SetConfigName(configSet[1]) // name of config file (without extension)
	viper.SetConfigType(configSet[2])   // REQUIRED if the config file does not have the extension in the name

	viper.AddConfigPath(configSet[0])    // optionally look for config in the working directory
	viper.AddConfigPath("..")  // optionally look for config in aprent of the working directory

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil{
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found

			os.WriteFile("nmea_config.yaml", yamlExample, 0644)
    
			err = viper.ReadInConfig()
			if err != nil{
				err = fmt.Errorf("fatal error in config file: %w", err)
				return err
			}

		} else {
			// Config file was found but another error was produced
			err = fmt.Errorf("fatal error in config file: %w", err)
			return err
		}
	}
	
	Sentences = viper.GetStringMapStringSlice("sentences")
	Variables = viper.GetStringMapStringSlice("variables")
	return err
}

func Parse(nmea string)  map[string]string {
	end_byte:= len(nmea)
	if nmea[end_byte -3] == '*'{
		check_code := checksum(nmea)
		end_byte -= 2
		if check_code != nmea[end_byte:]{
			fmt.Printf("check sum error: %s != %s\n", check_code, nmea[end_byte:])
		}
		end_byte --
	}

	parts := strings.Split(nmea[1:end_byte], ",")
	preFix := parts[0][:2]
	sentenceType := strings.ToLower(parts[0][2:])
	key, varList := findInMap(sentenceType, Sentences)
	results := map[string]string {"device": preFix, "sentence": sentenceType}

	if len(key) > 0{
		fieldPointer := 1

		for varPointer :=0; varPointer < len(varList); varPointer++ {
			varName, templateList := findInMap(varList[varPointer], Variables)
			conVar := ""
			if len(varName) > 0{
				for _, template := range(templateList){
					conVar = convert(parts[fieldPointer], template, conVar)
					fieldPointer ++
				}
				results[varName] = conVar
			}else{
				fieldPointer ++
			}
		}
	}
	return results	
}

func convert(data string, template string, conVar string) string{
	switch template{
	case "hhmmss.ss":
		h, e := strconv.ParseInt(data[:2], 10, 16)
		m, e1 := strconv.ParseInt(data[2:4], 10, 16)
		s, e2 := strconv.ParseFloat(data[4:], 32)
		if e == nil && e1 == nil && e2 == nil {
			return fmt.Sprintf("%02d:%02d:%02.2f", h, m, s)
		}
		return ""
	case "x.x":
		return data
	case "x":
		return data
	case "tz_h":
		return data
	case "tz_m":
		h, e := strconv.ParseFloat(conVar, 32)
		m, e1 := strconv.ParseFloat(data, 32)
		if e == nil && e1 == nil{
			h += m/60
			return fmt.Sprintf("%02.2f", h)
		}
	case "A":
		return data

	case "str":
		return data

	case "T":
		return data

	case "w":
		if data == "W" || data == "w"{
			return "-" + conVar
		}
		return conVar

	case "s":
		if data == "S" || data == "s"{
			return "-" + conVar
		}
		return conVar

	case "R":
		return data + conVar

	case "llll.lll":
		d, _ := strconv.ParseInt(data[:2], 10, 32)
		m, _ := strconv.ParseFloat(data[2:], 32)
		return fmt.Sprintf("%02d° %02.4f'", d, m)

	case "yyyyy.yyyy":
		d, _ := strconv.ParseInt(data[:3], 10, 32)
		m, _ := strconv.ParseFloat(data[3:], 32)
		return fmt.Sprintf("%03d° %02.4f'", d, m)

	case "WE":
		return conVar + data

	case "NS":
		return conVar + data

	case "ddmmyy":
		d, _ := strconv.ParseInt(data[:2], 10, 32)
		m, _ := strconv.ParseInt(data[2:4], 10, 32)
		y, _ := strconv.ParseInt(data[4:], 10, 32)
		if y < 60 {
			y += 2000
		}else{
			y += 1900
		}
		return fmt.Sprintf("%d-%02d-%02d", y ,m, d)
	
	}	
	return ""
}

func checksum(s string) string {
	check_sum:= 0

    nmea_data := []byte(s)

    for  i:= 1; i < len(s) - 3; i ++ {
        check_sum ^= (int)(nmea_data[i])
    }

    return fmt.Sprintf("%02X", check_sum)
}

func findInMap(k string, m map[string][]string) (string, []string){
	
	for i, v := range(m){
		if i == k{
			return i, v
		}
	}
	return "", []string{""}
}

func LatLongToFloat(lat string, long string)(float64, float64){
	
	lenL := len(lat)-1
	symbol := lat[lenL]
	lenL--
	dlat, _ := strconv.ParseFloat(lat[:2], 64)
	mlat, _ := strconv.ParseFloat(lat[5:lenL], 64)
	dlat += mlat/60
	if symbol == 'S' {
		dlat = -dlat
	}
	lenL = len(long) - 1
	symbol = long[lenL]
	lenL--	
	dlong, _ := strconv.ParseFloat(long[:3], 32)
	mlong, _ := strconv.ParseFloat(long[6:lenL], 32)
	dlong += mlong/60
	if symbol == 'W' {
		dlong = -dlong
	}

	return dlat, dlong
}

func LatLongToString(latFloat float64, longFloat float64)(string, string){
	latAbs := math.Abs(latFloat)
	latInt := int(latAbs)
	latMins := (latAbs - float64(latInt)) * 60
	symbol := "N"
	if latFloat < 0 {
		symbol = "S"
	}

	lat :=  fmt.Sprintf("%02d° %02.4f'%s", latInt, latMins, symbol)

	longAbs := math.Abs(longFloat)
	longInt := int(longAbs)
	longMins := (longAbs - float64(longInt)) * 60
	symbol = "E"
	if longFloat < 0 {
		symbol = "W"
	}

	long := fmt.Sprintf("%03d° %02.4f'%s", longInt, longMins, symbol)

	return lat, long
}
