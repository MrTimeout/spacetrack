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

// SpaceOrbitalObj is used to persist all the information in format xml or json
type SpaceOrbitalObj struct {
	XMLName            xml.Name `xml:"row"`
	CcsdsOmmVers       string   `json:"CCSDS_OMM_VERS" xml:"CCSDS_OMM_VERS"`
	Comment            string   `json:"COMMENT" xml:"COMMENT"`
	CreationDate       string   `json:"CREATION_DATE" xml:"CREATION_DATE"`
	Originator         string   `json:"ORIGINATOR" xml:"ORIGINATOR"`
	ObjectName         string   `json:"OBJECT_NAME" xml:"OBJECT_NAME"`
	ObjectId           string   `json:"OBJECT_ID" xml:"OBJECT_ID"`
	CenterName         string   `json:"CENTER_NAME" xml:"CENTER_NAME"`
	RefFrame           string   `json:"REF_FRAME" xml:"REF_FRAME"`
	TimeSystem         string   `json:"TIME_SYSTEM" xml:"TIME_SYSTEM"`
	MeanElementTheory  string   `json:"MEAN_ELEMENT_THEORY" xml:"MEAN_ELEMENT_THEORY"`
	Epoch              string   `json:"EPOCH" xml:"EPOCH"`
	MeanMotion         string   `json:"MEAN_MOTION" xml:"MEAN_MOTION"`
	Eccentricity       string   `json:"ECCENTRICITY" xml:"ECCENTRICITY"`
	Inclination        string   `json:"INCLINATION" xml:"INCLINATION"`
	RaOfAscNode        string   `json:"RA_OF_ASC_NODE" xml:"RA_OF_ASC_NODE"`
	ArgOfPericenter    string   `json:"ARG_OF_PERICENTER" xml:"ARG_OF_PERICENTER"`
	MeanAnomaly        string   `json:"MEAN_ANOMALY" xml:"MEAN_ANOMALY"`
	EphemerisType      string   `json:"EPHEMERIS_TYPE" xml:"EPHEMERIS_TYPE"`
	ClassificationType string   `json:"CLASSIFICATION_TYPE" xml:"CLASSIFICATION_TYPE"`
	NoradCatId         string   `json:"NORAD_CAT_ID" xml:"NORAD_CAT_ID"`
	ElementSetNo       string   `json:"ELEMENT_SET_NO" xml:"ELEMENT_SET_NO"`
	RevAtEpoch         string   `json:"REV_AT_EPOCH" xml:"REV_AT_EPOCH"`
	Bstar              string   `json:"BSTAR" xml:"BSTAR"`
	MeanMotionDot      string   `json:"MEAN_MOTION_DOT" xml:"MEAN_MOTION_DOT"`
	MeanMotionDdot     string   `json:"MEAN_MOTION_DDOT" xml:"MEAN_MOTION_DDOT"`
	SemimajorAxis      string   `json:"SEMIMAJOR_AXIS" xml:"SEMIMAJOR_AXIS"`
	Period             string   `json:"PERIOD" xml:"PERIOD"`
	Apoasis            string   `json:"APOASIS" xml:"APOASIS"`
	Periapsis          string   `json:"PERIAPSIS" xml:"PERIAPSIS"`
	ObjectType         string   `json:"OBJECT_TYPE" xml:"OBJECT_TYPE"`
	RcsSize            string   `json:"RCS_SIZE" xml:"RCS_SIZE"`
	CountryCode        string   `json:"COUNTRY_CODE" xml:"COUNTRY_CODE"`
	LaunchDate         string   `json:"LAUNCH_DATE" xml:"LAUNCH_DATE"`
	Site               string   `json:"SITE" xml:"SITE"`
	DecayDate          string   `json:"DECAY_DATE" xml:"DECAY_DATE"`
	File               string   `json:"FILE" xml:"FILE"`
	GpId               string   `json:"GP_ID" xml:"GP_ID"`
	TleLine0           string   `json:"TLE_LINE0" xml:"TLE_LINE0"`
	TleLine1           string   `json:"TLE_LINE1" xml:"TLE_LINE1"`
	TleLine2           string   `json:"TLE_LINE2" xml:"TLE_LINE2"`
}
