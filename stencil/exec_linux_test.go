package stencil_test

import (
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
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
		factory = NewFactory(rootDir, logger)
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
			Args:    []string{"/usr/bin/curl", "--version"},
			Stencil: stencil,
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(string(out)).To(HavePrefix("curl"))
	})

	It("can reach google.com from container", func() {
		stencil, err := factory.Get("curl:latest")
		Expect(err).ToNot(HaveOccurred())
		out, err := Run(Exec{
			Args:    []string{"/usr/bin/curl", "-I", "google.com"},
			Stencil: stencil,
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(string(out)).To(HavePrefix("HTTP/1.1 301 Moved Permanently"))
	})

	It("can mount a directory into container", func() {
		testMount, err := ioutil.TempDir("", "testmount")
		Expect(err).ToNot(HaveOccurred())
		testFile := path.Join(testMount, "foo")
		f, err := os.Create(testFile)
		Expect(err).ToNot(HaveOccurred())
		_, err = f.WriteString("Hello World from mount")
		Expect(err).ToNot(HaveOccurred())
		err = f.Close()
		Expect(err).ToNot(HaveOccurred())
		stencil, err := factory.Get("curl:latest")
		Expect(err).ToNot(HaveOccurred())
		out, err := Run(Exec{
			Args:    []string{"/bin/cat", testFile},
			Stencil: stencil,
			Mounts:  []string{testMount},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(string(out)).To(HavePrefix("Hello World from mount"))
	})

	AfterSuite(func() {
		err := os.RemoveAll(rootDir)
		Expect(err).ToNot(HaveOccurred())
		err = dockerRegistry.Process.Kill()
		Expect(err).ToNot(HaveOccurred())
	})

})
