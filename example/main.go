/*
Copyright © 2022 Martin Marsh martin@marshtrio.com

*/
package main

import (
	"fmt"
	"github.com/martinmarsh/nmea0183"
)

func main() {
	// Load config file from working directory or create an example one if not found
	nm, err := nmea0183.Load()
	if err != nil{
		fmt.Println(fmt.Errorf("**** Error config: %w", err))
		nmea0183.SaveConfig()
		nm = nmea0183.Create()
	}

	// set time period in seconds to remove old values (<= 0 to disable) and if real time processing
	nm.Preferences(60, true)

	// use returned handle to Parse NMEA statements
	nm.Parse("$GPZDA,110910.59,15,09,2020,00,00*6F")
	// values parsed are merged into a Data map
	fmt.Println(nm.GetMap())

	nm.Parse("$GPRMC,110910.59,A,5047.3986,N,00054.6007,W,0.08,0.19,150920,0.24,W,D,V*75")
	fmt.Println(nm.GetMap())

	//Format of lat and long are readable strings
	fmt.Println(nm.Get("position"))

	//Can convert position variable to floats
	latFloat, longFloat, _ := nm.LatLongToFloat("position")
	fmt.Printf("%f %f\n",latFloat, longFloat)

	//Can write a variable from float lat and long
	nm.LatLongToString(latFloat, longFloat, "new_position")
	fmt.Println(nm.Get("new_position"))

	//examples of other sentances passed
	nm.Parse("$HCHDM,172.5,M*28")
	nm.Parse("$GPAPB,A,A,5,L,N,V,V,359.,T,1,359.1,T,6,T,A*7C")
	nm.Parse("$SSDPT,2.8,-0.7")

	//Data is built and updated as sentances parsed
	fmt.Println(nm.GetMap())

}
