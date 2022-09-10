package network_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	networkv1alpha1 "github.com/tbnco/tbnco/apis/network/v1alpha1"
	controller "github.com/tbnco/tbnco/controllers/network"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Network controller", func() {
	const (
		timeout  = time.Second * 5
		interval = time.Millisecond * 250
	)

	When("CR is instantiated", func() {
		const (
			namespace = "default"
			name      = "network"
		)

		var network *networkv1alpha1.Network
		var ctx context.Context

		BeforeEach(func() {
			ctx = context.Background()
			network = &networkv1alpha1.Network{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
			}

			Expect(k8sClient.Create(ctx, network)).Should(Succeed())
		})

		AfterEach(func() {
			Expect(k8sClient.Delete(ctx, network)).Should(Succeed())
		})

		It("should set a finalizer", func() {
			Eventually(func(g Gomega) {
				err := k8sClient.Get(ctx, client.ObjectKeyFromObject(network), network)
				g.Expect(err).ToNot(HaveOccurred())

				g.Expect(controller.NetworkFinalizer).ToNot(BeEmpty())
				g.Expect(network.GetFinalizers()).To(ContainElement(controller.NetworkFinalizer))
			}, timeout, interval).Should(Succeed())
		})
	})
})
