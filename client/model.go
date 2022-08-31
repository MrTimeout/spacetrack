package client

import (
	"encoding/xml"

	"github.com/MrTimeout/spacetrack/model"
)

type RequestController string

const (
	BasicSpaceData    RequestController = "basicspacedata"
	ExpandedSpaceData RequestController = "expandedspacedata"
	FileShare         RequestController = "fileshare"
	CombinedOpsData   RequestController = "combineopsdata"
)

type RequestAction string

const (
	Query    RequestAction = "query"
	ModelDef RequestAction = "modeldef"
)

type RequestClass string

const (
	GP     RequestClass = "gp"
	SatCat RequestClass = "satcat"
)

type RequestPredicate string

const (
	Limit   RequestPredicate = "limit"
	OrderBy RequestPredicate = "orderby"
)

// SpaceRequest is used to hold all the information needed to build the query to fetch SpaceTrack
type SpaceRequest struct {
	Predicates      []model.Predicate
	Format          model.Format
	Limit           model.Limit
	OrderBy         model.OrderBy
	ShowEmptyResult bool
}

// BuildQuery is used to build the query of Spacetrack
func (s SpaceRequest) BuildQuery() string {
	var pathResult string
	m := []model.SpaceToPath{s.OrderBy, s.Limit, s.Format}

	// TODO we have to fix this lines
	for _, v := range s.Predicates {
		pathResult += v.ToPath()
	}

	for _, v := range m {
		pathResult += v.ToPath()
	}

	if s.ShowEmptyResult {
		pathResult += "/emptyresult/show"
	}

	return pathResult
}

// SpaceOrbitalObjArr is used to bump the response of spacetrack into a single file
type SpaceOrbitalObjArr struct {
	XMLName          xml.Name          `json:"-" xml:"spacetrack-data"`
	SpaceOrbitalObjs []SpaceOrbitalObj `json:"row" xml:"row" html:"row"`
}

// SpaceOrbitalObj is used to persist all the information in format xml or json
type SpaceOrbitalObj struct {
	XMLName            xml.Name `json:"-" xml:"row"`
	CcsdsOmmVers       string   `json:"CCSDS_OMM_VERS" xml:"CCSDS_OMM_VERS" csv:"CCSDS_OMM_VERS" html:"l=CCSDS_OMM_VERS,e=span"`
	Comment            string   `json:"COMMENT" xml:"COMMENT" csv:"COMMENT" html:"l=COMMENT,e=span"`
	CreationDate       string   `json:"CREATION_DATE" xml:"CREATION_DATE" csv:"CREATION_DATE" html:"l=CREATION_DATE,e=span"`
	Originator         string   `json:"ORIGINATOR" xml:"ORIGINATOR" csv:"ORIGINATOR" html:"l=ORIGINATOR,e=span"`
	ObjectName         string   `json:"OBJECT_NAME" xml:"OBJECT_NAME" csv:"OBJECT_NAME" html:"l=OBJECT_NAME,e=span"`
	ObjectId           string   `json:"OBJECT_ID" xml:"OBJECT_ID" csv:"OBJECT_ID" html:"l=OBJECT_ID,e=span"`
	CenterName         string   `json:"CENTER_NAME" xml:"CENTER_NAME" csv:"CENTER_NAME" html:"l=CENTER_NAME,e=span"`
	RefFrame           string   `json:"REF_FRAME" xml:"REF_FRAME" csv:"REF_FRAME" html:"l=REF_FRAME,e=span"`
	TimeSystem         string   `json:"TIME_SYSTEM" xml:"TIME_SYSTEM" csv:"TIME_SYSTEM" html:"l=TIME_SYSTEM,e=span"`
	MeanElementTheory  string   `json:"MEAN_ELEMENT_THEORY" xml:"MEAN_ELEMENT_THEORY" csv:"MEAN_ELEMENT_THEORY" html:"l=MEAN_ELEMENT_THEORY,e=span"`
	Epoch              string   `json:"EPOCH" xml:"EPOCH" csv:"EPOCH" html:"l=EPOCH,e=span"`
	MeanMotion         string   `json:"MEAN_MOTION" xml:"MEAN_MOTION" csv:"MEAN_MOTION" html:"l=MEAN_MOTION,e=span"`
	Eccentricity       string   `json:"ECCENTRICITY" xml:"ECCENTRICITY" csv:"ECCENTRICITY" html:"l=ECCENTRICITY,e=span"`
	Inclination        string   `json:"INCLINATION" xml:"INCLINATION" csv:"INCLINATION" html:"l=INCLINATION,e=span"`
	RaOfAscNode        string   `json:"RA_OF_ASC_NODE" xml:"RA_OF_ASC_NODE" csv:"RA_OF_ASC_NODE" html:"l=RA_OF_ASC_NODE,e=span"`
	ArgOfPericenter    string   `json:"ARG_OF_PERICENTER" xml:"ARG_OF_PERICENTER" csv:"ARG_OF_PERICENTER" html:"l=ARG_OF_PERICENTER,e=span"`
	MeanAnomaly        string   `json:"MEAN_ANOMALY" xml:"MEAN_ANOMALY" csv:"MEAN_ANOMALY" html:"l=MEAN_ANOMALY,e=span"`
	EphemerisType      string   `json:"EPHEMERIS_TYPE" xml:"EPHEMERIS_TYPE" csv:"EPHEMERIS_TYPE" html:"l=EPHEMERIS_TYPE,e=span"`
	ClassificationType string   `json:"CLASSIFICATION_TYPE" xml:"CLASSIFICATION_TYPE" csv:"CLASSIFICATION_TYPE" html:"l=CLASSIFICATION_TYPE,e=span"`
	NoradCatId         string   `json:"NORAD_CAT_ID" xml:"NORAD_CAT_ID" csv:"NORAD_CAT_ID" html:"l=NORAD_CAT_ID,e=span"`
	ElementSetNo       string   `json:"ELEMENT_SET_NO" xml:"ELEMENT_SET_NO" csv:"ELEMENT_SET_NO" html:"l=ELEMENT_SET_NO,e=span"`
	RevAtEpoch         string   `json:"REV_AT_EPOCH" xml:"REV_AT_EPOCH" csv:"REV_AT_EPOCH" html:"l=REV_AT_EPOCH,e=span"`
	Bstar              string   `json:"BSTAR" xml:"BSTAR" csv:"BSTAR" html:"l=BSTAR,e=span"`
	MeanMotionDot      string   `json:"MEAN_MOTION_DOT" xml:"MEAN_MOTION_DOT" csv:"MEAN_MOTION_DOT" html:"l=MEAN_MOTION_DOT,e=span"`
	MeanMotionDdot     string   `json:"MEAN_MOTION_DDOT" xml:"MEAN_MOTION_DDOT" csv:"MEAN_MOTION_DDOT" html:"l=MEAN_MOTION_DDOT,e=span"`
	SemimajorAxis      string   `json:"SEMIMAJOR_AXIS" xml:"SEMIMAJOR_AXIS" csv:"SEMIMAJOR_AXIS" html:"l=SEMIMAJOR_AXIS,e=span"`
	Period             string   `json:"PERIOD" xml:"PERIOD" csv:"PERIOD" html:"l=PERIOD,e=span"`
	Apoasis            string   `json:"APOASIS" xml:"APOASIS" csv:"APOASIS" html:"l=APOASIS,e=span"`
	Periapsis          string   `json:"PERIAPSIS" xml:"PERIAPSIS" csv:"PERIAPSIS" html:"l=PERIAPSIS,e=span"`
	ObjectType         string   `json:"OBJECT_TYPE" xml:"OBJECT_TYPE" csv:"OBJECT_TYPE" html:"l=OBJECT_TYPE,e=span"`
	RcsSize            string   `json:"RCS_SIZE" xml:"RCS_SIZE" csv:"RCS_SIZE" html:"l=RCS_SIZE,e=span"`
	CountryCode        string   `json:"COUNTRY_CODE" xml:"COUNTRY_CODE" csv:"COUNTRY_CODE" html:"l=COUNTRY_CODE,e=span"`
	LaunchDate         string   `json:"LAUNCH_DATE" xml:"LAUNCH_DATE" csv:"LAUNCH_DATE" html:"l=LAUNCH_DATE,e=span"`
	Site               string   `json:"SITE" xml:"SITE" csv:"SITE" html:"l=SITE,e=span"`
	DecayDate          string   `json:"DECAY_DATE" xml:"DECAY_DATE" csv:"DECAY_DATE" html:"l=DECAY_DATE,e=span"`
	File               string   `json:"FILE" xml:"FILE" csv:"FILE" html:"l=FILE,e=span"`
	GpId               string   `json:"GP_ID" xml:"GP_ID" scv:"GP_ID" html:"l=GP_ID,e=span"`
	TleLine0           string   `json:"TLE_LINE0" xml:"TLE_LINE0" csv:"TLE_LINE0" html:"l=TLE_LINE0,e=span"`
	TleLine1           string   `json:"TLE_LINE1" xml:"TLE_LINE1" csv:"TLE_LINE1" html:"l=TLE_LINE1,e=span"`
	TleLine2           string   `json:"TLE_LINE2" xml:"TLE_LINE2" csv:"TLE_LINE2" html:"l=TLE_LINE2,e=span"`
}
