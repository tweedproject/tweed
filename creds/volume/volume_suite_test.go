package volume_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestVolume(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Volume Suite")
}
