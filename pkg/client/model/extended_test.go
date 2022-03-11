package model_test

import (
	"github.com/bakito/adguardhome-sync/pkg/client/model"
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
})
