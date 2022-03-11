package model_test

import (
	"github.com/bakito/adguardhome-sync/pkg/client/model"
	. "github.com/bakito/adguardhome-sync/pkg/pointer"
	"github.com/google/uuid"
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

	Context("Clients", func() {
		Context("Merge", func() {
			var (
				originClients  *model.Clients
				replicaClients model.Clients
				name           string
			)
			BeforeEach(func() {
				originClients = &model.Clients{}
				replicaClients = model.Clients{}
				name = uuid.NewString()
			})

			It("should add a missing client", func() {
				*originClients.Clients = append(*originClients.Clients, model.Client{Name: &name})
				a, u, d := replicaClients.Merge(originClients)
				Ω(a).Should(HaveLen(1))
				Ω(u).Should(BeEmpty())
				Ω(d).Should(BeEmpty())

				Ω(a[0].Name).Should(Equal(name))
			})

			It("should remove additional client", func() {
				*replicaClients.Clients = append(*replicaClients.Clients, model.Client{Name: &name})
				a, u, d := replicaClients.Merge(originClients)
				Ω(a).Should(BeEmpty())
				Ω(u).Should(BeEmpty())
				Ω(d).Should(HaveLen(1))

				Ω(d[0]).Should(Equal(name))
			})

			It("should update existing client when name differs", func() {
				disallowed := true
				*originClients.Clients = append(*originClients.Clients, model.Client{Name: &name, FilteringEnabled: ToB(disallowed)})
				*replicaClients.Clients = append(*replicaClients.Clients, model.Client{Name: &name, FilteringEnabled: ToB(!disallowed)})
				a, u, d := replicaClients.Merge(originClients)
				Ω(a).Should(BeEmpty())
				Ω(u).Should(HaveLen(1))
				Ω(d).Should(BeEmpty())

				Ω(*u[0].FilteringEnabled).Should(Equal(disallowed))
			})
		})
	})
	Context("BlockedServices", func() {
		Context("Equals", func() {
			It("should be equal", func() {
				s1 := model.BlockedServicesArray([]string{"a", "b"})
				s2 := model.BlockedServicesArray([]string{"b", "a"})
				Ω(s1.Equals(s2)).Should(BeTrue())
			})
			It("should not be equal different values", func() {
				s1 := model.BlockedServicesArray([]string{"a", "b"})
				s2 := model.BlockedServicesArray([]string{"B", "a"})
				Ω(s1.Equals(s2)).ShouldNot(BeTrue())
			})
			It("should not be equal different length", func() {
				s1 := model.BlockedServicesArray([]string{"a", "b"})
				s2 := model.BlockedServicesArray([]string{"b", "a", "c"})
				Ω(s1.Equals(s2)).ShouldNot(BeTrue())
			})
		})
	})
})
