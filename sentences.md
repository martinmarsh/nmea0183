
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

This is one of the sentences commonly emitted by GPS units.

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