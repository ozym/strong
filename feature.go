package strong

import (
	"strconv"
	"strings"
	"time"
)

type Feature struct {
	Properties struct {
		EventType             *string    `json:"eventtype"`
		PublicID              *string    `json:"publicid"`
		ModificationTime      *time.Time `json:"modificationtime"`
		OriginTime            *time.Time `json:"origintime"`
		OriginError           *float64   `json:"originerror"`
		EarthModel            *string    `json:"earthmodel"`
		EvaluationMethod      *string    `json:"evaluationmethod"`
		EvaluationStatus      *string    `json:"evaluationstatus"`
		EvaluationMode        *string    `json:"evaluationmode"`
		Latitude              *float64   `json:"latitude"`
		Longitude             *float64   `json:"longitude"`
		Depth                 *float64   `json:"depth"`
		DepthType             *string    `json:"depthtype"`
		UsedPhaseCount        *int32     `json:"usedphasecount"`
		UsedStationCount      *int32     `json:"usedstationcount"`
		AzimuthalGap          *float64   `json:"azimuthalgap"`
		MinimumDistance       *float64   `json:"minimumdistance"`
		Magnitude             *float64   `json:"magnitude"`
		MagnitudeType         *string    `json:"magnitudetype"`
		MagnitudeStationCount *int32     `json:"magnitudestationcount"`
		MagnitudeUncertainty  *float64   `json:"magnitudeuncertainty"`
	} `json:"properties"`
}

type Search struct {
	Features []Feature `json:"features"`
}

//Process *string `xml:"process"`
//Site    *string `xml:"site"`
//Uid     *string `xml:"uid"`

func (f *Feature) Event(agency *string) (*Event, error) {

	var mag string
	if f.Properties.Magnitude != nil {
		mag = strconv.FormatFloat(*f.Properties.Magnitude, 'f', -1, 64)
	}
	var uncertainty string
	if f.Properties.MagnitudeUncertainty != nil {
		uncertainty = strconv.FormatFloat(*f.Properties.MagnitudeUncertainty, 'f', -1, 64)
	}

	var update string
	if f.Properties.ModificationTime != nil {
		update = f.Properties.ModificationTime.Format(RFC3339Micro)
	}

	p := []string{
		*f.Properties.PublicID,
		*f.Properties.EvaluationStatus,
		f.Properties.OriginTime.Format(RFC3339Micro),
		strconv.FormatFloat(*f.Properties.Latitude, 'f', -1, 64),
		strconv.FormatFloat(*f.Properties.Longitude, 'f', -1, 64),
		strconv.FormatFloat(*f.Properties.Depth, 'f', -1, 64),
		mag,
		*f.Properties.MagnitudeType,
	}

	uid := strings.Join(p, ":")

	e := Event{
		Uid:                   &uid,
		PublicID:              f.Properties.PublicID,
		AgencyID:              agency,
		Type:                  f.Properties.EventType,
		UpdateTime:            &update,
		Status:                f.Properties.EvaluationStatus, // for want of an alternative
		Time:                  f.Properties.OriginTime,
		StandardError:         f.Properties.OriginError,
		Latitude:              f.Properties.Latitude,
		Longitude:             f.Properties.Longitude,
		Depth:                 f.Properties.Depth,
		DepthType:             f.Properties.DepthType,
		UsedPhaseCount:        f.Properties.UsedPhaseCount,
		UsedStationCount:      f.Properties.UsedStationCount,
		AzimuthalGap:          f.Properties.AzimuthalGap,
		MinimumDistance:       f.Properties.MinimumDistance,
		Magnitude:             &mag,
		MagnitudeType:         f.Properties.MagnitudeType,
		MagnitudeUncertainty:  &uncertainty,
		MagnitudeStationCount: f.Properties.MagnitudeStationCount,
		EarthModelID:          f.Properties.EarthModel,
		MethodID:              f.Properties.EvaluationMethod,
		EvaluationMode:        f.Properties.EvaluationMode,
		EvaluationStatus:      f.Properties.EvaluationStatus,
	}

	return &e, nil
}
