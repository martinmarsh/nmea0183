# nmea0183
## Externally configurable nmea0183 sentence analyser in a go package
### Reason to choose:
- Able configure per message source parsing of defined sentances into named variables
- Collects data across sentances to merge into one results set
- Keeps results as strings for ease of logging, communication and visualisation
- Easy conversion to own data stuctures
### Reason to consider other packages
- Do not wnat to look at your sentances to define parsing config and own variable names
- Want a solution that collects data into pre-defined package structures 
- Want a pre-configured extensive sentence parser

### Current Status:
Essentially functional but undergoing development and testing. Some experimental features may be added or removed as the package is being developed for use in my navigation system.
Features marked with * are under development and refinement, templates are being refined as further sentences tested.

### Features:
- Sentences fully customisable and configurable via a Yaml file 
- No built in sentence assumptions and only a wide range of field formats built in.
- Check sum verification - Can be optional, automatic if present, or mandatory *
- Sentence definitions and config can be preloaded and multiple sentences passed
- Minimal processing for speed / processor use
- Results returned in string format for ease of print out, logging, transfer by UDP, channels etc 
- Results mapped to user named map keys.
- Sentence fields can be ignored if not required
- Variables can be mapped to different sources for example when using two GPS systems
- Designed for continuous processing and data logging
- Optional Results processing functions*
    - Can read latest date and time across multiple sentences -useful to set Raspberry Pi clock
    - Post parsing result to map to key names referencing channel, sentence and/or device
    - Current status processing to collect data from multiple sentences and remove expired data
    - Lat/long position from string to float
    - Lat/Long Position from float to string

### Limitations

- No plans to support AIS
- Only supports comma delimited fields and messages starting with $
- Limited to passing sentence fields which match built in templates

## Install:

### Assuming you have installed go and are outide of GOPATH:
    go mod init  your_module
    go get github.com/martinmarsh/nmea0183
	go get github.com/spf13/viper

    from download in GOPATH or from git hub either:
        copy /example/main to same directory to same directory as go.mod
        copy nmea_config.yaml to same directory as main.go
    or 
        write main as shown below - on 1st run a nmea_config file will be created for you

    go install your_module
 
    run your code and update yaml file as required 


### Basic use:

See main.go in example

eg

    package main

    import (
	    "fmt"
	    "github.com/martinmarsh/nmea0183"
    )

    func main() {
        // Load config file from working directory or create an example one if not found
	    nm, err := nmea0183.Load()
	    if err != nil{
		    fmt.Println("Error config not found")
	    }

	    // use returned handle to Parse NMEA statements
	    nm.Parse("$GPZDA,110910.59,15,09,2020,00,00*6F")
	    // values parsed are merged into a Data map
	    fmt.Println(nm.Data)

	    nm.Parse("$GPRMC,110910.59,A,5047.3986,N,00054.6007,W,0.08,0.19,150920,0.24,W,D,V*75")
	    fmt.Println(nm.Data)

	    //Format of lat and long are readable strings
	    fmt.Println(nm.Data["lat"] + " "+ nm.Data["long"])

	    //Can convert lat and long to floats
	    latFloat, longFloat := nm.LatLongToFloat(nm.Data["lat"], nm.Data["long"])
	    fmt.Printf("%f %f\n",latFloat, longFloat)

	    //examples of other sentances passed
	    nm.Parse("$HCHDM,172.5,M*28")
	    nm.Parse("$GPAPB,A,A,5,L,N,V,V,359.,T,1,359.1,T,6,T,A*7C")
	    nm.Parse("$SSDPT,2.8,-0.7")

	    //Data is built and updated as sentances parsed
	    fmt.Println(nm.Data)
    }
### More advances use:

Can specify location of config file and type eg ymal or json

    handle, err := nmea0183.Load(".", "filename", "ymal") 

By choosing different config fies can use differnet handles to parse sentences differently. Filename1 may select different parts or names to filename2

    chan1, err1 := nmea0183.Load(".", "filename1", "ymal") 
    chan2, err2 := nmea0183.Load(".", "filename2", "ymal")

    results1 := chan1.Parse("$GPAPB,A,A,5,L,N,V,V,359.,T,1,359.1,T,6,T,A*7C") 

    results2 := chan2.Parse("$GPAPB,A,A,5,L,N,V,V,359.,T,1,359.1,T,6,T,A*7C")

    for example different configs might result in:
    results1:
     map[acir:V aper:V bod:359. bod_true:T bpd:359.1 device:GP did:1 hts:6 hts_true:T sentence:apb status:A xte:L5 xte_units:N]

    results2:
    map[mybod:359. my_hts:6 sentence:apb  device:GP]

If don't wnat to use a config file can hard code into your module settings and use Create instead of
load

    handle = NMEA0183.Create(Sentances, Variables)
