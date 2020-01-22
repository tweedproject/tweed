package kubernetes_test

import (
	"code.cloudfoundry.org/lager/lagertest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tweedproject/tweed/creds"
	"github.com/tweedproject/tweed/creds/kubernetes"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

var _ = Describe("Kubernetes", func() {
	var fakeClientset *fake.Clientset
	var secretName = "some-secret-name"
	var secrets creds.Secrets

	BeforeEach(func() {
		fakeClientset = fake.NewSimpleClientset()

		factory := kubernetes.NewKubernetesFactory(
			lagertest.NewTestLogger("test"),
			fakeClientset,
			"secrets-namespace",
		)

		secrets = factory.NewSecrets()
	})

	Describe("Get()", func() {
		It("should get a secret by name in a given namespace", func() {
			fakeClientset.CoreV1().Secrets("secrets-namespace").Create(&v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: secretName,
				},
				Data: map[string][]byte{
					"value": []byte("some-value"),
				},
			})

			value, exist, err := secrets.Get(secretName)
			Expect(err).ToNot(HaveOccurred())
			Expect(exist).To(BeTrue())
			Expect(value).To(Equal("some-value"))
		})

		It("should get a non string secret by name in a given namespace", func() {
			fakeClientset.CoreV1().Secrets("secrets-namespace").Create(&v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: secretName,
				},
				Data: map[string][]byte{
					"foo": []byte("bar"),
				},
			})

			value, exist, err := secrets.Get(secretName)
			Expect(err).ToNot(HaveOccurred())
			Expect(exist).To(BeTrue())
			Expect(value).To(BeEquivalentTo(map[string]interface{}{"foo": "bar"}))
		})

		It("when secret does not exist", func() {
			value, exist, err := secrets.Get(secretName)
			Expect(err).ToNot(HaveOccurred())
			Expect(exist).To(BeFalse())
			Expect(value).To(BeNil())
		})
	})

	Describe("Set()", func() {
		It("should set a secret by name in a given namespace", func() {
			value := "data"
			err := secrets.Set(secretName, value)
			Expect(err).ToNot(HaveOccurred())

			secret, err := fakeClientset.CoreV1().Secrets("secrets-namespace").
				Get(secretName, metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())
			val, found := secret.Data["value"]
			Expect(found).To(BeTrue())
			Expect(string(val)).To(Equal("data"))
		})

		It("should update a secret which already exists", func() {
			fakeClientset.CoreV1().Secrets("secrets-namespace").Create(&v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: secretName,
				},
				Data: map[string][]byte{
					"bar": []byte("bar"),
				},
			})

			value := "data"
			err := secrets.Set(secretName, value)
			Expect(err).ToNot(HaveOccurred())

			secret, err := fakeClientset.CoreV1().Secrets("secrets-namespace").
				Get(secretName, metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())
			val, found := secret.Data["value"]
			Expect(found).To(BeTrue())
			Expect(string(val)).To(Equal("data"))
		})

		It("should set a secret to a non string value", func() {
			value := map[string]string{"foo": "bar"}
			err := secrets.Set(secretName, value)
			Expect(err).ToNot(HaveOccurred())

			secret, err := fakeClientset.CoreV1().Secrets("secrets-namespace").
				Get(secretName, metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())
			val, found := secret.Data["foo"]
			Expect(found).To(BeTrue())
			Expect(string(val)).To(Equal("bar"))
		})

	})
})
