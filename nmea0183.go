package nmea0183

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

var Sentences = make(map[string][]string)
var Variables = make(map[string][]string)


func Config(){
	viper.SetConfigName("nmea_config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name

	viper.AddConfigPath(".")    // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := /*  */ err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			panic(fmt.Errorf("fatal error config file err: %w not found = %t", err, ok))
		} else {
			// Config file was found but another error was produced
			panic(fmt.Errorf("fatal error in config file: %w", err))
		}
	}


	Sentences = viper.GetStringMapStringSlice("sentences")
	Variables = viper.GetStringMapStringSlice("variables")
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
			varName, templateList := findInMap(varList[varPointer],Variables)
			conVar := ""
			for _, template := range(templateList){
				conVar = convert(parts[fieldPointer], template, conVar)
				fieldPointer ++
			}
			results[varName] = conVar
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
		if data == "L" || data == "l"{
			return "-" + conVar
		}
		return conVar
	case "llll.lll":
		d, _ := strconv.ParseFloat(data[:2], 32)
		m, _ := strconv.ParseFloat(data[2:], 32)
		d += m/60
		return fmt.Sprintf("%2.8f", d)

	case "yyyyy.yyyy":
		d, _ := strconv.ParseFloat(data[:3], 32)
		m, _ := strconv.ParseFloat(data[3:], 32)
		d += m/60
		return fmt.Sprintf("%3.8f", d)

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

