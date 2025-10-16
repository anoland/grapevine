package e2e

import (
	"context"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Grapevine E2E Suite")
}

var (
	kubeconfig string
	clientset  *kubernetes.Clientset
	terraformOptions *terraform.Options
)

var _ = BeforeSuite(func() {
	token := os.Getenv("RACKSPACE_TOKEN")
	Expect(token).ToNot(BeEmpty(), "RACKSPACE_TOKEN environment variable must be set")

	terraformOptions = &terraform.Options{
		TerraformDir: "./terraform",
		Vars: map[string]interface{}{
			"rackspace_token": token,
		},
	}

	By("initializing and applying terraform")
	terraform.InitAndApply(GinkgoT(), terraformOptions)

	By("getting kubeconfig from terraform outputs")
	kubeconfig = terraform.Output(GinkgoT(), terraformOptions, "kubeconfig")

	By("creating kubernetes clientset")
	config, err := clientcmd.NewClientConfigFromBytes([]byte(kubeconfig))
	Expect(err).NotTo(HaveOccurred())

	restConfig, err := config.ClientConfig()
	Expect(err).NotTo(HaveOccurred())

	clientset, err = kubernetes.NewForConfig(restConfig)
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	By("destroying terraform resources")
	terraform.Destroy(GinkgoT(), terraformOptions)
})

// Helper function to create a test namespace
func createTestNamespace(clientset *kubernetes.Clientset, baseName string) *corev1.Namespace {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: baseName + "-",
		},
	}
	ns, err := clientset.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
	Expect(err).NotTo(HaveOccurred())
	return ns
}

// Helper function to delete a namespace
func deleteTestNamespace(clientset *kubernetes.Clientset, namespace *corev1.Namespace) {
	err := clientset.CoreV1().Namespaces().Delete(context.TODO(), namespace.Name, metav1.DeleteOptions{})
	Expect(err).NotTo(HaveOccurred())
}

// Helper function to create an Nginx deployment
func createNginxDeployment(clientset *kubernetes.Clientset, namespace string) *appsv1.Deployment {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "nginx",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "nginx",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx:1.21.6",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	By("creating nginx deployment")
	dep, err := clientset.AppsV1().Deployments(namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	Expect(err).NotTo(HaveOccurred())
	return dep
}

func int32Ptr(i int32) *int32 { return &i }