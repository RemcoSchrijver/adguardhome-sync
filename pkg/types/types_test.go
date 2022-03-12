package types_test

import (
	"github.com/bakito/adguardhome-sync/pkg/types"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Types", func() {
	var (
		url     string
		apiPath string
	)
	BeforeEach(func() {
		url = "https://" + uuid.NewString()
		apiPath = "/" + uuid.NewString()
	})

	Context("AdGuardInstance", func() {
		It("should build a key with url and api apiPath", func() {
			i := &types.AdGuardInstance{URL: url, APIPath: apiPath}
			Ω(i.Key()).Should(Equal(url + "#" + apiPath))
		})
	})
	Context("RewriteEntry", func() {
		It("should build a key with url and api apiPath", func() {
			domain := uuid.NewString()
			answer := uuid.NewString()
			re := &types.RewriteEntry{Domain: domain, Answer: answer}
			Ω(re.Key()).Should(Equal(domain + "#" + answer))
		})
	})

	Context("RewriteEntries", func() {
		Context("Merge", func() {
			var (
				originRE  types.RewriteEntries
				replicaRE types.RewriteEntries
				domain    string
			)
			BeforeEach(func() {
				originRE = types.RewriteEntries{}
				replicaRE = types.RewriteEntries{}
				domain = uuid.NewString()
			})

			It("should add a missing rewrite entry", func() {
				originRE = append(originRE, types.RewriteEntry{Domain: domain})
				a, r, d := replicaRE.Merge(&originRE)
				Ω(a).Should(HaveLen(1))
				Ω(r).Should(BeEmpty())
				Ω(d).Should(BeEmpty())

				Ω(a[0].Domain).Should(Equal(domain))
			})

			It("should remove additional ewrite entry", func() {
				replicaRE = append(replicaRE, types.RewriteEntry{Domain: domain})
				a, r, d := replicaRE.Merge(&originRE)
				Ω(a).Should(BeEmpty())
				Ω(r).Should(HaveLen(1))
				Ω(d).Should(BeEmpty())

				Ω(r[0].Domain).Should(Equal(domain))
			})

			It("should have no changes", func() {
				originRE = append(originRE, types.RewriteEntry{Domain: domain})
				replicaRE = append(replicaRE, types.RewriteEntry{Domain: domain})
				a, r, d := replicaRE.Merge(&originRE)
				Ω(a).Should(BeEmpty())
				Ω(r).Should(BeEmpty())
				Ω(d).Should(BeEmpty())
			})

			It("should remove target duplicate", func() {
				originRE = append(originRE, types.RewriteEntry{Domain: domain})
				replicaRE = append(replicaRE, types.RewriteEntry{Domain: domain})
				replicaRE = append(replicaRE, types.RewriteEntry{Domain: domain})
				a, r, d := replicaRE.Merge(&originRE)
				Ω(a).Should(BeEmpty())
				Ω(r).Should(HaveLen(1))
				Ω(d).Should(BeEmpty())
			})

			It("should remove target duplicate", func() {
				originRE = append(originRE, types.RewriteEntry{Domain: domain})
				originRE = append(originRE, types.RewriteEntry{Domain: domain})
				replicaRE = append(replicaRE, types.RewriteEntry{Domain: domain})
				a, r, d := replicaRE.Merge(&originRE)
				Ω(a).Should(BeEmpty())
				Ω(r).Should(BeEmpty())
				Ω(d).Should(HaveLen(1))
			})
		})
	})

	Context("Config", func() {
		var cfg *types.Config
		BeforeEach(func() {
			cfg = &types.Config{}
		})
		Context("UniqueReplicas", func() {
			It("should be empty if noting defined", func() {
				r := cfg.UniqueReplicas()
				Ω(r).Should(BeEmpty())
			})
			It("should be empty if replica url is not set", func() {
				cfg.Replica = types.AdGuardInstance{URL: ""}
				r := cfg.UniqueReplicas()
				Ω(r).Should(BeEmpty())
			})
			It("should be empty if replicas url is not set", func() {
				cfg.Replicas = []types.AdGuardInstance{{URL: ""}}
				r := cfg.UniqueReplicas()
				Ω(r).Should(BeEmpty())
			})
			It("should return only one replica if same url and apiPath", func() {
				cfg.Replica = types.AdGuardInstance{URL: url, APIPath: apiPath}
				cfg.Replicas = []types.AdGuardInstance{{URL: url, APIPath: apiPath}, {URL: url, APIPath: apiPath}}
				r := cfg.UniqueReplicas()
				Ω(r).Should(HaveLen(1))
			})
			It("should return 3 one replicas if urls are different", func() {
				cfg.Replica = types.AdGuardInstance{URL: url, APIPath: apiPath}
				cfg.Replicas = []types.AdGuardInstance{{URL: url + "1", APIPath: apiPath}, {URL: url, APIPath: apiPath + "1"}}
				r := cfg.UniqueReplicas()
				Ω(r).Should(HaveLen(3))
			})
			It("should set default api apiPath if not set", func() {
				cfg.Replica = types.AdGuardInstance{URL: url}
				cfg.Replicas = []types.AdGuardInstance{{URL: url + "1"}}
				r := cfg.UniqueReplicas()
				Ω(r).Should(HaveLen(2))
				Ω(r[0].APIPath).Should(Equal(types.DefaultAPIPath))
				Ω(r[1].APIPath).Should(Equal(types.DefaultAPIPath))
			})
		})
	})
})
