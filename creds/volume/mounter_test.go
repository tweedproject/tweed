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
		var (
			secretFile string
		)
		BeforeEach(func() {
			secrets.GetReturns(map[string][]byte{"foo/bar": []byte("test")}, true, nil)
			secretFile = path.Join(target, "foo/bar")
		})

		It("Can read secret from mounted volume", func() {
			volume, err := mounter.Mount(ctx, target, secret)
			Expect(err).ToNot(HaveOccurred())
			defer func() {
				err = volume.Unmount()
				Expect(err).ToNot(HaveOccurred())
			}()
			content, err := ioutil.ReadFile(secretFile)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(content)).To(Equal("test"))
		})

		It("Can update secret on the mounted volume", func() {
			volume, err := mounter.Mount(ctx, target, secret)
			Expect(err).ToNot(HaveOccurred())
			defer func() {
				err = volume.Unmount()
				Expect(err).ToNot(HaveOccurred())
			}()
			err = ioutil.WriteFile(secretFile, []byte("updated"), 0775)
			Expect(err).ToNot(HaveOccurred())
			Expect(secrets.SetCallCount()).To(Equal(1))
			path, data := secrets.SetArgsForCall(0)
			Expect(path).To(Equal(secret))
			Expect(data).To(Equal(map[string][]byte{"foo/bar": []byte("updated")}))
			content, err := ioutil.ReadFile(secretFile)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(content)).To(Equal("updated"))

		})

	})

	AfterEach(func() {
		err := os.RemoveAll(target)
		Expect(err).ToNot(HaveOccurred())

	})

})
