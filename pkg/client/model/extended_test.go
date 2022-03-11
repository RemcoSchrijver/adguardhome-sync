package model_test

import (
	"github.com/bakito/adguardhome-sync/pkg/client/model"
	. "github.com/bakito/adguardhome-sync/pkg/pointer"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Types", func() {
	Context("DNSConfig", func() {
		Context("Equals", func() {
			It("should be equal", func() {
				dc1 := &model.DNSConfig{LocalPtrUpstreams: &[]string{"a"}}
				dc2 := &model.DNSConfig{LocalPtrUpstreams: &[]string{"a"}}
				Ω(dc1.Equals(dc2)).Should(BeTrue())
			})
			It("should not be equal", func() {
				dc1 := &model.DNSConfig{LocalPtrUpstreams: &[]string{"a"}}
				dc2 := &model.DNSConfig{LocalPtrUpstreams: &[]string{"b"}}
				Ω(dc1.Equals(dc2)).ShouldNot(BeTrue())
			})
		})
	})
	Context("QueryLogConfig", func() {
		Context("Equal", func() {
			var (
				a *model.QueryLogConfig
				b *model.QueryLogConfig
			)
			BeforeEach(func() {
				a = &model.QueryLogConfig{}
				b = &model.QueryLogConfig{}
			})
			It("should be equal", func() {
				i := model.QueryLogConfigInterval(1)
				a.Enabled = ToB(true)
				a.Interval = &i
				a.AnonymizeClientIp = ToB(true)
				b.Enabled = ToB(true)
				b.Interval = &i
				b.AnonymizeClientIp = ToB(true)
				Ω(a.Equals(b)).Should(BeTrue())
			})
			It("should not be equal when enabled differs", func() {
				a.Enabled = ToB(true)
				b.Enabled = ToB(false)
				Ω(a.Equals(b)).ShouldNot(BeTrue())
			})
			It("should not be equal when interval differs", func() {
				ia := model.QueryLogConfigInterval(1)
				ib := model.QueryLogConfigInterval(2)
				a.Interval = &ia
				b.Interval = &ib
				Ω(a.Equals(b)).ShouldNot(BeTrue())
			})
			It("should not be equal when anonymizeClientIP differs", func() {
				a.AnonymizeClientIp = ToB(true)
				b.AnonymizeClientIp = ToB(false)
				Ω(a.Equals(b)).ShouldNot(BeTrue())
			})
		})
	})
})
