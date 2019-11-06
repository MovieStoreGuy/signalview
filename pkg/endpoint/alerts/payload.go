package alerts

type BundledMutingRuleResults struct {
	Count   int64              `json:"count"`
	Results []MutingRuleResult `json:"results"`
}

type MutingRuleResult struct {
	Created       int64               `json:"created"`
	Creator       string              `json:"creator,omitempty"`
	Description   string              `json:"description,omitempty"`
	Filters       []map[string]string `json:"filters,omitempty"`
	ID            string              `json:"id"`
	LastUpdated   uint64              `json:"lastUpdated"`
	LastUpdatedBy string              `json:"lastUpdatedBy,omitempty"`
	StartTime     uint64              `json:"startTime"`
	StopTime      uint64              `json:"stopTime"`
}
