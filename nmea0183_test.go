package nmea0183

import (
	"testing"
    "fmt"
	"math"
)

func TestConfig(t *testing.T) {
    _, e := Load()
	if e != nil {
		errStr := fmt.Errorf("Config failed with error: %w", e)
		t.Error(errStr)
	}
}

func TestCheckSum(t *testing.T) {	
	check := checksum("$1111111*45")
	expect := "31"
	if check != expect {
		t.Errorf("CheckSum was incorrect, got: %s, should be: %s.", check, expect)
	}

	check2 := checksum("$111111*45")
	expect2 := "00"
	if check2 != expect2 {
		t.Errorf("CheckSum with even ones incorrect, got: %s, should be: %s", check2, expect2)
	}
}

func TestConvetLatLong(t *testing.T) {
	latStr := "50° 47.3986'N"
	longStr := "000° 54.6007'W"

	latFloat, longFloat := LatLongToFloat(latStr, longStr)

	latr, longr := LatLongToString(latFloat, longFloat)

	if latr != latStr || longr != longStr{
		t.Errorf("lat long conversion error %s != %s or %s != %s ", latStr, latr, longStr, longr)
	}

	if !(longFloat < 0) {
		t.Error("Westerly lat float must be neagrive")
	}
}
	
func TestConvetFloatLatLong(t *testing.T) {
	latf := -76.12345
	longf := 170.54321

	lats, longs := LatLongToString(latf, longf)
	latr, longr := LatLongToFloat(lats, longs)


	if math.Abs(latr-latf) > 1E-6 || math.Abs(longr-longf) > 1e-6 {
		t.Errorf("lat long conversion error %f != %f or %f != %f ", latf, latr, longf, longr)
	}

	if lats[len(lats)-1:] != "S" {
		t.Error("Negative lat float must be South")
	}
	if longs[len(longs)-1:] != "E" {
		t.Error("Positive long float must be East")
	}

}


func TestZDA(*testing.T){
	nm, _ := Load()
	nm.Parse("$GPZDA,110910.59,15,09,2020,00,00*6F")
	fmt.Println(nm.Data)

}

func TestZDA2(*testing.T){
	zda := []string {"time","day","month","year","tz"}
	sentences := map[string][]string {"zda": zda}
	variables := map[string][]string {
		 "time": {"hhmmss.ss"},
		 "day": {"x"},
		 "month": {"x"},
		 "year": {"x"},
		 "tz": {"tz_h", "tz_m"},	   
		}
	nm := Create(sentences, variables)
	nm.Parse("$GPZDA,110910.59,15,09,2020,00,00*6F")
	fmt.Println(nm.Data)

}