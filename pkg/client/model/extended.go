package model

import (
	"encoding/json"
	"sort"
)

// Equals dhcp config equal check
func (c *DhcpStatus) Equals(o *DhcpStatus) bool {
	a, _ := json.Marshal(c)
	b, _ := json.Marshal(o)
	return string(a) == string(b)
}

func (c *DhcpStatus) ToConfig() *DhcpConfig {
	return &DhcpConfig{
		Enabled:       c.Enabled,
		InterfaceName: c.InterfaceName,
		V4:            c.V4,
		V6:            c.V6,
	}
}

// DhcpStaticLeaseMerge merge the leases
func DhcpStaticLeaseMerge(src *[]DhcpStaticLease, dest *[]DhcpStaticLease) ([]DhcpStaticLease, []DhcpStaticLease) {
	current := make(map[string]DhcpStaticLease)

	var adds []DhcpStaticLease
	var removes []DhcpStaticLease
	if src != nil {
		for _, le := range *src {
			current[le.Mac] = le
		}
	}
	if dest != nil {
		for _, le := range *dest {
			if _, ok := current[le.Mac]; ok {
				delete(current, le.Mac)
			} else {
				adds = append(adds, le)
			}
		}
	}

	for _, rr := range current {
		removes = append(removes, rr)
	}

	return adds, removes
}

// Equals dns config equal check
func (c *DNSConfig) Equals(o *DNSConfig) bool {
	c.Sort()
	o.Sort()

	a, _ := json.Marshal(c)
	b, _ := json.Marshal(o)
	return string(a) == string(b)
}

// Sort sort dns config
func (c *DNSConfig) Sort() {
	safeSort(c.LocalPtrUpstreams)
	safeSort(c.BootstrapDns)
	safeSort(c.LocalPtrUpstreams)
}

func safeSort(s *[]string) {
	if s != nil {
		sort.Strings(*s)
	}
}

// Equals access list equal check
func (al *AccessList) Equals(o *AccessList) bool {
	return safeEquals(al.AllowedClients, o.AllowedClients) &&
		safeEquals(al.DisallowedClients, o.DisallowedClients) &&
		safeEquals(al.BlockedHosts, o.BlockedHosts)
}

func safeEquals(a *[]string, b *[]string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return equals(*a, *b)
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

// Equals QueryLogConfig equal check
func (qlc *QueryLogConfig) Equals(o *QueryLogConfig) bool {
	return qlc.Enabled == o.Enabled && qlc.AnonymizeClientIp == o.AnonymizeClientIp && qlc.Interval == o.Interval
}

// Sort sort clients
func (cl *Client) Sort() {
	safeSort(cl.Ids)
	safeSort(cl.Tags)
	safeSort(cl.BlockedServices)
	safeSort(cl.Upstreams)
}

// Equals Clients equal check
func (cl *Client) Equals(o *Client) bool {
	cl.Sort()
	o.Sort()

	a, _ := json.Marshal(cl)
	b, _ := json.Marshal(o)
	return string(a) == string(b)
}

// Merge merge Clients
func (clients *Clients) Merge(other *Clients) ([]Client, []Client, []string) {
	current := make(map[string]Client)
	if clients.Clients != nil {
		for _, client := range *clients.Clients {
			current[*client.Name] = client
		}
	}

	expected := make(map[string]Client)
	if other.Clients != nil {
		for _, client := range *other.Clients {
			expected[*client.Name] = client
		}
	}

	var adds []Client
	var removes []string
	var updates []Client

	for _, cl := range expected {
		if oc, ok := current[*cl.Name]; ok {
			if !cl.Equals(&oc) {
				updates = append(updates, cl)
			}
			delete(current, *cl.Name)
		} else {
			adds = append(adds, cl)
		}
	}

	for _, rr := range current {
		removes = append(removes, *rr.Name)
	}

	return adds, updates, removes
}
