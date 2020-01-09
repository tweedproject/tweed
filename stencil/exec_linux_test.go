package stencil_test

import (
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/tweedproject/tweed/stencil"
)

var _ = Describe("Exec", func() {
	var (
		factory *Factory
		//		stencil        *Stencil
		rootDir        string
		dockerRegistry *exec.Cmd
	)

	BeforeSuite(func() {
		rootDir, err := ioutil.TempDir("", "tweedroot")
		Expect(err).ToNot(HaveOccurred())
		logger := log.New(GinkgoWriter, "", 0)
		factory = NewFactory(rootDir, "/tweed/tweed", logger)
		dockerRegistry = exec.Command("registry", "serve", "/etc/docker/registry/config.yml")
		err = dockerRegistry.Start()
		for {
			_, err := net.Dial("tcp", "localhost:5000")
			if err == nil {
				return
			}
			time.Sleep(1 * time.Second)
		}
		Expect(err).ToNot(HaveOccurred())
	})

	It("loads stencil image", func() {
		err := factory.Load("curl:latest")
		Expect(err).ToNot(HaveOccurred())

	})

	It("can run Exec in container", func() {
		stencil, err := factory.Get("curl:latest")
		Expect(err).ToNot(HaveOccurred())
		out, err := Run(Exec{
			Run:     "curl",
			Stencil: stencil,
		})
		Expect(err).ToNot(HaveOccurred())

		Expect(out).To(Equal("foo"))
	})

	AfterSuite(func() {
		err := os.RemoveAll(rootDir)
		Expect(err).ToNot(HaveOccurred())
		err = dockerRegistry.Process.Kill()
		Expect(err).ToNot(HaveOccurred())
	})

})
