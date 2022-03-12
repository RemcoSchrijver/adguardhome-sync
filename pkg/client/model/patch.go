package model

type FilterSetUrlData struct { //nolint
	Enabled *bool   `json:"enabled,omitempty"`
	Name    *string `json:"name,omitempty"`
	Url     *string `json:"url,omitempty"` //nolint
}

type RemoveUrlRequestPatch struct { //nolint
	RemoveUrlRequest
	Whitelist bool `json:"whitelist"`
}

type FilterStatusPatch struct {
	FilterStatus
	WhitelistFilters *[]Filter `json:"whitelist_filters"`
}
