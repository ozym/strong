package strong

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

const RFC3339Micro = "2006-01-02T15:04:05.999999Z"

type Event struct {
	XMLName string `xml:"event"`

	PublicID *string `xml:"publicID,attr"`

	AgencyID   *string `xml:"creationInfo>agencyID"`
	UpdateTime *string `xml:"creationInfo>updateTime"`

	Process *string `xml:"process"`
	Site    *string `xml:"site"`
	Uid     *string `xml:"uid"`
	Type    *string `xml:"type"`
	Status  *string `xml:"status"`

	MethodID         *string `xml:"methodID"`
	EarthModelID     *string `xml:"earthModelID"`
	EvaluationMode   *string `xml:"evaluationMode"`
	EvaluationStatus *string `xml:"evaluationStatus"`

	Time                  *time.Time `xml:"preferredOrigin>time>value"`
	Latitude              *float64   `xml:"preferredOrigin>latitude>value"`
	Longitude             *float64   `xml:"preferredOrigin>longitude>value"`
	Depth                 *float64   `xml:"preferredOrigin>depth>value"`
	DepthType             *string    `xml:"preferredOrigin>depthType"`
	UsedPhaseCount        *int32     `xml:"preferredOrigin>quality>usedPhaseCount"`
	UsedStationCount      *int32     `xml:"preferredOrigin>quality>usedStationCount"`
	StandardError         *float64   `xml:"preferredOrigin>quality>standardError"`
	AzimuthalGap          *float64   `xml:"preferredOrigin>quality>azimuthalGap"`
	MinimumDistance       *float64   `xml:"preferredOrigin>quality>minimumDistance"`
	Magnitude             *string    `xml:"preferredOrigin>preferredMagnitude>magnitude>value,omitempty"`
	MagnitudeUncertainty  *string    `xml:"preferredOrigin>preferredMagnitude>magnitude>uncertainty,omitempty"`
	MagnitudeType         *string    `xml:"preferredOrigin>preferredMagnitude>type,omitempty"`
	MagnitudeStationCount *int32     `xml:"preferredOrigin>preferredMagnitude>stationCount,omitempty"`
}

func (e *Event) GetPublicID() (string, error) {
	if e.PublicID == nil {
		return "", errors.New("missing event publicID")
	}
	return *e.PublicID, nil
}

func (e *Event) GetType() (string, error) {
	if e.Type == nil {
		return "", errors.New("missing event type")
	}
	return *e.Type, nil
}

func (e *Event) GetUpdateTime() (time.Time, error) {

	if e.UpdateTime == nil {
		return time.Time{}, errors.New("missing event update time")
	}

	return time.Parse(RFC3339Micro, *e.UpdateTime)
}

func (e *Event) GetMagnitude() (float64, error) {
	if e.Magnitude == nil {
		return 0.0, errors.New("missing event preferred origin magnitude")
	}
	return strconv.ParseFloat(*e.Magnitude, 64)
}

func (e *Event) After(event *Event) (bool, error) {
	t1, err := e.GetUpdateTime()
	if err != nil {
		return false, err
	}
	t2, err := event.GetUpdateTime()
	if err != nil {
		return false, err
	}

	return t1.After(t2), nil
}

func (e *Event) Before(event *Event) (bool, error) {
	t1, err := e.GetUpdateTime()
	if err != nil {
		return false, err
	}
	fmt.Println(t1)
	t2, err := event.GetUpdateTime()
	if err != nil {
		return false, err
	}
	fmt.Println(t2)

	return t1.Before(t2), nil
}

func (e *Event) Marshal() ([]byte, error) {

	res := ([]byte)(xml.Header)

	b, err := xml.MarshalIndent(e, "", "   ")
	if err != nil {
		return nil, err
	}

	return append(res, b...), nil
}

func (e *Event) Write(outfile string) (int, error) {

	b, err := e.Marshal()
	if err != nil {
		return 0, err
	}

	file, err := os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return file.Write(b)
}
