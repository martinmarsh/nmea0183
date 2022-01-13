/*
Copyright © 2022 Martin Marsh martin@marshtrio.com

*/
package main

import (
	"fmt"
	"github.com/martinmarsh/nmea0183"
)

func main() {
	err := nmea0183.Config()
	if err != nil{
		fmt.Println("Error config not found")
	}
	results := nmea0183.Parse("$GPZDA,110910.59,15,09,2020,00,00*6F")
	fmt.Println(results)

	results = (nmea0183.Parse("$GPRMC,110910.59,A,5047.3986,N,00054.6007,W,0.08,0.19,150920,0.24,W,D,V*75"))
	fmt.Println(results)
	fmt.Println(results["lat"] + " "+ results["long"])
	latFloat, longFloat := nmea0183.LatLongToFloat(results["lat"], results["long"])
	fmt.Printf("%f %f\n",latFloat, longFloat)
	fmt.Println(nmea0183.Parse("$HCHDM,172.5,M*28"))
	fmt.Println(nmea0183.Parse("$GPAPB,A,A,5,L,N,V,V,359.,T,1,359.1,T,6,T,A*7C"))
	fmt.Println(nmea0183.Parse("$SSDPT,2.8,-0.7"))

}

