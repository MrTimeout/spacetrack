package model

import "github.com/MrTimeout/spacetrack/utils"

const (
	BY_CCSDS_OMM_VERS      string = "CCSDS_OMM_VERS"
	BY_COMMENT             string = "COMMENT"
	BY_CREATION_DATE       string = "CREATION_DATE"
	BY_ORIGINATOR          string = "ORIGINATOR"
	BY_OBJECT_NAME         string = "OBJECT_NAME"
	BY_OBJECT_ID           string = "OBJECT_ID"
	BY_CENTER_NAME         string = "CENTER_NAME"
	BY_REF_FRAME           string = "REF_FRAME"
	BY_TIME_SYSTEM         string = "TIME_SYSTEM"
	BY_MEAN_ELEMENT_THEORY string = "MEAN_ELEMENT_THEORY"
	BY_EPOCH               string = "EPOCH"
	BY_MEAN_MOTION         string = "MEAN_MOTION"
	BY_ECCENTRICITY        string = "ECCENTRICITY"
	BY_INCLINATION         string = "INCLINATION"
	BY_RA_OF_ASC_NODE      string = "RA_OF_ASC_NODE"
	BY_ARG_OF_PERICENTER   string = "ARG_OF_PERICENTER"
	BY_MEAN_ANOMALY        string = "MEAN_ANOMALY"
	BY_EPHEMERIS_TYPE      string = "EPHEMERIS_TYPE"
	BY_CLASSIFICATION_TYPE string = "CLASSIFICATION_TYPE"
	BY_NORAD_CAT_ID        string = "NORAD_CAT_ID"
	BY_ELEMENT_SET_NO      string = "ELEMENT_SET_NO"
	BY_REV_AT_EPOCH        string = "REV_AT_EPOCH"
	BY_BSTAR               string = "BSTAR"
	BY_MEAN_MOTION_DOT     string = "MEAN_MOTION_DOT"
	BY_MEAN_MOTION_DDOT    string = "MEAN_MOTION_DDOT"
	BY_SEMIMAJOR_AXIS      string = "SEMIMAJOR_AXIS"
	BY_PERIOD              string = "PERIOD"
	BY_APOAPSIS            string = "APOAPSIS"
	BY_PERIAPSIS           string = "PERIAPSIS"
	BY_OBJECT_TYPE         string = "OBJECT_TYPE"
	BY_RCS_SIZE            string = "RCS_SIZE"
	BY_COUNTRY_CODE        string = "COUNTRY_CODE"
	BY_LAUNCH_DATE         string = "LAUNCH_DATE"
	BY_SITE                string = "SITE"
	BY_DECAY_DATE          string = "DECAY_DATE"
	BY_FILE                string = "FILE"
	BY_GP_ID               string = "GP_ID"
	BY_TLE_LINE0           string = "TLE_LINE0"
	BY_TLE_LINE1           string = "TLE_LINE1"
	BY_TLE_LINE2           string = "TLE_LINE2"
)

var ByPossibleValues = []string{
	BY_CCSDS_OMM_VERS,
	BY_COMMENT,
	BY_CREATION_DATE,
	BY_ORIGINATOR,
	BY_OBJECT_NAME,
	BY_OBJECT_ID,
	BY_CENTER_NAME,
	BY_REF_FRAME,
	BY_TIME_SYSTEM,
	BY_MEAN_ELEMENT_THEORY,
	BY_EPOCH,
	BY_MEAN_MOTION,
	BY_ECCENTRICITY,
	BY_INCLINATION,
	BY_RA_OF_ASC_NODE,
	BY_ARG_OF_PERICENTER,
	BY_MEAN_ANOMALY,
	BY_EPHEMERIS_TYPE,
	BY_CLASSIFICATION_TYPE,
	BY_NORAD_CAT_ID,
	BY_ELEMENT_SET_NO,
	BY_REV_AT_EPOCH,
	BY_BSTAR,
	BY_MEAN_MOTION_DOT,
	BY_MEAN_MOTION_DDOT,
	BY_SEMIMAJOR_AXIS,
	BY_PERIOD,
	BY_APOAPSIS,
	BY_PERIAPSIS,
	BY_OBJECT_TYPE,
	BY_RCS_SIZE,
	BY_COUNTRY_CODE,
	BY_LAUNCH_DATE,
	BY_SITE,
	BY_DECAY_DATE,
	BY_FILE,
	BY_GP_ID,
	BY_TLE_LINE0,
	BY_TLE_LINE1,
	BY_TLE_LINE2,
}

type OrderBy struct {
	By   string
	Sort Sort
}

func (o OrderBy) ToPath() string {
	if !utils.CheckStringsIsIn(ByPossibleValues, o.By) {
		return ""
	}
	var orderByQuery = "/orderby/" + o.By

	if sort := o.Sort.String(); sort != "" {
		orderByQuery += " " + sort
	}

	return orderByQuery
}
