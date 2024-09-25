package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var dryC bool

var testteEmptyNamespacesCmd = &cobra.Command{
	Use:   "ns-clean",
	Short: "List or delete all namespaces that have no running pods",
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize Kubernetes client
		config, err := clientcmd.BuildConfigFromFlags("", "/Users/declan/.kube/config")
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

		// Iterate over each namespace and check if it has pods
		emptyNamespaces := []string{}
		for _, ns := range namespaces.Items {
			fmt.Printf("Checking namespace: %s\n", ns.Name) // 여기서 각 네임스페이스를 확인 중이라는 메시지 출력

			pods, err := clientset.CoreV1().Pods(ns.Name).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				fmt.Printf("Error getting pods for namespace %s: %s\n", ns.Name, err.Error())
				continue
			}

			// If no pods are found, consider the namespace empty
			if len(pods.Items) == 0 {
				emptyNamespaces = append(emptyNamespaces, ns.Name)
			}
		}

		// Print empty namespaces
		if len(emptyNamespaces) > 0 {
			fmt.Println("Namespaces with no running pods:")
			for _, ns := range emptyNamespaces {
				fmt.Printf("- %s\n", ns)
			}

			// If not in dry-run mode, prompt the user before deletion
			if !dryRun {
				fmt.Print("\nDo you want to delete all these namespaces? (yes/no): ")
				var input string
				fmt.Scanln(&input)

				// Normalize input to lowercase and check if it's "yes"
				if strings.ToLower(input) == "yes" {
					fmt.Println("\nDeleting namespaces with no running pods...")
					for _, ns := range emptyNamespaces {
						fmt.Printf("Deleting namespace: %s\n", ns)
						err := clientset.CoreV1().Namespaces().Delete(context.TODO(), ns, metav1.DeleteOptions{})
						if err != nil {
							fmt.Printf("Error deleting namespace %s: %s\n", ns, err.Error())
						} else {
							fmt.Printf("Successfully deleted namespace: %s\n", ns)
						}
					}
				} else {
					fmt.Println("\nNo namespaces were deleted.")
				}
			} else {
				fmt.Println("\nDry-run mode: no namespaces were deleted.")
			}
		} else {
			fmt.Println("No namespaces with zero running pods found.")
		}
	},
}

func init() {
	// Add --dry-run flag to prevent actual deletion
	deleteEmptyNamespacesCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "List namespaces without deleting them")
	rootCmd.AddCommand(deleteEmptyNamespacesCmd)
}
