package nmea0183

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

var Sentences = make(map[string][]string)
var Variables = make(map[string][]string)

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
	viper.AddConfigPath(".")  // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		return err
	}
	
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := /*  */ err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			err = fmt.Errorf("fatal error config file err: %w not found = %t", err, ok)
			return err
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

    return fmt.Sprintf("%2X", check_sum)
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

