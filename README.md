# nmea0183

## Post build configurable NMEA 0183 sentence analyser and writer in a go package

### Reason to choose

- Configurable Sentences definitions and variables either built in or via external yaml/json/toml files
- Collects data across sentences to merge into one results set
- Keeps results as strings for ease of logging, communication and visualisation
- Handles Removing old data

### Reason to consider other packages

- Do not want to look at your sentences to define parsing config and own variable names
- Want a solution that collects data into pre-defined structures
- Need to parse more easily multiple line sentences

### Current Status

Essentially functional but undergoing development and testing. Some experimental features may be added or removed as the package is being developed for use in my navigation system.
Features marked with * are under development and refinement, templates are being refined as further sentences tested.

### Features

- Sentences fully customisable and configurable via a Yaml/Json files
- Built in basic sentences can be used if preferred
- Check sum verification automatic if present
- Sentence definitions and config can be preloaded and multiple sentences passed
- Minimal processing for speed / processor use
- Results returned in string format for ease of print out, logging, transfer by UDP, channels etc
- Sentence values are collected in a GO map with user named keys
- Sentence fields can be ignored if not required
- Variables can be mapped to different sources for example when using two GPS systems
- Designed for continuous processing and data logging
- Can read extract date and time from multiple sentences - useful to set Raspberry Pi clock
- Auto removal of expired data
- Lat/long position from string to float
- Lat/Long Position from float to string

### Limitations

- No plans to support AIS
- 
- Only supports comma delimited fields and messages starting with $
- Limited to passing sentence fields which match built in templates

## Install

### Assuming you have installed go and are outide of GOPATH

    install go on your system

    go mod init  your_module
    go get github.com/martinmarsh/nmea0183

    write a main file as shown below or copy /example/main.go:

    go install your_module

### Basic use

To start write the following in main.go in your modules root directory
a copy of this file is in demo/main.go

    package main

    import (
     "fmt"
     "github.com/martinmarsh/nmea0183"
    )

    func main() {
        // For easy start use built in sentences and variable definitions
        nm:= nmea0183.Create()

        // Now parse a sentence
        nm.Parse("$GPZDA,110910.59,15,09,2020,00,00*6F")

        // values parsed are merged into a Data map
        fmt.Println(nm.Data)

        nm.Parse("$GPRMC,110910.59,A,5047.3986,N,00054.6007,W,0.08,0.19,150920,0.24,W,D,V*75")
        fmt.Println(nm.Data)

        //Format of lat and long are readable strings
        fmt.Println(nm.Data["position"])

        //Can convert position variable to floats
        latFloat, longFloat, _ := nm.LatLongToFloat("position")
        fmt.Printf("%f %f\n",latFloat, longFloat)

        //Can write a variable from float lat and long
        nm.LatLongToString(latFloat, longFloat, "new_position")
        fmt.Println(nm.Data["new_position"])

        //examples of other sentances passed
        nm.Parse("$HCHDM,172.5,M*28")
        nm.Parse("$GPAPB,A,A,5,L,N,V,V,359.,T,1,359.1,T,6,T,A*7C")
        nm.Parse("$SSDPT,2.8,-0.7")

        //Data is built and updated as sentences are parsed
        fmt.Println(nm.Data)

        }

### Defing your own sentences and mapping to variable

Built in sentences and variables may be sufficient for some applications but sooner or later some
bespoke configuration will be required. If you don't need the flexiblity of an external configuration
and are happy to build in the configuration into your compiled code you can parse your own sentence
definitions.  This is how:

    // example of 2 user defined sentences, just list the variable names to collect the data in
    zda := []string {"time","day","month","year","tz"}
    dpt := []string {"dbt","toff"}

    // These defintions are placed in map with the key matching the sentence name 
    sentences := map[string][]string {"zda": zda, "dpt": dpt}

    // Use create to load your only your sentences - no built in ones will be added
    nm := Create(sentences)

    // Now just parse sentences
    nm.Parse("$GPZDA,110910.59,15,09,2020,01,30*6F")

In the above example "zda" refers to any NMEA 0183 sentence after the first 2 digits, the manufacturer's code is
removed. The sentence definition refers to variables: "time","day","month","year","tz".  These variables
will use the built in definition.  When parsing a sentence each variable consumes the NMEA string in order ie
starting after the $GPZDA in our example "time" is taken from "110910.59", the "day" from "15", "month" from "09"
and "year" from "2020".  "tz" is an example of a variable which is taken from 2 fields ie 01,30
which is hours and minutes of time offset from UTC.

Other common formats consuming multiple fields are:  position, lat, long, mag_var. Creating variables from multiple fields is safer.
Position for example holds both lat and long and takes 4 fields eg in "$GPRMC,110910.59,A,5047.3986,N,00054.6007,W,0.08,0.19,150920,0.24,W,D,V*75"  the sequence  "5047.3986,N,00054.6007,W" is parsed to position which would be set to "50° 47.3986'N, 000° 54.6007'W".
This ensures that the position comes from one just sentence and the parts cannot be mixed up.

### Defing your own variable

The above uses built in variable definitions but you can configure this too.
Here is an example of variables which you could set to use instead of the default ones.

    // Lets use our own time variables pos_time amd zda_time amd map them in sentences
    sentences := map[string][]string {
        "rmc": {"pos_time", "status", "position", "sog", "tmg", "date", "mag_var"},
        "zda": {"zda_time", "day", "month", "year", "tz"},
    }
    sentences := map[string][]string {"zda": zda, "rmc": rmc}

    // now define some variables we might like to use mapping them to internal templates
    variables := map[string][]string {
        "pos_time": {"hhmmss.ss"},
        "zda_time": {"hhmmss.ss"},
        "date": {"ddmmyy"},
        "mag_var": {"x.x", "w"},              // Mag Var E positive, W negative
        "day": {"DD_day"},
        "month": {"DD_month"},
        "year": {"DD_year"},
        "tz":  {"tz_h", "tz:m"},              // Datetime from ZDA if available - tz:m returns hrs:mins
        "position": {"lat", "lat_NS", "pos_long", "pos_WE"}, 
        "sog": {"x.x"},                       // Speed Over Ground  float knots
        "tmg": {"x.x"},                       // Track Made Good
        "date": {"ddmmyy"},
        "status": {"A"},                      // status of fix A = ok ie 1 V = fail ie 0
        }


        nm := Create(sentences, variables)
        nm.Parse("$GPZDA,110910.59,15,09,2020,01,30*6F")
        nm.Parse("$GPRMC,110910.59,A,5047.3986,N,00054.6007,W,0.08,0.19,150920,0.24,W,D,V*75")

In the above example "pos_time" maps to built in format "hhmmss.ss" and so does "zda_time"
On parsing the strings the above 2 strings you will get both nm.Data["pos_time"] and nm.Data["zda_time"].
The default configuation just uses time in both sentences so you would get only nm.Data["time"]

### Using a config file instead of building in sentance definitions

For more fexibilty and the advantage of being able to change the parsing without having to install
go and rebuild configuations can be read and written to nmea_config.yaml file in the working directory.
Instead of Create but use Load.

    handle, err := nmea0183.Load()
    handle.Parse("$GPZDA,110910.59,15,09,2020,00,00*6F")

Can specify location of config file and type eg ymal or json

    handle, err := nmea0183.Load(".", "filename", "ymal")

A default Config file can be written using default settings:

    nmea0183.SaveConfig()
    or
    nmea0183.SaveConfig(".", "filename", "ymal")

### Cleaning up old data

By default Parse and Merge build Sentence data into a Go map called handle.Data

This allows a record of status to be easily maintained for example on a boat application
simply parsing all sentences as they are recieved allows the boat navigation status to be updated. This is
especially useful if the sentences are merged from different devices/inputs

However, there is a risk of being mislead by old data for example devices such as a GPS might loose signal
and the position might not be included in the sentence. The application will still see data but the position
may get out of date.

To reduce this risk the package notes the time of each variables and old ones deleted:

    handle.Preferences(seconds, real_time)   // seconds is time to keep variables or <= 0 is forever, real time mode true or false

use real_time = false if your system does not have a real time clock or historic data is being passed and of course the
messages contain a date.  Set seconds to zero or less to disable auto clear

This is not intended to cope with complete loss of connections when no sentences are Parsed ie for this
to work some sentances must still be Parsed or a Merge calls made.  

### Different channels

By choosing different config fies can use different handles to parse sentences differently. Filename1 may select different parts or names to filename2

    chan1, err1 := nmea0183.Load(".", "filename1", "ymal") 
    chan2, err2 := nmea0183.Load(".", "filename2", "ymal")

    results1 := chan1.Parse("$GPAPB,A,A,5,L,N,V,V,359.,T,1,359.1,T,6,T,A*7C") 

    results2 := chan2.Parse("$GPAPB,A,A,5,L,N,V,V,359.,T,1,359.1,T,6,T,A*7C")

    for example different configs might result in:
    results1:
     map[acir:V aper:V bod:359. bod_true:T bpd:359.1 device:GP did:1 hts:6 hts_true:T sentence:apb status:A xte:L5 xte_units:N]

    results2:
    map[mybod:359. my_hts:6 sentence:apb  device:GP]
