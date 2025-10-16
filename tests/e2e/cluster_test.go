package e2e

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Grapevine Cluster", func() {

	It("should have ready nodes", func() {
		By("listing nodes")
		nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		Expect(err).NotTo(HaveOccurred())
		Expect(nodes.Items).ToNot(BeEmpty())

		By("checking for node readiness")
		for _, node := range nodes.Items {
			isReady := false
			for _, condition := range node.Status.Conditions {
				if condition.Type == "Ready" && condition.Status == "True" {
					isReady = true
					break
				}
			}
			Expect(isReady).To(BeTrue(), "Node %s should be ready", node.Name)
		}
	})

	Context("when deploying a sample application", func() {
		It("should run and be healthy", func() {
			By("creating a test namespace")
			namespace := createTestNamespace(clientset, "nginx-test")
			defer deleteTestNamespace(clientset, namespace)

			By("deploying an nginx application")
			deployment := createNginxDeployment(clientset, namespace.Name)

			By("waiting for the deployment to be available")
			Eventually(func() bool {
				dep, err := clientset.AppsV1().Deployments(namespace.Name).Get(context.TODO(), deployment.Name, metav1.GetOptions{})
				if err != nil {
					return false
				}
				return dep.Status.AvailableReplicas == 1
			}, time.Minute*2, time.Second*5).Should(BeTrue())

			By("verifying the pod is running")
			pods, err := clientset.CoreV1().Pods(namespace.Name).List(context.TODO(), metav1.ListOptions{LabelSelector: "app=nginx"})
			Expect(err).NotTo(HaveOccurred())
			Expect(pods.Items).To(HaveLen(1))
			Expect(pods.Items[0].Status.Phase).To(Equal(BeEquivalentTo("Running")))
		})
	})
})