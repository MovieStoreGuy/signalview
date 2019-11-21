package detector

// BundledPayload is the result returned from the query endpoint
type BundledPayload struct {
	// Count is the total number of detectors that match the given query, not the size of the results
	Count   int64                  `json:"count"`
	Results []DetailResultsPayload `json:"results"`
}

// DetailResultsPayload contains all the data for a V2 detector
type DetailResultsPayload struct {
	AuthorizedWriters struct {
		Teams []string `json:"teams"`
		Users []string `json:"users"`
	} `json:"authorizedWriters,omitempty"`
	// CustomProperties string        `json:"customProperties,omitempty"`
	Created         int64          `json:"created"`
	Creator         string         `json:"creator"`
	Description     string         `json:"description,omitempty"`
	ID              string         `json:"id"`
	LabelResolution map[string]int `json:"labelResolution"`
	LastUpdated     int64          `json:"lastUpdated"`
	Locked          bool           `json:"locked,omitempty"`
	Name            string         `json:"name"`
	OverMTSLimit    bool           `json:"overMTSLimit,omitempty"`
	Rules           []struct {
		Description          *string                  `json:"description,omitempty"`
		DetectLabel          string                   `json:"detectLabel"`
		Disabled             bool                     `json:"disabled,omitempty"`
		Notifications        []map[string]interface{} `json:"notifications,omitempty"`
		ParameterizedBody    string                   `json:"parameterizedBody,omitempty"`
		ParameterizedSubject string                   `json:"parameterizedSubject,omitempty"`
		RunbookURL           string                   `json:"runbookUrl,omitempty"`
		Severity             string                   `json:"severity,omitempty"`
	} `json:"rules"`
	Tags  []string `json:"tags,omitempty"`
	Teams []string `json:"teams,omitempty"`
}
