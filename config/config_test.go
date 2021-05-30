package config_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"vladusenko.io/home-torrent/config"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

var _ = Describe("Config", func() {
	BeforeEach(func() {
		config.Reset()
	})

	It("Should parse and validate a config file", func() {
		conf, err := config.GetConfig("./testfixtures/simple_config.json")

		Expect(err).To(BeNil())
		Expect(conf).To(Equal(&config.Config{
			HttpPort: 8080,
		}))
	})

	It("Should return error for an invalid config file", func() {
		conf, err := config.GetConfig("./testfixtures/invalid_config.json")

		Expect(conf).To(BeNil())
		Expect(err).NotTo(BeNil())
	})
})
