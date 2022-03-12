package pointer

import "github.com/bakito/adguardhome-sync/pkg/client/model"

func ToB(b bool) *bool {
	return &b
}

func FromB(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func FromI(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func ToI(i int) *int {
	return &i
}

func FromQueryLogConfigInterval(i *model.QueryLogConfigInterval) model.QueryLogConfigInterval {
	if i == nil {
		return 0
	}
	return *i
}
