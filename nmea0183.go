package nmea0183

import (
	"fmt"
	"github.com/spf13/viper"
	"math"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Sentences, Variables map[string][]string
	Data                 map[string]string
	History				 map[string]int32
	Counter				 int32
}

func setUp() *Config{
	var c Config
	c.Data = make(map[string]string)
	c.History = make(map[string]int32) 
	c.Counter = 0
	return &c
}

func Create(sentances, variables map[string][]string) *Config {
	c := setUp()
	c.Sentences = sentances
	c.Variables = variables
	return c
}

func Load(setting ...string) (*Config, error) {
	c := setUp()
	configSet := []string{".", "nmea_config", "yaml"}
	copy(configSet, setting)

	viper.SetConfigName(configSet[1]) // name of config file (without extension)
	viper.SetConfigType(configSet[2]) // REQUIRED if the config file does not have the extension in the name

	viper.AddConfigPath(configSet[0]) // optionally look for config in the working directory
	viper.AddConfigPath("..")         // optionally look for config in aprent of the working directory

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found and looked for yaml type trying writing a default one
			if configSet[2] == "yaml" {
				os.WriteFile(configSet[0]+"/"+configSet[1]+".yaml", yamlExample, 0644)

				err = viper.ReadInConfig()
				if err != nil {
					err = fmt.Errorf("fatal error in config file: %w", err)
					return c, err
				}
			}

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

func (c *Config) Merge(results map[string]string) {
	c.Counter ++
	for n, v := range results {
		c.Data[n] = v
		c.History[n] = c.Counter
	}
}

func (c *Config) Parse(nmea string) {
	end_byte := len(nmea)
	if nmea[end_byte-3] == '*' {
		check_code := checksum(nmea)
		end_byte -= 2
		if check_code != nmea[end_byte:] {
			fmt.Printf("check sum error: %s != %s\n", check_code, nmea[end_byte:])
		}
		end_byte--
	}

	parts := strings.Split(nmea[1:end_byte], ",")
	preFix := parts[0][:2]
	sentenceType := strings.ToLower(parts[0][2:])
	key, varList := findInMap(sentenceType, c.Sentences)
	results := map[string]string{"device": preFix, "sentence": sentenceType}

	if len(key) > 0 {
		fieldPointer := 1

		for varPointer := 0; varPointer < len(varList); varPointer++ {
			varName, templateList := findInMap(varList[varPointer], c.Variables)
			conVar := ""
			if len(varName) > 0 {
				for _, template := range templateList {
					conVar = convert(parts[fieldPointer], template, conVar)
					fieldPointer++
				}
				results[varName] = conVar
			} else {
				fieldPointer++
			}
		}
	}
	c.Merge(results)
}

func convert(data string, template string, conVar string) string {
	switch template {
	case "hhmmss.ss":
		h, e := strconv.ParseInt(data[:2], 10, 16)
		m, e1 := strconv.ParseInt(data[2:4], 10, 16)
		s, e2 := strconv.ParseFloat(data[4:], 32)
		if e == nil && e1 == nil && e2 == nil {
			return conVar + fmt.Sprintf("%02d:%02d:%02.2f", h, m, s)
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
		if e == nil && e1 == nil {
			h += m / 60
			return fmt.Sprintf("%02.2f", h)
		}
	case "A":
		return data

	case "str":
		return data

	case "T":
		return data

	case "w":
		if data == "W" || data == "w" {
			return "-" + conVar
		}
		return conVar

	case "s":
		if data == "S" || data == "s" {
			return "-" + conVar
		}
		return conVar

	case "R":
		return data + conVar

	case "llll.llll":
		d, _ := strconv.ParseInt(data[:2], 10, 32)
		m, _ := strconv.ParseFloat(data[2:], 32)
		return fmt.Sprintf("%02d° %02.4f'", d, m)

	case "yyyyy.yyyy":
		d, _ := strconv.ParseInt(data[:3], 10, 32)
		m, _ := strconv.ParseFloat(data[3:], 32)
		return fmt.Sprintf("%03d° %02.4f'", d, m)
	
	case ",yyyyy.yyyy":
		d, _ := strconv.ParseInt(data[:3], 10, 32)
		m, _ := strconv.ParseFloat(data[3:], 32)
		return conVar + ", " + fmt.Sprintf("%03d° %02.4f'", d, m)

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
		} else {
			y += 1900
		}
		return fmt.Sprintf("%d-%02d-%02d", y, m, d)

	}
	return ""
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

	if len(params) == 1 {
		c.Data[params[0]] = latStr + ", " + longStr
	} else {
		c.Data[params[0]] = latStr
		c.Data[params[1]] = longStr
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
