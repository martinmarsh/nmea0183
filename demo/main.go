package main

import (
	"fmt"
	"github.com/martinmarsh/nmea0183"
)

func main() {

	// For easy start use built in sentences
	sentences := nmea0183.DefaultSentences()

	// Now create a handle to sentence processing, config and processed data
	nm := sentences.MakeHandle()

	// if you try to create a sentence without parsing information
	if gprmc, err := nm.WriteSentence("gp", "rmc"); err == nil {
		fmt.Println(gprmc)
	} else {
		// you will get an error and best attempt to create a sentence
		// only missing required fields not marked n/a will trigger an error
		fmt.Println(err, gprmc)
	}

	//fmt.Println(gprmc)

	// Now parse a sentence
	nm.Parse("$GPZDA,110910.59,15,09,2020,00,00*6F")

	// values parsed are merged into a Data map
	fmt.Println(nm.GetMap())

	nm.Parse("$GPRMC,110910.59,A,5047.3986,N,00054.6007,W,0.08,0.19,150920,0.24,W,D,V*75")
	fmt.Println(nm.GetMap())

	//Format of lat and long are readable strings
	fmt.Println(nm.Get("position"))

	//Can convert position variable to floats
	latFloat, longFloat, _ := nm.LatLongToFloat("position")
	fmt.Printf("%f %f\n", latFloat, longFloat)

	//Can write a variable from float lat and long
	nm.LatLongToString(latFloat, longFloat, "new_position")
	fmt.Println(nm.Get("new_position"))
	fmt.Println(nm.Date("position"))

	//examples of other sentences parsed
	nm.Parse("$HCHDM,172.5,M*28")
	nm.Parse("$GPAPB,A,A,5,L,N,V,V,359.,T,1,359.1,T,6,T,A*7C")
	nm.Parse("$SSDPT,2.8,-0.7")

	//Data is built and updated as sentences are parsed
	fmt.Println(nm.GetMap())

	gprmc, _ := nm.WriteSentence("gp", "rmc")
	fmt.Println(gprmc)

	hdm, _ := nm.WriteSentence("hc", "hdm")
	fmt.Println(hdm)

	apb, _ := nm.WriteSentence("gp", "apb")
	fmt.Println(apb)

	dpt, _ := nm.WriteSentence("ss", "dpt")
	fmt.Println(dpt)

	zda, _ := nm.WriteSentence("gp", "zda")
	fmt.Println(zda)

	//Data is built and updated as sentences are parsed also can
	fmt.Println(nm.GetMap())
	//also can get date of change
	fmt.Println(nm.DateMap())

}
