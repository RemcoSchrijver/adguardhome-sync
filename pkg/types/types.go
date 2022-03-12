package types

import (
	"fmt"
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
