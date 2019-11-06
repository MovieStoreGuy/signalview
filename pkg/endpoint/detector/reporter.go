package detector

import (
	"fmt"
	"bytes"
	"text/template"
	"time"
)

const report = `======================
Date of report: {{.time}}
Number of detectors: {{.count}}
======================
{{ if .details.OverMts }}
======================
| Over MTS detectors |
======================
{{ range $detector := .details.OverMts }} {{ range $field, $value := $detector}}| {{ $field }} : {{ $value }} {{ end }}
{{ end }}
======================
{{ end }}
{{ if .details.DisabledRules }}
======================
| Disabled Detectors |
======================
{{ range $detector := .details.DisabledDetectors }} {{ range $field, $value := $detector}}| {{ $field }} : {{ $value }} {{ end }}
{{ end }}
======================
{{ end }}
{{ if .details.OverMts }}
======================
| Disabled Rules     |
======================
{{ range $detector := .details.DisabledRules }} {{ range $field, $value := $detector}}| {{ $field }} : {{$value}} {{ end }}
{{ end }}
======================
{{ end }}
{{ if .details.NoNotifications }}
======================
| No notifications   |
======================
{{ range $detector := .details.DisabledRules }} {{ range $field, $value := $detector}}| {{ $field }} : {{$value}} {{ end }}
{{ end }}
======================
{{ end }}
`

type Reporter struct {
	Payload *BundledPayload
	Details struct {
		OverMts           []map[string]string
		DisabledRules     []map[string]string
		DisabledDetectors []map[string]string
		NoNotifications   []map[string]string
	}
}

func (rp *Reporter) Generate() (string, error) {
	for _, payload := range rp.Payload.Results {
		if payload.OverMTSLimit {
			rp.Details.OverMts = append(rp.Details.OverMts, map[string]string{
				"id":   payload.ID,
				"name": payload.Name,
			})
		}
		var disabled uint
		for _, rules := range payload.Rules {
			if rules.Disabled {
				disabled++
				rp.Details.DisabledRules = append(rp.Details.DisabledRules, map[string]string{
					"id":       payload.ID,
					"name":     payload.Name,
					"label":    rules.DetectLabel,
					"created":  fmt.Sprint(payload.Created),
					"modified": fmt.Sprint(payload.LastUpdated),
				})
			}
			if len(rules.Notifications) == 0 && !rules.Disabled {
				rp.Details.NoNotifications = append(rp.Details.NoNotifications, map[string]string{
					"id":       payload.ID,
					"name":     payload.Name,
					"label":    rules.DetectLabel,
					"created":  fmt.Sprint(payload.Created),
					"modified": fmt.Sprint(payload.LastUpdated),
				})
			}
		}
		if disabled == uint(len(payload.Rules)) {
			rp.Details.DisabledDetectors = append(rp.Details.DisabledDetectors, map[string]string{
				"id":       payload.ID,
				"name":     payload.Name,
				"created":  fmt.Sprint(payload.Created),
				"modified": fmt.Sprint(payload.LastUpdated),
			})
		}
	}
	tmpl, err := template.New("detector report").Parse(report)
	if err != nil {
		return "", err
	}
	buff := bytes.NewBuffer(nil)
	if err = tmpl.Execute(buff, map[string]interface{}{
		"time":    time.Now(),
		"count":   rp.Payload.Count,
		"details": rp.Details,
	}); err != nil {
		return "", err
	}
	return buff.String(), nil
}
