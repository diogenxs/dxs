/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net/url"
	"strings"

	capi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
)

// consulCmd represents the consul command
var consulCmd = &cobra.Command{
	Use:   "consul",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("consul called")
	},
}

var consulSvcDeregisterCmd = &cobra.Command{
	Use:   "deregister",
	Short: "deregister a service from consul",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url, err := url.Parse(args[0])
		if err != nil {
			fmt.Println("pass the entire url of consul ui")
			return
		}

		urlParts := strings.Split(url.Path, "/")

		dc := urlParts[2]
		service := urlParts[4]
		node := urlParts[6]
		instance := urlParts[7]

		fmt.Println("deregistering in ...")
		fmt.Println("datacenter: ", dc)
		fmt.Println("service: ", service)
		fmt.Println("node: ", node)
		fmt.Println("instance: ", instance)

		client, err := capi.NewClient(&capi.Config{
			Address:    url.Host,
			Scheme:     string(url.Scheme),
			Datacenter: dc,
		})
		// client, err = capi.NewClient(capi.DefaultConfig())
		if err != nil {
			fmt.Println("error creating consul client")
			panic(err)
		}
		// resp, _, err := client.Health().Node(node, &capi.QueryOptions{})
		// fmt.Println(resp.AggregatedStatus())
		_, err = client.Catalog().Deregister(&capi.CatalogDeregistration{
			Node:      node,
			ServiceID: service,
			CheckID:   instance,
		}, &capi.WriteOptions{})

		if err != nil {
			fmt.Println("error deregistering consul client")
			panic(err)
		}
		fmt.Println("service deregistered successfully")
		// for _, v := range resp {
		// 	fmt.Println(v.ID)
		// 	fmt.Println(v.Service)
		// 	fmt.Println(v.Tags)
		// 	fmt.Println(v.Port)
		// }
	},
}

func init() {
	rootCmd.AddCommand(consulCmd)

	consulCmd.AddCommand(consulSvcDeregisterCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// consulCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// consulCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
