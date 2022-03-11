package types

import (
	"fmt"
	"sort"
	"strings"
)

const (
	// DefaultAPIPath default api path
	DefaultAPIPath = "/control"
)

// Config application configuration struct
type Config struct {
	Origin     AdGuardInstance   `json:"origin" yaml:"origin"`
	Replica    AdGuardInstance   `json:"replica,omitempty" yaml:"replica,omitempty"`
	Replicas   []AdGuardInstance `json:"replicas,omitempty" yaml:"replicas,omitempty"`
	Cron       string            `json:"cron,omitempty" yaml:"cron,omitempty"`
	RunOnStart bool              `json:"runOnStart,omitempty" yaml:"runOnStart,omitempty"`
	API        API               `json:"api,omitempty" yaml:"api,omitempty"`
	Features   Features          `json:"features,omitempty" yaml:"features,omitempty"`
}

// API configuration
type API struct {
	Port     int    `json:"port,omitempty" yaml:"port,omitempty"`
	Username string `json:"username,omitempty" yaml:"username,omitempty"`
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
	DarkMode bool   `json:"darkMode,omitempty" yaml:"darkMode,omitempty"`
}

// UniqueReplicas get unique replication instances
func (cfg *Config) UniqueReplicas() []AdGuardInstance {
	dedup := make(map[string]AdGuardInstance)
	if cfg.Replica.URL != "" {
		dedup[cfg.Replica.Key()] = cfg.Replica
	}
	for _, replica := range cfg.Replicas {
		if replica.URL != "" {
			dedup[replica.Key()] = replica
		}
	}

	var r []AdGuardInstance
	for _, replica := range dedup {
		if replica.APIPath == "" {
			replica.APIPath = DefaultAPIPath
		}
		r = append(r, replica)
	}
	return r
}

// AdGuardInstance AdguardHome config instance
type AdGuardInstance struct {
	URL                string `json:"url" yaml:"url"`
	APIPath            string `json:"apiPath,omitempty" yaml:"apiPath,omitempty"`
	Username           string `json:"username,omitempty" yaml:"username,omitempty"`
	Password           string `json:"password,omitempty" yaml:"password,omitempty"`
	InsecureSkipVerify bool   `json:"insecureSkipVerify" yaml:"insecureSkipVerify"`
	AutoSetup          bool   `json:"autoSetup" yaml:"autoSetup"`
}

// Key AdGuardInstance key
func (i *AdGuardInstance) Key() string {
	return fmt.Sprintf("%s#%s", i.URL, i.APIPath)
}

// Protection API struct
type Protection struct {
	ProtectionEnabled bool `json:"protection_enabled"`
}

// Status API struct
type Status struct {
	Protection
	DNSAddresses  []string `json:"dns_addresses"`
	DNSPort       int      `json:"dns_port"`
	HTTPPort      int      `json:"http_port"`
	DhcpAvailable bool     `json:"dhcp_available"`
	Running       bool     `json:"running"`
	Version       string   `json:"version"`
	Language      string   `json:"language"`
}

// RewriteEntries list of RewriteEntry
type RewriteEntries []RewriteEntry

// Merge RewriteEntries
func (rwe *RewriteEntries) Merge(other *RewriteEntries) (RewriteEntries, RewriteEntries, RewriteEntries) {
	current := make(map[string]RewriteEntry)

	var adds RewriteEntries
	var removes RewriteEntries
	var duplicates RewriteEntries
	processed := make(map[string]bool)
	for _, rr := range *rwe {
		if _, ok := processed[rr.Key()]; !ok {
			current[rr.Key()] = rr
			processed[rr.Key()] = true
		} else {
			// remove duplicate
			removes = append(removes, rr)
		}
	}

	for _, rr := range *other {
		if _, ok := current[rr.Key()]; ok {
			delete(current, rr.Key())
		} else {
			if _, ok := processed[rr.Key()]; !ok {
				adds = append(adds, rr)
				processed[rr.Key()] = true
			} else {
				//	skip duplicate
				duplicates = append(duplicates, rr)
			}
		}
	}

	for _, rr := range current {
		removes = append(removes, rr)
	}

	return adds, removes, duplicates
}

// RewriteEntry API struct
type RewriteEntry struct {
	Domain string `json:"domain"`
	Answer string `json:"answer"`
}

// Key RewriteEntry key
func (re *RewriteEntry) Key() string {
	return fmt.Sprintf("%s#%s", re.Domain, re.Answer)
}

// Filters list of Filter
type Filters []Filter

// Merge merge Filters
func (f Filters) Merge(other Filters) (Filters, Filters, Filters) {
	current := make(map[string]Filter)

	var adds Filters
	var updates Filters
	var removes Filters
	for _, f := range f {
		current[f.URL] = f
	}

	for i := range other {
		rr := other[i]
		if c, ok := current[rr.URL]; ok {
			if !c.Equals(&rr) {
				updates = append(updates, rr)
			}
			delete(current, rr.URL)
		} else {
			adds = append(adds, rr)
		}
	}

	for _, rr := range current {
		removes = append(removes, rr)
	}

	return adds, updates, removes
}

// Filter API struct
type Filter struct {
	ID         int    `json:"id"`
	Enabled    bool   `json:"enabled"`
	URL        string `json:"url"`  // needed for add
	Name       string `json:"name"` // needed for add
	RulesCount int    `json:"rules_count"`
	Whitelist  bool   `json:"whitelist"` // needed for add
}

// Equals Filter equal check
func (f *Filter) Equals(o *Filter) bool {
	return f.Enabled == o.Enabled && f.URL == o.URL && f.Name == o.Name
}

// FilterUpdate  API struct
type FilterUpdate struct {
	URL       string `json:"url"`
	Data      Filter `json:"data"`
	Whitelist bool   `json:"whitelist"`
}

// FilteringStatus API struct
type FilteringStatus struct {
	FilteringConfig
	Filters          Filters   `json:"filters"`
	WhitelistFilters Filters   `json:"whitelist_filters"`
	UserRules        UserRules `json:"user_rules"`
}

// UserRules API struct
type UserRules []string

// String toString of Users
func (ur UserRules) String() string {
	return strings.Join(ur, "\n")
}

// EnableConfig API struct
type EnableConfig struct {
	Enabled bool `json:"enabled"`
}

// IntervalConfig API struct
type IntervalConfig struct {
	Interval float64 `json:"interval"`
}

// FilteringConfig API struct
type FilteringConfig struct {
	EnableConfig
	IntervalConfig
}

// RefreshFilter API struct
type RefreshFilter struct {
	Whitelist bool `json:"whitelist"`
}

// Services API struct
type Services []string

// Sort sort Services
func (s Services) Sort() {
	sort.Strings(s)
}

// Equals Services equal check
func (s Services) Equals(o Services) bool {
	s.Sort()
	o.Sort()
	return equals(s, o)
}

// Clients API struct
type Clients struct {
	Clients     []Client `json:"clients"`
	AutoClients []struct {
		IP        string   `json:"ip"`
		Name      string   `json:"name"`
		Source    string   `json:"source"`
		WhoisInfo struct{} `json:"whois_info"`
	} `json:"auto_clients"`
	SupportedTags []string `json:"supported_tags"`
}

// Client API struct
type Client struct {
	Ids             []string `json:"ids,omitempty"`
	Tags            []string `json:"tags,omitempty"`
	BlockedServices []string `json:"blocked_services,omitempty"`
	Upstreams       []string `json:"upstreams,omitempty"`

	UseGlobalSettings        bool   `json:"use_global_settings"`
	UseGlobalBlockedServices bool   `json:"use_global_blocked_services"`
	Name                     string `json:"name"`
	FilteringEnabled         bool   `json:"filtering_enabled"`
	ParentalEnabled          bool   `json:"parental_enabled"`
	SafesearchEnabled        bool   `json:"safesearch_enabled"`
	SafebrowsingEnabled      bool   `json:"safebrowsing_enabled"`
	Disallowed               bool   `json:"disallowed"`
	DisallowedRule           string `json:"disallowed_rule"`
}

// ClientUpdate API struct
type ClientUpdate struct {
	Name string `json:"name"`
	Data Client `json:"data"`
}

func equals(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// InstallPort AdguardHome install config port
type InstallPort struct {
	IP         string `json:"ip"`
	Port       int    `json:"port"`
	Status     string `json:"status"`
	CanAutofix bool   `json:"can_autofix"`
}
