package stencil_test

import (
	"bufio"
	"bytes"
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

	It("can eval Exec in container", func() {
		stencil, err := factory.Get("curl:latest")
		Expect(err).ToNot(HaveOccurred())
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		exec := Exec{
			Args:    []string{"/usr/bin/curl", "--version"},
			Stencil: stencil,
			Stdout:  bufio.NewWriter(&stdout),
			Stderr:  bufio.NewWriter(&stderr),
		}
		state, err := exec.Eval()
		Expect(err).ToNot(HaveOccurred())
		Expect(string(stderr.Bytes())).To(Equal(""))
		Expect(state.ExitCode).To(Equal(0))
		Expect(string(stdout.Bytes())).To(HavePrefix("curl"))
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
		defer os.RemoveAll(testMount)
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
			Mounts: []Mount{{
				Source:      testMount,
				Destination: testMount,
			}},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(string(out)).To(HavePrefix("Hello World from mount"))
	})

	It("multiple processes can have conflicting mounts", func() {
		testMount1, err := ioutil.TempDir("", "testmount1")
		Expect(err).ToNot(HaveOccurred())
		defer os.RemoveAll(testMount1)
		testMount2, err := ioutil.TempDir("", "testmount2")
		Expect(err).ToNot(HaveOccurred())
		defer os.RemoveAll(testMount2)
		stencil, err := factory.Get("curl:latest")
		Expect(err).ToNot(HaveOccurred())
		go Run(Exec{
			Args:    []string{"/bin/sh", "-c", "echo 'test1' > /data/test"},
			Stencil: stencil,
			Mounts: []Mount{{
				Source:      testMount1,
				Destination: "/data",
				Writable:    true,
			}},
		})
		_, err = Run(Exec{
			Args:    []string{"/bin/sh", "-c", "echo 'test2' > /data/test"},
			Stencil: stencil,
			Mounts: []Mount{{
				Source:      testMount2,
				Destination: "/data",
				Writable:    true,
			}},
		})
		Expect(err).ToNot(HaveOccurred())
		out, err := Run(Exec{
			Args:    []string{"/bin/sh", "-c", "cat /data1/test /data2/test"},
			Stencil: stencil,
			Mounts: []Mount{{
				Source:      testMount1,
				Destination: "/data1",
				Writable:    true,
			}, {
				Source:      testMount2,
				Destination: "/data2",
				Writable:    true,
			}},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(string(out)).To(Equal("test1\ntest2\n"))
	})

	It("returns exit code of a exited process", func() {
		stencil, err := factory.Get("curl:latest")
		Expect(err).ToNot(HaveOccurred())
		var stderr bytes.Buffer
		exec := Exec{
			Args:    []string{"/usr/bin/curl", "-I", "https://localhost:123"},
			Stencil: stencil,
			Stderr:  bufio.NewWriter(&stderr),
		}
		state, err := exec.Eval()
		Expect(err).ToNot(HaveOccurred())
		Expect(state.ExitCode).ToNot(Equal(0))
		Expect(string(stderr.Bytes())).To(ContainSubstring("Connection refused"))
	})

	It("returns exit code of a exited process", func() {
		stencil, err := factory.Get("curl:latest")
		Expect(err).ToNot(HaveOccurred())
		_, err = Run(Exec{
			Args:    []string{"/usr/bin/curl", "-I", "https://localhost:123"},
			Stencil: stencil,
		})
		Expect(err.Error()).To(ContainSubstring("Connection refused"))
	})

	AfterSuite(func() {
		err := os.RemoveAll(rootDir)
		Expect(err).ToNot(HaveOccurred())
		err = dockerRegistry.Process.Kill()
		Expect(err).ToNot(HaveOccurred())
	})

})
