package stencil_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestStencil(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Stencil Suite")
}
