package nmea0183

func GetDefaultVars() *map[string][]string{
    vars := map[string][]string {
        "arrived_circle": {"A"},
        "passed_waypt": {"A"},
        "arrival_radius": {"x.x"},
        "radius_units":{"A"},
        "waypt_id": {"str"},
        "ap_status": {"A"},
        "ap_loran": {"A"},
        "bearing_to_waypt": {"x.xT","T"},
        "bearing_origin_to_waypt": {"x.xT","T"},
        "bearing_position_to_waypt": {"x.xT","T"},
        "hts": {"x.xT","T"},    // Heading to Steer T True or M magnetic
        "ap_mode": {"A"},
        "faa_mode": {"A"},
        "nav_status": {"A"},

        "time": {"hhmmss.ss"},
        "status": {"A"},                 // status of fix A = ok ie 1 V = fail ie 0
        "lat": {"lat", "lat_NS"},      // formated latitude
        "long": {"long","long_WE"},    // formated longitude
        "position": {"lat", "lat_NS", "pos_long", "pos_WE"}, //formated lat, long
        "sog": {"x.x"},                 // Speed Over Ground  float knots
        "tmg": {"x.x"},                 // Track Made Good
        "date": {"ddmmyy"},
        "mag_var": {"x.x", "w"},       // Mag Var E positive, W negative
        "day": {"DD_day"},
        "month": {"DD_month"},
        "year": {"DD_year"},
        "tz": {"tz_h", "tz:m"},   // Datetime from ZDA if available - tz:m returns minutes part of tx as hh:mm format
        "xte": {"Lx.xN", "R", "N"},      // Cross Track Error turn R or L eg prefix L12.3N post fix  N = Nm
        "acir": {"A"},           // Arrived at way pt circle
        "aper": {"A"},           // Perpendicular passing of way pt
        "bod": {"x.x"},           // Bearing origin to destination
        "bod_true": {"T"},        // Bearing origin to destination True
        "did": {"str"},           //Destination Waypoint ID as a str
        "bpd": {"x.x"},
        "bdp_true": {"T"},        // Bearing, present position to Destination True
           // Heading to Steer True
        "hdm": {"x.x"},          // Heading Magnetic
        "hdm_true": {"T"},      // heading true or Manetic
        "dbt": {"x.x"},          // Depth below transducer
        "toff": {"-x.x"},         // Transducer offset -ve from transducer to keel +ve transducer to water line
        "stw": {"x.x"},          // Speed Through Water float knots
        "dw":  {"x.x"},          // Water distance since reset float knots
       
    }

    return &vars
}

func GetDefaultSentences() *map[string][]string{
    sent := map[string][]string {
        "aam": {"arrived_circle", "passed_waypt", "arrival_radius", "radius_units", "waypt_id"},
        "apa": {"ap_status","ap_loran", "xte", "arrived_circle", "passed_waypt", "bearing_to_waypt", "waypt_id"},
        "apb": {"ap_status", "ap_loran", "xte", "arrived_circle", "passed_waypt", "bearing_origin_to_waypt", "waypt_id", "bearing_position_to_waypt", "hts", "ap_mode"},
        "rmc": {"time", "status", "position", "sog", "tmg", "date", "mag_var", "faa_mode","nav_status"},
       
        "zda": {"time", "day", "month", "year", "tz"},
        "hdg": {"n/a", "n/a", "n/a", "mag_var"},
        "hdm": {"hdm", "hdm_true"}, 
        "dpt": {"dbt", "toff"},
        "vhm": {"n/a", "n/a", "n/a", "n/a", "stw"},
        "vlw": {"n/a", "n/a", "wd"},
    }

    return &sent
}