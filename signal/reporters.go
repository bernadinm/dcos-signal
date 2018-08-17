package signal

import (
	"github.com/dcos/dcos-signal/config"
)

func makeReporters(c config.Config) (chan Reporter, error) {

	var reporters = []Reporter{
		&Diagnostics{
			Name:      "diagnostics",
			Endpoints: c.DiagnosticsURLs,
			Method:    "GET",
			Headers: map[string]string{
				"content-type": "application/json",
			},
		},
		&Cosmos{
			Name:      "cosmos",
			Endpoints: c.CosmosURLs,
			Method:    "POST",
			Headers: map[string]string{
				"content-type": "application/vnd.dcos.package.list-request+json;charset=utf-8;version=v1",
				"accept":       "application/vnd.dcos.package.list-response+json;charset=utf-8;version=v1",
			},
		},
		&Mesos{
			Name:      "mesos",
			Endpoints: c.MesosURLs,
			Method:    "GET",
			Headers: map[string]string{
				"content-type": "application/json",
			},
		},
		&License{
			Name:      "license",
			Endpoints: c.LicenseURLs,
			Method:    "GET",
			Headers: map[string]string{
				"content-type": "application/json",
			},
		},
	}

	reportChan := make(chan Reporter, len(reporters))
	for _, r := range reporters {
		r.addHeaders(c.ExtraHeaders)
		reportChan <- r
	}

	return reportChan, nil
}
