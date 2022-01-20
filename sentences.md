
# NMEA Sentences

## Multi-line

ALR (gps_almanac_data): not NGW-1, not supported in v1

RTE (Routes): not NGW-1, not supported in v1

---

## Messages with variable content

---

## Actisense NGW-1 0183 to 2000

AAM - Waypoint Arrival Alarm

APB - Autopilot Sentence "B"


---

## Actisense NGW-1 2000 to 0183

AAM - Waypoint Arrival Alarm

APB - Autopilot Sentence "B"


TW,
HDG, HDM, HDT
THS,
VDR, VHW, VDM
VDO ABM, BBM,ABM, BBM VDM, VDO, VDM
BWC, BWR, RMB, XTE
AAM,  BWC, BWR, RMB, ABM
GSA,
GRS, GSV,
GST,
MDA, MWD, MWV, VWR
MDA, MTW,
RPM,
VHW,
DBT, DPT,
DTM,
VLW,
GGA, GLL, GNS, RMC,
VTG,
GGA, GLL, GNS, GRS, GSA, RMC, ZDA,
GRS, ZDA,
RSA,
HDG, HDM or HDT*, VHW,
ROT,
HDG,
GRS, ZDA,
VBW

---

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

### aam:arrived_circle

1 Status A = Arrival Circle Entered,  V = not enterred

### passed_way_pt

2 Status A = Perpendicular passed at waypoint,  V = not passed

### arrival_radius

3 Arrival circle radius

### radius_units

4 Units of radius, nautical miles

### waypt_id

5 Waypoint ID

6  Checksum

Example: $GPAAM,A,A,0.10,N,WPTNME*43

WPTNME is the waypoint name.

---

## APA - Autopilot Sentence "A"

This sentence is sent by some GPS receivers to allow them to be used to control an autopilot unit. This sentence is commonly used by autopilots   and contains navigation receiver warning flag status, cross-track-error, waypoint arrival status, initial bearing from origin waypoint to the destination, continuous bearing from present position to destination and recommended heading-to-steer to destination waypoint for the active navigation leg of the journey.

               1 2  3   4 5 6 7  8  9 10    11
               | |  |   | | | |  |  | |     |
        $--APA,A,A,x.xx,L,N,A,A,xxx,M,c---c*hh<CR><LR>

Sentance def:

       APA: ap_status,n/a,xte,arrived_cicle,passed_wpt,bear_to_waypt,to_waypt_id

| Field | Name          | Format | No of Fields matched | Example Value |
| ----- | ------------- | ------ | -------------------- | ------------- |
| 1     | ap_status     | A      | 1                    | V             |
| 2     | n/a           |        | 1                    |               |
| 3,4   | xte           | x.x    | 3                    | L3.4N         |
| 6     | arrived_cicle | A      | 1                    | A             |
| 7     | passed_wpt    |        | 2                    | A             |
| 8,9   | bear_to_waypt |        | 2                    | 123M          |
| 10    | waypt_id      |        | 1                    | Chich         |

### ap_status:

Status V = Loran-C Blink or SNR warning V = general warning flag or other navigation systems when a reliable fix is not available

### n/a:

field 2 not mapped n/a -  means Status V  Loran-C Cycle Lock warning flag A = OK or not used

### xte

uses fields 3, 4 and 5 to produce Cross Track Error Magnitude  prefixed with L or R eg R2.3N steer right 2.3 N milesoff course see units.   Cross Track Units are Nautical miles or kilometers

### arrived_circle

uses field 6 Status A = Arrival Circle Entered

### passed_way_pt

uses field 7 Status A = Perpendicular passed at waypoint

### bear_to_waypt

uses fields 8 and 9 Bearing origin to destination. Post fix with 
M = Magnetic, T = True

## to_waypt_id

uses field 10 Destination Waypoint ID

field 11 checksum automatic

Example Sentence: $GPAPA,A,A,0.10,R,N,V,V,011,M,DEST,011,M*82

---

## APB - Autopilot Sentence "B"

NGW-1: 2000 -> 183, 0183 -> 2000


This is a fixed form of the APA sentence with some ambiguities removed.

Note: Some autopilots, Robertson in particular, misinterpret "bearing from origin to destination" as "bearing from present position to destination". This is likely due to the difference between the APB sentence and the APA sentence. for the APA sentence this would be the correct thing to do for the data in the same field. APA only differs from APB in this one field and APA leaves off the last two fields where this distinction is clearly spelled out. This will result in poor performance if the boat is sufficiently off-course that the two bearings are different. 13 15

               1 2 3   4 5 6 7 8   9 10   11  12|   14|
               | | |   | | | | |   | |    |   | |   | |
        $--APB,A,A,x.x,a,N,A,A,x.x,a,c--c,x.x,a,x.x,a*hh<CR><LF>

Field Number:

Status A = Data valid V = Loran-C Blink or SNR warning V = general warning flag or other navigation systems when a reliable fix is not available

Status V = Loran-C Cycle Lock warning flag A = OK or not used

Cross Track Error Magnitude

Direction to steer, L or R

Cross Track Units, N = Nautical Miles

Status A = Arrival Circle Entered

Status A = Perpendicular passed at waypoint

Bearing origin to destination

M = Magnetic, T = True

Destination Waypoint ID

Bearing, present position to Destination

M = Magnetic, T = True

Heading to steer to destination waypoint

M = Magnetic, T = True

Checksum

Example: $GPAPB,A,A,0.10,R,N,V,V,011,M,DEST,011,M,011,M*82


---

BOD - Bearing - Waypoint to Waypoint
        1   2 3   4 5    6    7
        |   | |   | |    |    |
 $--BOD,x.x,T,x.x,M,c--c,c--c*hh<CR><LF>
Field Number:

Bearing Degrees, True

T = True

Bearing Degrees, Magnetic

M = Magnetic

Destination Waypoint

origin Waypoint

Checksum

Example 1: $GPBOD,099.3,T,105.6,M,POINTB*01

Waypoint ID: "POINTB" Bearing 99.3 True, 105.6 Magnetic This sentence is transmitted in the GOTO mode, without an active route on your GPS. WARNING: this is the bearing from the moment you press enter in the GOTO page to the destination waypoint and is NOT updated dynamically! To update the information, (current bearing to waypoint), you will have to press enter in the GOTO page again.

Example 2: $GPBOD,097.0,T,103.2,M,POINTB,POINTA*52

This sentence is transmitted when a route is active. It contains the active leg information: origin waypoint "POINTA" and destination waypoint "POINTB", bearing between the two points 97.0 True, 103.2 Magnetic. It does NOT display the bearing from current location to destination waypoint! WARNING Again this information does not change until you are on the next leg of the route. (The bearing from POINTA to POINTB does not change during the time you are on this leg.)

This sentence has been replaced by BWW in NMEA 4.00 (and possibly earlier versions) [ANON].

BWC - Bearing & Distance to Waypoint - Great Circle
                                                         12
        1         2       3 4        5 6   7 8   9 10  11|   13
        |         |       | |        | |   | |   | |   | |    |
 $--BWC,hhmmss.ss,llll.ll,a,yyyyy.yy,a,x.x,T,x.x,M,x.x,N,c--c*hh<CR><LF>
NMEA 2.3:
 $--BWC,hhmmss.ss,llll.ll,a,yyyyy.yy,a,x.x,T,x.x,M,x.x,N,c--c,m*hh<CR><LF>
Field Number:

UTC Time or observation, hh is hours, mm is minutes, ss.ss is seconds

Waypoint Latitude

N = North, S = South

Waypoint Longitude

E = East, W = West

Bearing, degrees True

T = True

Bearing, degrees Magnetic

M = Magnetic

Distance, Nautical Miles

N = Nautical Miles

Waypoint ID

FAA mode indicator (NMEA 2.3 and later, optional)

Checksum

Example 1: $GPBWC,081837,,,,,,T,,M,,N*13

Example 2: GPBWC,220516,5130.02,N,00046.34,W,213.8,T,218.0,M,0004.6,N,EGLM*11

BWR - Bearing and Distance to Waypoint - Rhumb Line
                                                       11
        1         2       3 4        5 6   7 8   9 10  | 12  13
        |         |       | |        | |   | |   | |   | |    |
 $--BWR,hhmmss.ss,llll.ll,a,yyyyy.yy,a,x.x,T,x.x,M,x.x,N,c--c*hh<CR><LF>
NMEA 2.3:
 $--BWR,hhmmss.ss,llll.ll,a,yyyyy.yy,a,x.x,T,x.x,M,x.x,N,c--c,m*hh<CR><LF>
Field Number:

UTC Time of observation, hh is hours, mm is minutes, ss.ss is seconds

Waypoint Latitude

N = North, S = South

Waypoint Longitude

E = East, W = West

Bearing, degrees True

T = True

Bearing, degrees Magnetic

M = Magnetic

Distance, Nautical Miles

N = Nautical Miles

Waypoint ID

FAA mode indicator (NMEA 2.3 and later, optional)

Checksum

BWW - Bearing - Waypoint to Waypoint
Bearing calculated at the FROM waypoint.

        1   2 3   4 5    6    7
        |   | |   | |    |    |
 $--BWW,x.x,T,x.x,M,c--c,c--c*hh<CR><LF>
Field Number:

Bearing, degrees True

T = True

Bearing Degrees, Magnetic

M = Magnetic

TO Waypoint ID

FROM Waypoint ID

Checksum

DBK - Depth Below Keel
        1   2 3   4 5   6 7
        |   | |   | |   | |
 $--DBK,x.x,f,x.x,M,x.x,F*hh<CR><LF>
Field Number:

Depth, feet

f = feet

Depth, meters

M = meters

Depth, Fathoms

F = Fathoms

Checksum

DBS - Depth Below Surface
        1   2 3   4 5   6 7
        |   | |   | |   | |
 $--DBS,x.x,f,x.x,M,x.x,F*hh<CR><LF>
Field Number:

Depth, feet

f = feet

Depth, meters

M = meters

Depth, Fathoms

F = Fathoms

Checksum

DBT - Depth below transducer
        1   2 3   4 5   6 7
        |   | |   | |   | |
 $--DBT,x.x,f,x.x,M,x.x,F*hh<CR><LF>
Field Number:

Water depth, feet

f = feet

Water depth, meters

M = meters

Water depth, Fathoms

F = Fathoms

Checksum

In real-world sensors, sometimes not all three conversions are reported. So you might see something like $SDDBT,,f,22.5,M,,F*cs

Example: $SDDBT,7.8,f,2.4,M,1.3,F*0D


