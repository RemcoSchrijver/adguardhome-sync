package model

import "encoding/json"

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
