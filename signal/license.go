package signal

import (
	"encoding/json"
	"fmt"

	"github.com/dcos/dcos-signal/config"
	"github.com/segmentio/analytics-go"

	log "github.com/Sirupsen/logrus"
)

// Complete report used by signal service, composed of all requests
type LicenseReport struct {
	id                 string `json:"id"`
	number_of_nodes    int    `json:"number_of_nodes"`
	number_of_breaches int    `json:"number_of_breaches"`
	node_capacity      int    `json:"node_capacity"`
	start_timestamp    string `json:"start_timestamp"`
	end_timestamp      string `json:"end_timestamp"`
	current_timestamp  string `json:"current_timestamp"`
}

type License struct {
	Report    *LicenseReport
	Endpoints []string
	Method    string
	Headers   map[string]string
	Track     *analytics.Track
	Error     []string
	Name      string
}

func (d *License) getName() string {
	return d.Name
}

func (d *License) setReport(body []byte) error {
	if err := json.Unmarshal(body, &d.Report); err != nil {
		return err
	}
	return nil
}

func (d *License) getReport() interface{} {
	return d.Report
}

func (d *License) addHeaders(head map[string]string) {
	for k, v := range head {
		d.Headers[k] = v
	}
}
func (d *License) getHeaders() map[string]string {
	return d.Headers
}

func (d *License) getEndpoints() []string {
	if len(d.Endpoints) != 1 {
		log.Errorf("License needs 1 endpoints, got %d", len(d.Endpoints))
	}
	return d.Endpoints
}

func (d *License) getMethod() string {
	return d.Method
}

func (d *License) getError() []string {
	return d.Error
}

func (d *License) appendError(err string) {
	d.Error = append(d.Error, err)
}

func (d *License) setTrack(c config.Config) error {
	if d.Report == nil {
		return fmt.Errorf("%s report is nil, bailing out.", d.Name)
	}

	properties := map[string]interface{}{
		"source":             "cluster",
		"customerKey":        c.CustomerKey,
		"environmentVersion": c.DCOSVersion,
		"clusterId":          c.ClusterID,
		"variant":            c.DCOSVariant,
		"platform":           c.GenPlatform,
		"provider":           c.GenProvider,
		"id":                 d.Report.id,
		"number_of_nodes":    d.Report.number_of_nodes,
		"number_of_breaches": d.Report.number_of_breaches,
		"node_capacity":      d.Report.node_capacity,
		"start_timestamp":    d.Report.start_timestamp,
		"end_timestamp":      d.Report.end_timestamp,
		"current_timestamp":  d.Report.current_timestamp,
	}

	d.Track = &analytics.Track{
		Event:       "license_track",
		UserId:      c.CustomerKey,
		AnonymousId: c.ClusterID,
		Properties:  properties,
	}
	return nil
}

func (d *License) getTrack() *analytics.Track {
	return d.Track
}

func (d *License) sendTrack(c config.Config) error {
	ac := CreateSegmentClient(c.SegmentKey, c.FlagVerbose)
	defer ac.Close()
	err := ac.Track(d.Track)
	return err
}
