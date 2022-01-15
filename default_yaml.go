package nmea0183

/*
The following is a get started config. used to generate nmea_conf.yaml
if it is missing in root directory

*/


var yamlExample = []byte(`
sentences:
    RMC:
        - "time"
        - "status"
        - "lat"
        - "long"
        - "sog"
        - "tmg"
        - "date"
        - "mag_var"
    ZDA:
        - "time"
        - "day"
        - "month"
        - "year"
        - "tz"
    APB:
        - "status"
        - "n/a"
        - "xte"
        - "xte_units"
        - "acir"
        - "aper"
        - "bod"
        - "bod_true"
        - "did"
        - "bpd"
        - "bpd_true"
        - "hts"
        - "hts_true"
    HDG:
        - "n/a"
        - "n/a" 
        - "n/a"
        - "mag_var"
    HDM:
        - "hdm"
    DPT:
        - "dbt"
        - "toff"
    VHW:
        - "n/a"
        - "n/a"
        - "n/a"
        - "n/a"
        - "stw"
    VLW:
        - "n/a"
        - "n/a"
        - "wd"

variables:
    time: "hhmmss.ss"  # time of fix
    status: "A"   # status of fix A = ok ie 1 V = fail ie 0
    lat:
        - "llll.llll"
        - "NS"  # lat N / S postfix
    long:
        - "yyyyy.yyyy"
        - "WE"   # long float W/E postfix
    sog: "x.x"  # Speed Over Ground  float knots
    tmg: "x.x"  # Track Made Good
    date: "ddmmyy" # Date of fix may not be valid with some GPS
    mag_var:
        - "x.x"
        - "w"     # Mag Var E positive, W negative
    day: "x"
    month: "x"
    year: "x"
    tz:
        - "tz_h"
        - "tz_m"  # Datetime from ZDA if available - tz return hours and mins as a float
    xte:
        - "x.x"
        - "R"           # Cross Track Error turn R or L
        
    xte_units: "A"      # Units for XTE - N = Nm
    acir: "A"           # Arrived at way pt circle
    aper: "A"           # Perpendicular passing of way pt
    bod: "x.x"
    bod_true: "T"        # Bearing origin to destination True
    did: "str"             # Destination Waypoint ID as a str
    bpd: "x.x"
    bdp_true: "T"        # Bearing, present position to Destination True
    hts: "x.x"
    hts_true: "T"        # Heading to Steer True
    hdm: "x.x"          # Heading Magnetic
    dbt: "x.x"          # Depth below transducer
    toff: "x.x"         # Transducer offset -ve from transducer to keel +ve transducer to water line
    stw: "x.x"          # Speed Through Water float knots
    dw:  "x.x"             # Water distance since reset float knots

`)

func GetDefaultVars() *map[string][]string{
    vars := map[string][]string {
        "time": {"hhmmss.ss"},
        "status": {"A"},                 // status of fix A = ok ie 1 V = fail ie 0
        "lat": {"llll.lll", "NS"},      // formated latitude
        "long": {"yyyyy.yyyy","WE"},    // formated longitude
        "position": {"llll.lll", "NS", ",yyyyy.yyyy", "WE"}, //formated lat, long
        "sog": {"x.x"},                 // Speed Over Ground  float knots
        "tmg": {"x.x"},                 // Track Made Good
        "date": {"ddmmyy"},
        "mag_var": {"x.x", "w"},   // Mag Var E positive, W negative
        "day": {"x"},
        "month": {"x"},
        "year": {"x"},
        "tz": {"tz_h", "tz_m"},   // Datetime from ZDA if available - tz returns decimal hours as a float
        "xte": {"x.x", "R"},      // Cross Track Error turn R or L
        "xte_units": {"A"},       // Units for XTE - N = Nm
        "acir": {"A"},           // Arrived at way pt circle
        "aper": {"A"},           // Perpendicular passing of way pt
        "bod": {"x.x"},           // Bearing origin to destination
        "bod_true": {"T"},        // Bearing origin to destination True
        "did": {"str"},           //Destination Waypoint ID as a str
        "bpd": {"x.x"},
        "bdp_true": {"T"},        // Bearing, present position to Destination True
        "hts": {"x.x"},
        "hts_true": {"T"},        // Heading to Steer True
        "hdm": {"x.x"},          // Heading Magnetic
        "dbt": {"x.x"},          // Depth below transducer
        "toff": {"x.x"},         // Transducer offset -ve from transducer to keel +ve transducer to water line
        "stw": {"x.x"},          // Speed Through Water float knots
        "dw":  {"x.x"},          // Water distance since reset float knots
       
    }

    return &vars
}

func GetDefaultSentences() *map[string][]string{
    sent := map[string][]string {
        "rmc": {"time", "status", "position", "sog", "tmg", "date", "mag_var"},
        "zda": {"time", "day", "month", "year", "tz"},
        "apb": {"status","n/a", "xte","xte_units","acir", "aper", "bod", "bod_true", "did", "bpd", "bpd_true", "hts","hts_true"},
        "hdg": {"n/a", "n/a", "n/a", "mag_var"},
        "hdm": {"hdm"}, 
        "dpt": {"dbt", "toff"},
        "vhm": {"n/a", "n/a", "n/a", "n/a", "stw"},
        "vlw": {"n/a", "n/a", "wd"},
    }

    return &sent
}