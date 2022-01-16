package nmea0183

import (
	"testing"
    "fmt"
	"math"
)

func TestConfig(t *testing.T) {
    _, e := Load("./example")
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

	latFloat, longFloat, _ := LatLongToFloat(latStr, longStr)

	latr, longr, _ := LatLongToString(latFloat, longFloat)

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

	lats, longs, _ := LatLongToString(latf, longf)
	latr, longr, _ := LatLongToFloat(lats, longs)


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

func TestZDA(t *testing.T){
	nm, _ := Load("./example")
	nm.Parse("$GPZDA,110910.59,15,09,2020,00,00*6F")
	if nm.Data["time"] != "11:09:10.59" {
		t.Errorf("Error time incorrectly parsed got %s", nm.Data["time"])
	}
	if nm.Data["day"] != "15" || nm.Data["month"] != "09" ||nm.Data["year"] != "2020" {
		t.Errorf("Icorrect data got %s %s %s", nm.Data["day"], nm.Data["month"], nm.Data["year"])
	}
	if nm.Data["tz"] != "0.00" {
		t.Errorf("Error TZ time incorrectly parsed got %s", nm.Data["tz"])
	}

}


func TestZDACreate(t *testing.T){
	zda := []string {"time","day","month","year","tz"}
	dpt := []string {"dbt","toff"}

	sentences := map[string][]string {"zda": zda, "dpt": dpt}
	variables := map[string][]string {
		 "time": {"hhmmss.ss"},
		 "day": {"x"},
		 "month": {"x"},
		 "year": {"x"},
		 "tz": {"tz_h", "tz_m"},
		 "dpt": {},
		 "toff": {},	   
		}
	nm := Create(sentences, variables)
	nm.Parse("$GPZDA,110910.59,15,09,2020,01,30*6D")
	if nm.Data["time"] != "11:09:10.59" {
		t.Errorf("Error time incorrectly parsed got %s", nm.Data["time"])
	}
	if nm.Data["day"] != "15" || nm.Data["month"] != "09" ||nm.Data["year"] != "2020" {
		t.Errorf("Icorrect data got %s %s %s", nm.Data["day"], nm.Data["month"], nm.Data["year"])
	}
	if nm.Data["tz"] != "1.50" {
		t.Errorf("Error TZ time incorrectly parsed got %s", nm.Data["tz"])
	}
	
	// Test loading another create is independant
	nm2 := Create(sentences, variables)
	if len(nm2.Data) != 0 || len(nm2.History) !=0 || len(nm.Data) == 0 || len(nm.History) == 0 {
		t.Errorf("Second Create call failed - check that they are independan ")
	}
	nm3 := Create(sentences)
	nm3.Parse("$GPZDA,120910.59,15,09,2020,01,30*6E")
	if nm3.Data["time"] != "12:09:10.59" {
		t.Errorf("Error time incorrectly parsed got %s", nm.Data["time"])
	}
	nm4 := Create()
	nm4.Parse("$GPZDA,130910.59,15,09,2020,01,30*6F")
	if nm4.Data["time"] != "13:09:10.59" {
		t.Errorf("Error time incorrectly parsed got %s", nm.Data["time"])
	}
	
}

func TestNone(t *testing.T){
}