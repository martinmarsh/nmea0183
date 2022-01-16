package nmea0183

func GetDefaultVars() *map[string][]string{
    vars := map[string][]string {
        "time": {"hhmmss.ss"},
        "status": {"A"},                 // status of fix A = ok ie 1 V = fail ie 0
        "lat": {"lat", "lat_NS"},      // formated latitude
        "long": {"long","Long_WE"},    // formated longitude
        "position": {"lat", "lat_NS", "pos_long", "pos_WE"}, //formated lat, long
        "sog": {"x.x"},                 // Speed Over Ground  float knots
        "tmg": {"x.x"},                 // Track Made Good
        "date": {"ddmmyy"},
        "mag_var": {"x.x", "w"},   // Mag Var E positive, W negative
        "day": {"DD_day"},
        "month": {"DD_month"},
        "year": {"DD_year"},
        "tz":  {"tz_h", "tz:m"},   // Datetime from ZDA if available - tz:m returns hrs:mins
        "tzhrs": {"tz_h", "tz_m"},   // Datetime from ZDA if available - tz_m returns decimal hours as a float
        "xte": {"x.x", "R"},      // Cross Track Error turn R or L eg prefix L12.3
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