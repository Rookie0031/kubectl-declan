package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var dryRun bool

// Check if the namespace has Pods, Services, and other resources
func isNamespaceEmpty(clientset *kubernetes.Clientset, namespace string) (bool, map[string]int, error) {
	// Initialize a map to store resource counts
	resourceCounts := make(map[string]int)

	// Check for Pods
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false, resourceCounts, err
	}
	resourceCounts["Pods"] = len(pods.Items)

	// Check for Services
	services, err := clientset.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false, resourceCounts, err
	}
	resourceCounts["Services"] = len(services.Items)

	// Check for ConfigMaps
	configMaps, err := clientset.CoreV1().ConfigMaps(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false, resourceCounts, err
	}
	resourceCounts["ConfigMaps"] = len(configMaps.Items)

	// Check for Secrets
	secrets, err := clientset.CoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false, resourceCounts, err
	}
	resourceCounts["Secrets"] = len(secrets.Items)

	// Check for Deployments
	deployments, err := clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false, resourceCounts, err
	}
	resourceCounts["Deployments"] = len(deployments.Items)

	// Check for ReplicaSets
	replicaSets, err := clientset.AppsV1().ReplicaSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false, resourceCounts, err
	}
	resourceCounts["ReplicaSets"] = len(replicaSets.Items)

	// Check for StatefulSets
	statefulSets, err := clientset.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false, resourceCounts, err
	}
	resourceCounts["StatefulSets"] = len(statefulSets.Items)

	// Check for DaemonSets
	daemonSets, err := clientset.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false, resourceCounts, err
	}
	resourceCounts["DaemonSets"] = len(daemonSets.Items)

	// Check for Jobs
	jobs, err := clientset.BatchV1().Jobs(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false, resourceCounts, err
	}
	resourceCounts["Jobs"] = len(jobs.Items)

	// Check for CronJobs
	cronJobs, err := clientset.BatchV1().CronJobs(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false, resourceCounts, err
	}
	resourceCounts["CronJobs"] = len(cronJobs.Items)

	// Check for PersistentVolumeClaims
	pvcs, err := clientset.CoreV1().PersistentVolumeClaims(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false, resourceCounts, err
	}
	resourceCounts["PersistentVolumeClaims"] = len(pvcs.Items)

	// Check for Ingresses
	ingresses, err := clientset.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false, resourceCounts, err
	}
	resourceCounts["Ingresses"] = len(ingresses.Items)

	// If no Pods and Services are found, the namespace is considered empty of primary resources
	return resourceCounts["Pods"] == 0 && resourceCounts["Services"] == 0, resourceCounts, nil
}

var deleteEmptyNamespacesCmd = &cobra.Command{
	Use:   "ns-clean",
	Short: "List or delete all namespaces that have no Pods and Services",
	Run: func(cmd *cobra.Command, args []string) {
		// Get kubeconfig path from environment variable or use default
		kubeconfig := os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				log.Fatalf("Error getting user home directory: %v", err)
			}
			kubeconfig = filepath.Join(homeDir, ".kube", "config")
		}
		// Initialize Kubernetes client
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Fatalf("Error building kubeconfig: %s", err.Error())
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			log.Fatalf("Error creating Kubernetes client: %s", err.Error())
		}

		// Get all namespaces
		namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatalf("Error getting namespaces: %s", err.Error())
		}

		// Iterate over each namespace and check if it has any resources
		for _, ns := range namespaces.Items {
			fmt.Printf("Checking namespace: %s\n", ns.Name)

			// Check if the namespace is empty (no Pods and Services)
			empty, resourceCounts, err := isNamespaceEmpty(clientset, ns.Name)
			if err != nil {
				fmt.Printf("Error checking namespace %s: %s\n", ns.Name, err.Error())
				continue
			}

			if empty {
				fmt.Printf("Namespace %s is empty of Pods and Services.\n", ns.Name)

				// Show the remaining resources
				fmt.Printf("Resources remaining in namespace %s:\n", ns.Name)
				for resourceType, count := range resourceCounts {
					fmt.Printf("- %s: %d\n", resourceType, count)
				}

				// If not in dry-run mode, prompt the user before deleting the namespace
				if !dryRun {
					fmt.Printf("\nDo you want to delete the namespace '%s'? (yes/no): ", ns.Name)
					var input string
					fmt.Scanln(&input)

					// Normalize input to lowercase and check if it's "yes"
					if strings.ToLower(input) == "yes" {
						fmt.Printf("Deleting namespace: %s\n", ns.Name)
						err := clientset.CoreV1().Namespaces().Delete(context.TODO(), ns.Name, metav1.DeleteOptions{})
						if err != nil {
							fmt.Printf("Error deleting namespace %s: %s\n", ns.Name, err.Error())
						} else {
							fmt.Printf("Successfully deleted namespace: %s\n", ns.Name)
						}
					} else {
						fmt.Printf("Skipped deletion of namespace: %s\n", ns.Name)
					}
				}
			}
		}

		if dryRun {
			fmt.Println("\nDry-run mode: no namespaces were deleted.")
		}
	},
}

func init() {
	// Add --dry-run flag to prevent actual deletion
	if deleteEmptyNamespacesCmd.Flags().Lookup("dry-run") == nil {
		deleteEmptyNamespacesCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "List namespaces without deleting them")
	}
	rootCmd.AddCommand(deleteEmptyNamespacesCmd)
}
