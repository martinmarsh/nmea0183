# nmea0183

## Post build configurable NMEA 0183 sentence analyser and writer in a go package

This package was created to so solve a problem of logging data produced by many NMEA0183 and 2000 devices (via a bridge Actisense NGW-1). It has been proven on a system using USB com ports, OpenCPN, Actisence NGW-1, Raymarine Lighthouse, Garmin GPS, Digital Compass, AIS (message muliplexing only) and hosted on RaspberryPI and with a change to config files will run on Windows

The package has the following key features:

* Multiple Sentence types can be parsed into a variable map (dictionary) so as to creates a database of current status. 

* To avoid duplication in the data set and create a single point of truth sentences are parsed to logical variables as as position, sog etc. A Variable can be shared with different sentence types.

* For efficiency and accuracy the data map uses formatted strings such as position ==  "50° 47.3986'N, 000° 54.6007'W"
i.e. in a standard format with correct symbols

* Conversions to many float and int types when needed to process in code. Read and writing sentences uses internal string format.

* Can add configuration for new sentences and variables without re-compiling go files.

* Sentences parsed in are often reproduced exactly with same checksum



This package is used by the NMEA0183_MUX library which might also be useful to route, duplicate and generate NMEA messages to
different devices.

## Quick start:

```go

var sentences nmea0183.Sentences

nm := sentences.DefaultConfig()

// then just use Parse NMEA strings as they arrive
// using some examples of a couple of possible strings

nm.Parse("$GPZDA,110910.59,15,09,2020,00,00*6F")
nm.Parse("$GPRMC,110910.59,A,5047.3986,N,00054.6007,W,0.08,0.19,150920,0.24,W,D,V*75")

// To read the data for example current position ie the latest received during last 60 seconds

fmt.Println(nm.Get("position")
// 50° 47.3986'N, 000° 54.6007'W

// To access all data accumulated over last 60 seconds 
// ie stale items are deleted after 60 seconds and new data overwrites the last value

x:= nm.GetMap())

// then can use  x["position] to get position, x["datetime"] etc  ie the data obtained above is

fm.Println(x)
//map[datetime:2020-09-15T11:09:10.59+00:00 faa_mode:D fix_date:2020-09-15 fix_time:11:09:10.59 mag_var:-0.24 nav_status:V position:50° 47.3986'N, 000° 54.6007'W sog:0.08 status:A tmg:0.19]

// to make a new sentence from the data 

gprmc, _ := nm.WriteSentence("gp", "rmc")
fmt.Println(gprmc)
$GPRMC,110910.59,A,5047.3986,N,00054.6007,W,0.08,0.19,150920,0.24,W,D,V*75 

// note in this example same as parsed in but typical GPS units create multiple sentence types every second so
// in a real use case the generated message will have the very latest info from a combined source


```

The above will on first run create a nmea_sentences.yaml file in your current directory which will be used on the next run.
This provides an easy way customised configure the library but you can access the library without using DefaultConfig and
load config files from different locations or just use the in built ones.

The default data persistence is also configurable.

The data is always in readable string format and up to date so you can easily write logs to your own format or simply periodically
 write the map in JSON format to a log file.  There are functions to convert some data types such as position to float variables.

### Limitations

- No plans to support AIS
- Only supports comma delimited fields and messages starting with $
- Limited to passing sentences and fields which are fix format

## Full details

### Assuming you have installed go and are outside of GOPATH

```go
install go on your system

go mod init  your_module
go get github.com/martinmarsh/nmea0183

//write a main file as shown below or copy /example/main.go:

go install your_module
```

### Basic use

To start write the following in main.go in your modules root directory
a copy of this file is in demo/main.go

```go

    package main

    import (
     "fmt"
     "github.com/martinmarsh/nmea0183"
    )

    func main() {
        // For easy start we use built in sentences
        sentences := nmea0183.DefaultSentences()

        // Now create a handle to the structure used to hold sentence definitions and parsed data
        nm:= sentences.MakeHandle()

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
        fmt.Printf("%f %f\n",latFloat, longFloat)

        //Can write a variable from float lat and long
        nm.LatLongToString(latFloat, longFloat, "new_position")
        fmt.Println(nm.Data["new_position"])

        //examples of other sentences passed
        nm.Parse("$HCHDM,172.5,M*28")
        nm.Parse("$GPAPB,A,A,5,L,N,V,V,359.,T,1,359.1,T,6,T,A*7C")
        nm.Parse("$SSDPT,2.8,-0.7")

        //Data is built and updated as sentences are parsed
        fmt.Println(nm.GetMap())

        }
```

### Defing your own sentences and mapping to variable

Built in sentences and variables may be sufficient for some applications but sooner or later some
bespoke configuration will be required. If you don't need the flexibility of an external configuration
and are happy to build in the configuration into your compiled code you can parse your own sentence
definitions.  This is how:

```go

    // example of 2 user defined sentences, just list the variable names to collect the data in
    zda := []string {"time","day","month","year","tz"}
    dpt := []string {"dbt","toff"}

    // These definitions are placed in map with the key matching the sentence name 
    formats := map[string][]string {"zda": zda, "dpt": dpt}

    // Use create to load your only your sentences - no built in ones will be added
    sentences := MakeSentences(formats)

    // create handle
    nm := sentences.MakeHandle()

    // Now just parse sentences
    nm.Parse("$GPZDA,110910.59,15,09,2020,01,30*6F")
```

In the above example "zda" refers to any NMEA 0183 sentence after the first 2 digits, the manufacturer's code is
removed. The sentence definition refers to variables: "time","day","month","year","tz".  These variables
will use the built in definition.  When parsing a sentence each variable consumes the NMEA string in order ie
starting after the $GPZDA in our example "time" is taken from "110910.59", the "day" from "15", "month" from "09"
and "year" from "2020".  "tz" is an example of a variable which is taken from 2 fields ie 01,30
which is hours and minutes of time offset from UTC.  

The current built in sentence definition for ZDA is now uses a new format which allows a complete date to be read from 6 consecutive fields as defined in ZDA ie:
    zda := []string {"datetime"}
This parses just one user variable called datetime which is easier to handle but the old format is still supported.

Other common formats consuming multiple fields are:  position, lat, long, mag_var. Creating variables from multiple fields is safer.
Position for example holds both lat and long and takes 4 fields eg in "$GPRMC,110910.59,A,5047.3986,N,00054.6007,W,0.08,0.19,150920,0.24,W,D,V*75"  the sequence  "5047.3986,N,00054.6007,W" is parsed to position which would be set to "50° 47.3986'N, 000° 54.6007'W".
This ensures that the position comes from one just sentence and the parts cannot be mixed up.

### Defining your own variable

The above uses built in variable definitions but you can configure this too.
Here is an example of variables which you could set to use instead of the default ones.

```go
    // Lets use our own time variables d map them to sentences
    sentences := map[string][]string {
        "rmc": {"pos_time", "pos_status", "my_position", "sog", "tmg", "pos_date", "mag_var"},
        "zda": {"my_time_of_day"},
    }
    sentences := map[string][]string {"zda": zda, "rmc": rmc}

    // now define some variables we might like to use mapping them to internal templates
    variables := map[string]string {
        "pos_status": "A",
        "pos_time": "hhmmss.ss",
        "my_time_of_day": "hhmmss,day,month,year,tz",
        "my_position": "lat,NS,long,WE", //formatted lat, long
        "sog": "x.x",                 // Speed Over Ground 
        "tmg": "x.x",                 // Track Made Good
        "pos_date": "ddmmyy",
        "mag_var": "x.x,w",       // Mag Var E positive, W    
    }


        sentences := MakeSentences(formats, variables)

        // create handle
        nm := sentences.MakeHandle()
        nm := Create(sentences, variables)
        nm.Parse("$GPZDA,110910.59,15,09,2020,01,30*6F")
        nm.Parse("$GPRMC,110910.59,A,5047.3986,N,00054.6007,W,0.08,0.19,150920,0.24,W,D,V*75")
```

### Using a config file instead of building in sentence definitions

For more flexibility having to install
go and rebuild configurations can be read and written to nmea_sentences.yaml file in the working directory.
Instead of Create but use Load.

```go
    var sentences nmea0183.Sentences
    err := sentences.Load()
    handle.Parse("$GPZDA,110910.59,15,09,2020,00,00*6F")

//Can specify location of config file and type eg ymal or json

    err := sentences.Load(".", "filename", "ymal")

//A default Config file can be written using default settings:

    sentences.SaveDefault()
    or
    sentences.SaveDefault(".", "filename", "ymal")
```

### Cleaning up old data

By default Parse and Merge build Sentence data into a Go map called handle.Data

This allows a record of status to be easily maintained for example on a boat application
simply parsing all sentences as they are received allows the boat navigation status to be updated. This is
especially useful if the sentences are merged from different devices/inputs

However, there is a risk of being mislead by old data for example devices such as a GPS might loose signal
and the position might not be included in the sentence. The application will still see data but the position
may get out of date.

To reduce this risk the package notes the time of each variables and old ones deleted:

    handle.Preferences(seconds, real_time)   // seconds is time to keep variables or <= 0 is forever, real time mode true or false

use real_time = false if your system does not have a real time clock or historic data is being passed and of course the
messages contain a date.  Set seconds to zero or less to disable auto clear

This is not intended to cope with complete loss of connections when no sentences are Parsed ie for this
to work some sentences must still be Parsed or a Merge calls made.  

### Different channels

By choosing different definition files can use different handles to parse sentences differently. Filename1 may select different parts or names to filename

    var sentences1, sentences2 nmea0183.Sentence
    err1 := sentences1.Load(".", "filename1", "ymal") 
    err2 := sentences1.Load(".", "filename2", "ymal")
    chan1 := sentences1.MakeHandle()
    chan2 := sentences1.MakeHandle()

    chan1.Parse("$GPAPB,A,A,5,L,N,V,V,359.,T,1,359.1,T,6,T,A*7C")
    results1 := chan1.GetMap() 

    chan2.Parse("$GPAPB,A,A,5,L,N,V,V,359.,T,1,359.1,T,6,T,A*7C")
     results2 := chan2.GetMap() 

    for example different configs might result in:
    results1:
     map[acir:V aper:V bod:359. bod_true:T bpd:359.1 device:GP did:1 hts:6 hts_true:T sentence:apb status:A xte:L5 xte_units:N]

    results2:
    map[mybod:359. my_hts:6 sentence:apb  device:GP]

    # Sentences Defined in Default Config

## AAM - Waypoint Arrival Alarm

Actisense NGW-1 maps: 2000 -> 183, 0183 -> 2000

Indicates the status of arrival (entering the arrival circle, or passing the perpendicular of the course line) at the destination waypoint.

               1 2 3   4 5    6
               | | |   | |    |
        $--AAM,A,A,x.x,N,c--c*hh<CR><LF>

        aam:arrived_circle,passed_waypt,arrival_radius,radius_units,waypt_id

| Field | Name           | Parse Format | No of Fields used | Example Value |
| ----- | -------------- | ------------ | ----------------- | ------------- |
| 1     | arrived_circle | A            | 1                 | A             |
| 2     | passed_wpt     | A            | 1                 | V             |
| 3     | arrival_radius | x.x          | 1                 | 4.5           |
| 4     | radius_units   | A            | 1                 | N             |
| 5     | waypt_id       | A            | 1                 | Chich         |

**arrived_circle**:  Status A = Arrival Circle Entered,  V = not entered

**passed_way_pt**:  Status A = Perpendicular passed at waypoint,  V = not passed

**arrival_radius**: Arrival in circle radius

**radius_units**: Units of radius, N = nautical miles

**waypt_id**: Waypoint ID

Example: $GPAAM,A,A,0.10,N,WPTNME*43

WPTNME is the waypoint name.

---

## APA - Autopilot Sentence "A"

This sentence is sent by some GPS receivers to allow them to be used to control an autopilot unit. This sentence is commonly used by autopilots   and contains navigation receiver warning flag status, cross-track-error, waypoint arrival status, initial bearing from origin waypoint to the destination, continuous bearing from present position to destination and recommended heading-to-steer to destination waypoint for the active navigation leg of the journey.

               1 2  3   4 5 6 7  8  9 10    11
               | |  |   | | | |  |  | |     |
        $--APA,A,A,x.xx,L,N,A,A,xxx,M,c---c*hh<CR><LR>

Sentence def:

       APA: ap_status,ap_loran,xte,arrived_circle,passed_waypt,bearing_to_waypt,waypt_id

| Field | Name             | Format  | No of Fields matched | Example Value |
| ----- | ---------------- | ------- | -------------------- | ------------- |
| 1     | ap_status        | A       | 1                    | V             |
| 2     | ap_loran         | A       | 1                    | V             |
| 3,4,5 | xte              | x.x,R N | 3                    | L3.4N         |
| 6     | arrived_circle   | A       | 1                    | A             |
| 7     | passed_waypt     | A       | 1                    | A             |
| 8,9   | bearing_to_waypt | x.x,T   | 2                    | 123°M         |
| 10    | waypt_id         | A       | 1                    | Chich         |

**ap_status**:  Status V = Loran-C Blink or SNR warning V = general warning flag or other navigation systems when a reliable fix is not available

**ap_loran**: not used may see Status V  Loran-C Cycle Lock warning flag A = OK or not used

**xte**:  Cross Track Error eg R2.3N steer right  L2.3N is steer left, 2.3 is distance off course Units are N = Nautical miles or kilometers

**arrived_circle**: Status A = Arrival Circle Entered

**passed_waypt**: Status A = Perpendicular passed at waypoint

**bearing_to_waypt**: Bearing origin to destination. Post fix with M = Magnetic, T = True

**waypt_id**: Destination Waypoint ID

checksum automatic

Example Sentence: $GPAPA,A,A,8.30,L,M,V,V,11.7,T,Turning Track to Ijmuiden 1*1B"

---

## APB - Autopilot Sentence "B"

NGW-1: 2000 -> 183, 0183 -> 2000

This is a fixed form of the APA sentence with some ambiguities removed, used in later versions

Note: Some autopilots, Robertson in particular, misinterpret "bearing from origin to destination" as "bearing from present position to destination". This is likely due to the difference between the APB sentence and the APA sentence. for the APA sentence this would be the correct thing to do for the data in the same field. APA only differs from APB in this one field and APA leaves off the last two fields where this distinction is clearly spelled out. This will result in poor performance if the boat is sufficiently off-course that the two bearings are different. 13 15

               1 2 3   4 5 6 7 8   9 10   11  12|   14|
               | | |   | | | | |   | |    |   | |   | |
        $--APB,A,A,x.x,a,N,A,A,x.x,a,c--c,x.x,a,x.x,a*hh<CR><LF>

Sentence def:

       apb: ap_status,ap_loran,xte,arrived_circle,passed_waypt,bearing_origin_to_waypt,waypt_id,bearing_position_to_waypt,hts,ap_mode

| Field  | Name                      | Format  | No of Fields matched | Example Value |
| ------ | ------------------------- | ------- | -------------------- | ------------- |
| 1      | ap_status                 | A       | 1                    | V             |
| 2      | ap_loran                  | A       | 1                    | V             |
| 3,4,5  | xte                       | x.x,R,N | 3                    | L3.4N         |
| 6      | arrived_circle            | A       | 1                    | A             |
| 7      | passed_waypt              | A       | 1                    | A             |
| 8,9    | bearing_origin_to_waypt   | x.x,T   | 2                    |               |
| 10     | waypt_id                  | str     | 1                    | Chich         |
| 11, 12 | bearing_position_to_waypt | x.x,T   | 2                    | 123°T         |
| 13, 14 | hts                       | x.x,T   | 2                    | 121°M         |
| 15     | ap_mode                   | A       | 1                    |               |

**ap_status**: Status A = Data valid V = Loran-C Blink or SNR warning V = general warning flag or other navigation systems when a reliable fix is not available

**ap_loran**: Status V = Loran-C Cycle Lock warning flag A = OK or not used

**xte**: Cross Track Error Magnitude, Direction to steer, L or R, Cross Track Units, N = Nautical Miles example L12.3N

**arrived_circle**: Status A = Arrival Circle Entered

**passed_waypt**: Status A = Perpendicular passed at waypoint

**bearing_origin_to_waypt**: Bearing origin to destination M = Magnetic, T = True

**waypt_id**: Destination Waypoint ID

**bearing_position_to_waypt**: Bearing, present position to Destination M = Magnetic, T = True

**hts**: Heading to steer to destination waypoint M = Magnetic, T = True

**ap_mode**: Mode indicator ('D' when the position used for the XTE is valid, otherwise 'E')

Checksum

Examples:

       $GPAPB,A,A,0.00536,R,N,V,V,210.0,T,Vlissingen,213.4,T,213.4,T,D*5F
       $GPAPB,A,A,5,L,N,V,V,359.,T,1,359.1,T,6,T,A*7C
---

## RMC - Recommended Minimum Navigation Information

Actisense NGW-1 maps: 2000 -> 183, 0183 -> 2000

This is one of the sentences commonly emitted by GPS units. The default config writes NMEA 4.1 standard but will read in the older versions

              1         2 3       4 5        6  7   8   9    10 11
              |         | |       | |        |  |   |   |    |  |
       $--RMC,hhmmss.ss,A,ddmm.mm,a,dddmm.mm,a,x.x,x.x,xxxx,x.x,a*hh<CR><LF>
       NMEA 2.3:
       $--RMC,hhmmss.ss,A,ddmm.mm,a,dddmm.mm,a,x.x,x.x,xxxx,x.x,a,m*hh<CR><LF>
       NMEA 4.1:
       $--RMC,hhmmss.ss,A,ddmm.mm,a,dddmm.mm,a,x.x,x.x,xxxx,x.x,a,m,s*hh<CR><LF>

 Sentence def:
       "rmc": {"time", "status", "position", "sog", "tmg", "date", "mag_var", "faa_mode,"nav_status"},

| Field   | Name       | Format         | No of Fields matched | Example Value                 |
| ------- | ---------- | -------------- | -------------------- | ----------------------------- |
| 1       | time       | hhmmss.ss      | 1                    | V                             |
| 2       | gps_status | A              | 1                    | V                             |
| 3,4,5,6 | position   | lat,NS,long,WE | 4                    | 50° 10.3986'N, 000° 54.6007'W |
| 7       | sog        | x.x            | 1                    | 4.3                           |
| 8       | tmg        | x.x            | 1                    | 121                           |
| 9       | date       | ddmmyy         | 2                    | 2020-09-15                    |
| 10, 11  | mag_var    | x.x,w          | 2                    | -1.4                          |
| 12      | faa_mode   | A              | 1                    | A                             |
| 13      | nav_status | A              | 1                    | A                             |

**time**: UTC of position fix, hh is hours, mm is minutes, ss.ss is seconds.

**gps_status**: Status, A = Valid, V = Warning

**position**: Latitude, dd is degrees. mm.mm is minutes. N or S Longitude, ddd is degrees. mm.mm is minutes. E or W

**sog**: Speed over ground, knots

**tmg** Track made good, degrees true

**date** Date, ddmmyy

**mag_var** Magnetic Variation, degrees E or W

**faa_mode** FAA mode indicator (NMEA 2.3 and later)

**nav_status**: Nav Status (NMEA 4.1 and later) A=autonomous, D=differential, E=Estimated, M=Manual input mode N=not valid, S=Simulator, V = Valid

Checksum

A status of V means the GPS has a valid fix that is below an internal quality threshold, e.g. because the dilution of precision is too high or an elevation mask test failed.

The number of digits past the decimal point for Time, Latitude and Longitude is model dependent.

Example: $GNRMC,001031.00,A,4404.13993,N,12118.86023,W,0.146,,100117,,,A*7B

---

## ZDA - Time & Date - UTC, day, month, year and local time zone

Actisense NGW-1 maps: 2000 -> 183, 0183 -> 2000

    This is one of the sentences commonly emitted by GPS units.

                    1         2  3  4    5  6  7
                    |         |  |  |    |  |  |
            $--ZDA,hhmmss.ss,xx,xx,xxxx,xx,xx*hh<CR><LF>

    Sentence def:
       "zda": {"datetime"} 

**datetime**  Date time zone string as defined in RFC3339 Format

Alternative could be defined Sentence defined:

       "zda": {"time", "day", "month", "year", "tz"}

| Field | Name  | Format    | No of Fields matched | Example Value |
| ----- | ----- | --------- | -------------------- | ------------- |
| 1     | time  | hhmmss.ss | 1                    | 16:15:40.12   |
| 2     | day   | x         | DD_day               | 25            |
| 3     | month | x         | DD_month             | 12            |
| 4     | year  | x         | DD_year              | 2021          |
| 5,6   | tz    | tz_h,tz_m | 2                    | -12:23        |

**time** UTC time (hours, minutes, seconds, may have fractional sub-seconds)

**day** Day, 01 to 31

**month** Month, 01 to 12

**year** Year (4 digits)

**tz** Local zone description, 00 to +- 13 hours  Local zone minutes description, 00 to 59, apply same sign as local hours

Checksum

Example: $GPZDA,160012.71,11,03,2004,-1,00*7D

---
