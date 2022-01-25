# nmea0183

## Post build configurable NMEA 0183 sentence analyser and writer in a go package

### Reason to choose

- Configurable Sentences definitions
- Parsed values stored can be stored in user named variables
- Collection of pre-defined common sentences and variables for quick start
- Same configuration allows parsing and writing of defined sentences
- Optional use of yaml, json or toml config file to externally change sentence and/or variable definitions post build
- Collects data across multiple sentences to merge into one results set
- Readable values extracted eg position = "50° .3986'N, 000° 54.6007'W
- Conversion funtions to float values
- Handles Removing old data

### Reason to consider other packages

- Want a solution that collects data into pre-defined structures better tailored to your application
- Need to parse more  multi-line sentences

### Current Status

Essentially functional but undergoing development and testing. Some experimental features may be added or removed as the package is being developed for use in my navigation system.
Features marked with * are under development and refinement, templates are being refined as further sentences tested.

### Features

- Sentences fully customisable and configurable via a Yaml/Json files
- Built in basic sentences can be used if preferred
- Check sum verification automatic if present
- Sentence definitions and config can be preloaded and multiple se47ntences passed
- Results returned in readable string format eg bearing = 100°T  position = "50° 00.3986'N, 000° 54.6007'W"
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
| 5     | waypt_id       | str          | 1                 | Chich         |

**arrived_circle**:  Status A = Arrival Circle Entered,  V = not enterred

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

| Field | Name             | Format      | No of Fields matched | Example Value |
| ----- | ---------------- | ----------- | -------------------- | ------------- |
| 1     | ap_status        | A           | 1                    | V             |
| 2     | ap_loran         | A           | 1                    | V             |
| 3,4,5 | xte              | Lx.xN, R, N | 3                    | L3.4N         |
| 6     | arrived_circle   | A           | 1                    | A             |
| 7     | passed_waypt     | A           | 1                    | A             |
| 8,9   | bearing_to_waypt | x.xT, T     | 2                    | 123°M          |
| 10    | waypt_id         | str         | 1                    | Chich         |

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

| Field  | Name                      | Format      | No of Fields matched | Example Value |
| ------ | ------------------------- | ----------- | -------------------- | ------------- |
| 1      | ap_status                 | A           | 1                    | V             |
| 2      | ap_loran                  | A           | 1                    | V             |
| 3,4,5  | xte                       | Lx.xN, R, N | 3                    | L3.4N         |
| 6      | arrived_circle            | A           | 1                    | A             |
| 7      | passed_waypt              | A           | 1                    | A             |
| 8,9    | bearing_origin_to_waypt   | x.xT, T     | 2                    |               |
| 10     | waypt_id                  | str         | 1                    | Chich         |
| 11, 12 | bearing_position_to_waypt | x.xT, T     | 2                    | 123°T         |
| 13, 14 | hts                       | x.xT, T     | 2                    | 121°M         |
| 15     | ap_mode                   | A           | 1                    |               |

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


| Field   | Name       | Format                     | No of Fields matched | Example Value                 |
| ------- | ---------- | -------------------------- | -------------------- | ----------------------------- |
| 1       | time       | hhmmss.ss                  | 1                    | V                             |
| 2       | gps_status | A                          | 1                    | V                             |
| 3,4,5,6 | position   | lat,lat_NS,pos_long,pos_WE | 4                    | 50° 10.3986'N, 000° 54.6007'W |
| 7       | sog        | x.x                        | 1                    | 4.3                           |
| 8       | tmg        | x.x                        | 1                    | 121                           |
| 9       | date       | ddmmyy                     | 2                    | 2020-09-15                    |
| 10, 11  | mag_var    | x.x, w                     | 2                    | -1.4                          |
| 12      | faa_mode   | A                          | 1                    | A                             |
| 13      | nav_status | A                          | 1                    | A                             |

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
       "zda": {"time", "day", "month", "year", "tz"}


| Field | Name  | Format     | No of Fields matched | Example Value |
| ----- | ----- | ---------- | -------------------- | ------------- |
| 1     | time  | hhmmss.ss  | 1                    | 16:15:40.12   |
| 2     | day   | x          | DD_day               | 25            |
| 3     | month | x          | DD_month             | 12            |
| 7     | year  | x          | DD_year              | 2021          |
| 8     | tz    | tz_h, tz:m | 2                    | 12:23         |


**time** UTC time (hours, minutes, seconds, may have fractional subseconds)

**day** Day, 01 to 31

**month** Month, 01 to 12

**year** Year (4 digits)

**tz** Local zone description, 00 to +- 13 hours  Local zone minutes description, 00 to 59, apply same sign as local hours

Checksum

Example: $GPZDA,160012.71,11,03,2004,-1,00*7D

---