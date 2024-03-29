package nmea0183

import (
	"fmt"
	"math"
	"testing"
)

func verify_sentence(sentence string, t *testing.T) *Handle {

	sentences := DefaultSentences()
	nm := sentences.MakeHandle()

	preFix, postFix, err := nm.Parse(sentence)
	if err != nil {
		t.Error("parsing input sentence error: %w", err)
	}
	ret_sentence, err := nm.WriteSentence(preFix, postFix)
	if err != nil {
		t.Error("writing output sentence error: %w", err)
	}
	if sentence != ret_sentence {
		t.Error(fmt.Errorf("parsed sentence not equal write sentence : %s != %s", sentence, ret_sentence))
	}

	return nm
}

func verify_sentence_match(sentence, expected string, t *testing.T) *Handle {

	sentences := DefaultSentences()
	nm := sentences.MakeHandle()

	preFix, postFix, err := nm.Parse(sentence)
	if err != nil {
		t.Error("parsing input sentence error: %w", err)
	}
	ret_sentence, err := nm.WriteSentence(preFix, postFix)
	if err != nil {
		t.Error("writing output sentence error: %w", err)
	}
	if expected != ret_sentence {
		t.Error(fmt.Errorf("expected sentence not equal write sentence : %s != %s", expected, ret_sentence))
	}

	return nm
}

func verify_sentence_prefixvars(sentence, prefixVar string, nm *Handle, t *testing.T) *Handle {

	preFix, postFix, err := nm.ParsePrefixVar(sentence, prefixVar)
	if err != nil {
		t.Error("parsing input sentence error: %w", err)
	}
	ret_sentence, err := nm.WriteSentencePrefixVar(preFix, postFix, prefixVar)
	if err != nil {
		t.Error("writing output sentence error: %w", err)
	}
	if sentence != ret_sentence {
		t.Error(fmt.Errorf("parsed sentence not equal write sentence : %s != %s", sentence, ret_sentence))
	}

	return nm
}

func TestConfig(t *testing.T) {
	var sentences Sentences
	e := sentences.Load("./example")
	if e != nil {
		errStr := fmt.Errorf("Config failed with error: %w", e)
		t.Error(errStr)
	}
}

func TestCheckSum(t *testing.T) {
	check := checksum("$9")
	expect := "39"
	if check != expect {
		t.Errorf("CheckSum was incorrect, got: %s, should be: %s.", check, expect)
	}

	check2 := checksum("$99")
	expect2 := "00"
	if check2 != expect2 {
		t.Errorf("CheckSum with even ones incorrect, got: %s, should be: %s", check2, expect2)
	}

	check3 := checksum("$ZZZ")
	expect3 := "5A"
	if check3 != expect3 {
		t.Errorf("CheckSum with even ones incorrect, got: %s, should be: %s", check3, expect3)
	}
}

func TestConvetLatLong(t *testing.T) {
	latStr := "50° 47.3986'N"
	longStr := "000° 54.6007'W"

	latFloat, longFloat, _ := LatLongToFloat(latStr, longStr)

	latr, longr, _ := LatLongToString(latFloat, longFloat)

	if latr != latStr || longr != longStr {
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

	if math.Abs(latr-latf) > 1e-6 || math.Abs(longr-longf) > 1e-6 {
		t.Errorf("lat long conversion error %f != %f or %f != %f ", latf, latr, longf, longr)
	}

	if lats[len(lats)-1:] != "S" {
		t.Error("Negative lat float must be South")
	}
	if longs[len(longs)-1:] != "E" {
		t.Error("Positive long float must be East")
	}

}

func TestDateConv(t *testing.T) {
	commaString := dateTimeToCSV("2021-09-15T11:09:10.59+01:40")
	if commaString != "110910.59,15,09,2021,01,40" {
		t.Errorf("Error time incorrectly parsed got %s", commaString)
	}
}

func TestZDA(t *testing.T) {
	var sentences Sentences
	sentences.Load("./example")
	nm := sentences.MakeHandle()

	nm.Parse("$GPZDA,110910.59,15,09,2020,00,00*6F")
	if nm.data["datetime"] != "2020-09-15T11:09:10.59+00:00" {
		t.Errorf("Error time incorrectly parsed got %s", nm.data["datetime"])
	}

	nm.Parse("$GPZDA,110910.59,15,09,21,01,40*69")
	if nm.data["datetime"] != "2021-09-15T11:09:10.59+01:40" {
		t.Errorf("Error time incorrectly parsed got %s", nm.data["datetime"])
	}

	verify_sentence("$GPZDA,110910.59,15,09,2020,01,30*6D", t)
	verify_sentence("$GPZDA,110910.59,15,09,2020,-01,-30*6D", t)
}

func TestZDACreate(t *testing.T) {
	zda := []string{"time", "day", "month", "year", "tz"}
	dpt := []string{"dbt", "toff"}

	formats := map[string][]string{"zda": zda, "dpt": dpt}
	variables := map[string]string{
		"time":  "hhmmss.ss",
		"day":   "DD_day",
		"month": "DD_month",
		"year":  "DD_year",
		"tz":    "tz_h,tz_m",
		"dpt":   "",
		"toff":  "",
	}
	sentences := MakeSentences(formats, variables)
	nm := sentences.MakeHandle()
	nm.Parse("$GPZDA,110910.59,15,09,2020,01,30*6D")
	if nm.data["time"] != "11:09:10.59" {
		t.Errorf("Error time incorrectly parsed got %s", nm.data["time"])
	}
	if nm.data["day"] != "15" || nm.data["month"] != "09" || nm.data["year"] != "2020" {
		t.Errorf("Incorrect data got %s %s %s", nm.data["day"], nm.data["month"], nm.data["year"])
	}
	if nm.data["tz"] != "+01:30" {
		t.Errorf("Error TZ time incorrectly parsed got %s", nm.data["tz"])
	}

	// Test loading another create is independant
	nm2 := sentences.MakeHandle()
	if len(nm2.data) != 0 || len(nm2.history) != 0 || len(nm.data) == 0 || len(nm.history) == 0 {
		t.Errorf("Second Create call failed - check that they are independent ")
	}
	nm3 := sentences.MakeHandle()
	nm3.Parse("$GPZDA,120910.59,15,09,2020,01,30*6E")
	if nm3.data["time"] != "12:09:10.59" {
		t.Errorf("Error time incorrectly parsed got %s", nm.data["time"])
	}

	nm4 := sentences.MakeHandle()
	nm4.Parse("$GPZDA,130910.59,15,09,2020,01,30*6F")
	if nm4.data["time"] != "13:09:10.59" {
		t.Errorf("Error time incorrectly parsed got %s", nm.data["time"])
	}

}
func TestAAM(t *testing.T) {
	verify_sentence("$GPAAM,A,A,0.10,N,WPTNME*32", t)
}

func TestAPA(t *testing.T) {
	verify_sentence("$GPAPA,A,A,8.30,L,M,V,V,11.7,T,Turning Track to Ijmuiden 1*1B", t)
	verify_sentence("$GPAPA,A,A,8.99,L,M,V,V,11.7,T,Turning Track to Ijmuiden 1*18", t)
	verify_sentence("$GPAPA,A,A,9.78,L,M,V,V,11.7,T,Turning Track to Ijmuiden 1*16", t)
	verify_sentence("$GPAPA,A,A,10.35,L,M,V,V,11.7,T,Turning Track to Ijmuiden 1*27", t)
}

func TestAPB(t *testing.T) {
	verify_sentence("$GPAPB,A,A,0.02617,R,N,V,V,210.0,T,Vlissingen,236.6,T,236.6,T,D*5D", t)
	verify_sentence("$GPAPB,A,A,0.02620,R,N,V,V,210.0,T,Vlissingen,236.7,T,236.7,T,D*59", t)
	verify_sentence("$GPAPB,A,A,0.02003,R,N,V,V,210.0,T,Vlissingen,227.6,T,227.6,T,D*5E", t)
	verify_sentence("$GPAPB,A,A,0.00536,R,N,V,V,210.0,T,Vlissingen,213.4,T,213.4,T,D*5F", t)
	verify_sentence("$GPAPB,A,A,5,L,N,V,V,359.,T,1,359.1,T,6,T,A*7C", t)
}

func TestHDM(t *testing.T) {
	verify_sentence("$HCHDM,172.5,M*28", t)

}

func TestRMC(t *testing.T) {
	verify_sentence("$GPRMC,,A,5125.1974,N,00043.4154,W,0.12,360,,0.65,E,A,V*63", t)

	verify_sentence("$GPRMC,110910.59,A,5047.3986,N,00054.6007,W,0.08,0.19,150920,0.24,W,D,V*75", t)

	verify_sentence("$GPRMC,163354.17,A,5222.5109,N,00502.8805,E,4.5,271.1,130319,,,D,V*24", t)
	verify_sentence("$GPRMC,163355.67,A,5222.5110,N,00502.8773,E,4.5,272.3,130319,,,D,V*25", t)
	verify_sentence("$GPRMC,163400.19,A,5222.5111,N,00502.8679,E,4.4,272.6,130319,,,D,V*25", t)
	verify_sentence("$GPRMC,163400.19,A,5222.5111,N,00502.8679,E,4.4,272.6,130319,,,D,V*25", t)
	verify_sentence("$GPRMC,163401.70,A,5222.5111,N,00502.8649,E,4.4,272.1,130319,,,D,A*38", t)
	verify_sentence("$GNRMC,001031.00,A,4404.1399,N,12118.8602,W,0.146,,100117,,,A,*57", t)

}

func TestPrefixVar(t *testing.T) {
	nm := verify_sentence("$GPRMC,163400.19,A,5222.5111,N,00502.8679,E,4.4,272.6,130319,,,D,V*25", t)
	verify_sentence_prefixvars("$GPRMC,163354.17,A,5222.5109,N,00502.8805,E,4.5,271.1,130319,,,D,V*24", "test_", nm, t)
	if nm.data["test_fix_time"] != "16:33:54.17" || nm.data["fix_time"] != "16:34:00.19" {
		t.Errorf(
			"Did not distinguish in PrefixVar expect test_fix_time == 16:33:54.17 got %s and fix_time == 16:34:00.19 got %s",
			nm.data["test_fix_time"], nm.data["fix_time"])
	}

}

func TestBlankMessages(t *testing.T) {
	verify_sentence("$GPRMC,,,,,,,,,,,,,*67", t)
	verify_sentence("$HCHDM,,*4A", t)
	verify_sentence("$GPAPB,,,,,,,,,,,,,,,*68", t)
	verify_sentence("$GPAPA,,,,,,,,,,*47", t)
	verify_sentence("$GPAAM,,,,,*76", t)

}

func TestSampleFromRaymarine(t *testing.T) {
	//Raymaine does not appear send 2 decimal places for seconds or omits trailing zero
	verify_sentence_match(
		"$GPZDA,113156.3,22,02,2024,00,00*52",
		"$GPZDA,113156.30,22,02,2024,00,00*62",
		t)
	verify_sentence_match(
		"$GPRMC,113157.3,A,5125.1974,N,00043.4154,W,0.12,360,220224,0.65,E,A,V*7A",
		"$GPRMC,113157.30,A,5125.1974,N,00043.4154,W,0.12,360,220224,0.65,E,A,V*4A",
		t)

	verify_sentence("$HCHDG,,,,0.7,E*00", t)

	//These are not defined in the default config and apart from some tech details are
	//such as datum, HOP etc and are already covered and would just overwrite key nav data  
	//verify_sentence("$GPGLL,5125.1974,N,00043.4154,W,113157.3,A,A*43", t)
	//verify_sentence("$GPGNS,113157.3,5125.1974,N,00043.4154,W,A,14,0.82,126.8,,,,V*7B", t)
	//verify_sentence("$GPGRS,113157.3,,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,,*4C", t)
	//verify_sentence("$GPVTG,360,T,,M,0.12,N,0.22,K,A*15", t)
	//verify_sentence("$GPDTM,W84,,0.0000,N,0.0000,E,0,W84*71", t)
	//verify_sentence("$GPGGA,113157.3,5125.1974,N,00043.4154,W,1,14,,126.8,M,,M,,*7E", t)
}
