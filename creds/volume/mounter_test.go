package volume_test

import (
	"context"
	"io/ioutil"
	"os"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/tweedproject/tweed/creds/credsfakes"
	. "github.com/tweedproject/tweed/creds/volume"
)

var _ = Describe("Mounter", func() {
	var (
		mounter *Mounter
		secrets *credsfakes.FakeSecrets
		secret  string = "foo-secret"
		target  string
		ctx     context.Context
	)

	BeforeEach(func() {
		dir, err := ioutil.TempDir("", "target-dir")
		Expect(err).ToNot(HaveOccurred())
		target = dir
		secrets = &credsfakes.FakeSecrets{}
		mounter = NewMounter(secrets)
		ctx = context.Background()
	})

	Context("When secret does not exist", func() {
		It("Mounts an empty volume", func() {
			volume, err := mounter.Mount(ctx, target, secret)
			Expect(err).ToNot(HaveOccurred())
			err = volume.Unmount()
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("When given a secret", func() {
		BeforeEach(func() {
			secrets.GetReturns(map[string][]byte{"foo/bar": []byte("test")}, true, nil)
		})

		It("Can read secret from mounted volume", func() {
			volume, err := mounter.Mount(ctx, target, secret)
			Expect(err).ToNot(HaveOccurred())
			defer func() {
				err = volume.Unmount()
				Expect(err).ToNot(HaveOccurred())
			}()
			content, err := ioutil.ReadFile(path.Join(target, "foo/bar"))
			Expect(err).ToNot(HaveOccurred())
			Expect(string(content)).To(Equal("test"))
		})
	})

	AfterEach(func() {
		err := os.RemoveAll(target)
		Expect(err).ToNot(HaveOccurred())

	})

})
