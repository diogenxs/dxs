/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/diogenxs/dxs/pkg/k8s"
	"github.com/spf13/cobra"
)

// k8sCmd represents the k8s command
var k8sCmd = &cobra.Command{
	Use:     "k8s",
	Short:   "Useful k8s operations",
	Aliases: []string{"k"},
}

var listPendingPodsCmd = &cobra.Command{
	Use:     "lpp [namespace]",
	Short:   "List pending pods",
	Aliases: []string{"lpp"},
	Run: func(cmd *cobra.Command, args []string) {
		var namespace string

		if len(args) > 0 {
			namespace = args[0]
		}

		fmt.Println("listing pending pods...")
		client, err := k8s.NewK8sClient()
		if err != nil {
			fmt.Println(err)
		}

		pods, err := client.ListPendingPods(namespace)
		if err != nil {
			fmt.Println(err)
		}
		for _, pod := range pods.Items {
			fmt.Println(pod.Name)
			// fmt.Println(k8s.PrettyPrintJSON(pod))
		}
	},
}

var listNodesCmd = &cobra.Command{
	Use:     "listNodes",
	Short:   "List nodes",
	Aliases: []string{"ln"},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("listing nodes...")
		client, err := k8s.NewK8sClient()
		if err != nil {
			return err
		}
		label, err := cmd.Flags().GetString("label")
		if err != nil {
			return err
		}
		n, err := client.ListNodesByLabel(label)
		if err != nil {
			return err
		}
		fmt.Printf("found %d nodes\n", len(n.Items))
		for _, node := range n.Items {
			fmt.Println(node.Name)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(k8sCmd)

	k8sCmd.AddCommand(listPendingPodsCmd)
	k8sCmd.AddCommand(listNodesCmd)

	// global k8sCmd flags
	k8sCmd.PersistentFlags().String("context", "", "same as used with 'kubectl --context'")

	// specific subcommands flags
	listNodesCmd.PersistentFlags().StringP("label", "l", "", "same as used with 'kubectl get nodes -l'")
}
