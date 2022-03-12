package model

// FilterSetUrlPatch URL settings
type FilterSetUrlPatch struct {
	Data      Filter  `json:"data,omitempty"`
	Url       *string `json:"url,omitempty"`
	Whitelist *bool   `json:"whitelist,omitempty"`
}

type RemoveUrlRequestPatch struct { //nolint
	RemoveUrlRequest
	Whitelist bool `json:"whitelist"`
}

type FilterStatusPatch struct {
	FilterStatus
	WhitelistFilters *[]Filter `json:"whitelist_filters"`
}
