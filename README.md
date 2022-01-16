# nmea0183

## Externally configurable nmea0183 sentence analyser in a go package

### Reason to choose

- Configure per message and per source parsing of only defined sentances into named variables
- Collects data across sentances to merge into one results set
- Keeps results as strings for ease of logging, communication and visualisation
- Handles Removing old data
- Easy conversion to own data stuctures

### Reason to consider other packages

- Do not want to look at your sentances to define parsing config and own variable names
- Want a solution that collects data into pre-defined package structures
- Want a pre-configured extensive sentence parser

### Current Status

Essentially functional but undergoing development and testing. Some experimental features may be added or removed as the package is being developed for use in my navigation system.
Features marked with * are under development and refinement, templates are being refined as further sentences tested.

### Features

- Sentences fully customisable and configurable via a Yaml file
- Built in basic sentences can be used if preferred
- Check sum verification automatic if present
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

## Install

### Assuming you have installed go and are outide of GOPATH

    go mod init  your_module
    go get github.com/martinmarsh/nmea0183
    go get github.com/spf13/viper

    from download in GOPATH or from git hub either:
        copy /example/main to same directory to same directory as go.mod
        copy nmea_config.yaml to same directory as main.go
    or 
        write main as shown below:

    go install your_module
 
    run your code

### Basic use

See main.go in example

eg

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

Using built in sentences and variables is easy start but has is not configurable.
You can passs your own senatnce definitions instead of using the default ones:

    zda := []string {"time","day","month","year","tz"}
    dpt := []string {"dbt","toff"}
    sentences := map[string][]string {"zda": zda, "dpt": dpt}

    nm := Create(sentences)
    nm.Parse("$GPZDA,110910.59,15,09,2020,01,30*6F")

In the above example "zda" refers to any NMEA 0183 sentence starting with a 2 digit manufacturer's code eg
GPZDA.  The sentence definition refers to variables: "time","day","month","year","tz".  These variables
will use the built in definition.  When parsing a sentence each variable consumes the NMEA string in order ie
starting after the $GPZDA in our example "time" is taken from "110910.59", the "day" from "15", "month" from "09"
and "year" from "2020".  "tz" is an example of a variable which is taken from 2 fields ie 01,30
which is hours and minutes of time offset from UTC.  "tz" maps to a special format which takes the hour and the minutes
and returns a decimal hour offset.

Other common formats consuming multiple fields are position, lat, long, mag_var, and x
Creating variables from multiple fields is safer especially when mulitple sentences are analysed and there is a risk
of mixing up parts of what should be one element of date.

Consider a boat near a time line crossing E and W if east and west variables are separate from the longitude degrees and minutes as obtained from a sentance there is a risk when combining data that the variables may be updated differently. In fact position is even safer because it holds both lat and long and takes 4 fields eg in "$GPRMC,110910.59,A,5047.3986,N,00054.6007,W,0.08,0.19,150920,0.24,W,D,V*75"  the sequence  "5047.3986,N,00054.6007,W" is parsed to position which would be set to "50° 47.3986'N, 000° 54.6007'W"

### Defing your own variable

Here is an example of variables which you could set to use instead of the default ones.

    sentences := map[string][]string {
        "rmc": {"pos_time", "status", "position", "sog", "tmg", "date", "mag_var"},
        "zda": {"zda_time", "day", "month", "year", "tz"},
    }
    sentences := map[string][]string {"zda": zda, "rmc": rmc}

    variables := map[string][]string {
        "pos_time": {"hhmmss.ss"},
        "zda_time": {"hhmmss.ss"},
        "status": {"A"},
        "lat": {"llll.lll", "NS"},
        "long": {"yyyyy.yyyy","WE"},
        "position": {"llll.lll", "NS", "yyyyy.yyyy", "WE"},
        "day": {"x"},
        "month": {"x"},
        "year": {"x"},
        "tz": {"tz_h", "tz_m"},
        "dpt": {},
        "toff": {},
        }
        nm := Create(sentences, variables)
        nm.Parse("$GPZDA,110910.59,15,09,2020,01,30*6F")
        nm.Parse("$GPRMC,110910.59,A,5047.3986,N,00054.6007,W,0.08,0.19,150920,0.24,W,D,V*75")

In the above example "pos_time" maps to built in format "hhmmss.ss" and so does "zda_time"
On parsing the strings the above 2 strings you will get both nm.Data["pos_time"] and nm.Data["zda_time"].
The default configuation just uses time in both sentences so you would get only nm.Data["time"]

### Using a config file instead of building in sentance definitions

The problem with built in sentences and variables is that once a build has been done it is
impossible to change if sentences need to be updated. By using an external Yaml file
configuration and variable names can be read in at the start of your application. This is especially useful in all
applications need to process the values obtained from sentances and just log them or pass them on
via a file or communication channel.

To read the default nmea_config.yaml file from the working directory do not use Create but use instead.

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
and stop sending position in sentences. The application will still see data but the position in handle.Data may be
dangerously out of date.

To reduce this risk the package notes the time of each variables and old ones deleted:

    handle.AutoClear(seconds, true)   // automatically clears variables older than given number of seconds

use false if your system does not have a real time clock or historic data is being passed and of course the
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
