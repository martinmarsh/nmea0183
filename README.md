# nmea0183
## Externally configurable nmea0183 sentence analyser in a go package

### Status:
Essentially functional but undergoing development and testing. Some experimental features may be added or removed as the package is being developed for use in my navigation system.
Features marked with * are under development and refinement, templates are being refined as further sentences tested.

### Features and reason to develop:

- Sentences fully customisable and configurable via a Yaml file so can be adjusted for device differences
- No built in sentence assumptions only a wide range of field formats built in.
- Check sum verification - Can be optional, automatic if present, or mandatory *
- Sentence definitions and config can be preloaded and multiple sentences passed
- Minimal processing for speed / processor use
- Results returned in string format for ease of print out, logging, transfer by UDP, channels etc 
- Results mapped to user named map keys.
- Sentence fields can be ignored if not required
- Variables can be mapped to different sources for example when using two GPS systems
- Designed for continuous processing and data logging
- Optional Results processing functions*
    - Post parsing results to extract date and time across multiple sentences -useful to set Raspberry Pi clock
    - Post parsing result to map to key names referencing channel, sentance and/or device
    - Current status processing to collect data from multiple sentences and remove expired data
    - Lat/long position from string to float
    - Lat/Long Position from float to string

### Limitations

- No plans to support AIS
- Only supports comma delimited fields and messages starting with $
- Limited to passing sentence fields which match built templates

## Install:

### Assuming you have installed go and are outide of GOPATH:
    go mod init  your_module
    go get github.com/martinmarsh/nmea0183
	go get github.com/spf13/viper

    copy /example/main to same directory to same directory as go.mod
    copy /example/nmea_config.yaml to same directory as main.go

    go install your_module

    run your code and update yaml file as required 


To use:

See main.go in example

eg
    import (
	    "fmt"
	    "github.com/martinmarsh/nmea0183"
    )

    nmea0183.Config()
    
	results := nmea0183.Parse("$GPZDA,110910.59,15,09,2020,00,00*6F")
	fmt.Println(results)
	fmt.Println(nmea0183.Parse("$GPRMC,110910.59,A,5047.3986,N,00054.6007,W,0.08,0.19,150920,0.24,W,D,V*75"))
	fmt.Println(nmea0183.Parse("$HCHDM,172.5,M*28"))
	fmt.Println(nmea0183.Parse("$GPAPB,A,A,5,L,N,V,V,359.,T,1,359.1,T,6,T,A*7C"))
	fmt.Println(nmea0183.Parse("$SSDPT,2.8,-0.7"))

