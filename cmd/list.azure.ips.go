/*
Copyright © 2022 Adil ha

*/
package cmd

import (
	"bytes"
	"fmt"
	"net"
	"strings"

	"github.com/spf13/cobra"
)

// ipsCmd represents the ips command
var zipsCmd = &cobra.Command{
	Use:   "ips",
	Short: "ipv4 listing for azure",
	Long: `list azure ip addresses
	support different output formats
	type --help for more infomations`,
	Run: func(cmd *cobra.Command, args []string) {
		format, _ := cmd.Flags().GetString("format")
		switch format {
		case "raw":
			eps := GetAzure(azureRef)
			out := FormatAzIpRaw(eps)
			fmt.Println(out.String())
			break
		case "asa":
			eps := GetAzure(azureRef)
			out := FormatAzAsa(eps)
			fmt.Println(out.String())
			break
		default:
			fmt.Println(`unknown flag please use --help to see available options`)

		}
		//fmt.Println("ips called")
	},
}

func init() {
	azureCmd.AddCommand(zipsCmd)
	zipsCmd.PersistentFlags().StringP("format", "f", "raw", "output format")
	zipsCmd.MarkPersistentFlagRequired("format")

}

func FormatAzIpRaw(eps AzureEndpoint) bytes.Buffer {
	var out bytes.Buffer
	for _, v := range eps.Values {
		if v.Properties.AddressPrefixes != nil {
			for _, str := range v.Properties.AddressPrefixes {
				if !strings.Contains(str, ":") { // filter ipv6
					_, ipNet, _ := net.ParseCIDR(str)
					mask := ipv4MaskString(ipNet.Mask)
					out.WriteString(fmt.Sprintf("%v - %v \n", ipNet.IP.String(), mask))
					//ips = append(ips, fmt.Sprintf("%v - %v", ipNet.IP.String(), mask))
				}
			}
		}
	}
	return out
}

func FormatAzAsa(eps AzureEndpoint) bytes.Buffer {
	var out bytes.Buffer
	var ips []string
	var list = make(map[string][]string)
	for _, v := range eps.Values {
		//access to rules name
		name := fmt.Sprintf("AZ.%v.%v", v.Properties.RegionId, v.Name)
		if v.Properties.AddressPrefixes != nil {
			for _, str := range v.Properties.AddressPrefixes {
				if !strings.Contains(str, ":") { // filter ipv6
					_, ipNet, _ := net.ParseCIDR(str)
					mask := ipv4MaskString(ipNet.Mask)
					ips = append(ips, fmt.Sprintf("%v/%v", ipNet.IP.String(), mask))
				}
			}
		}
		list[name] = ips
		ips = ips[:0]

	}
	for key, value := range list {
		out.WriteString(fmt.Sprintf("no object-group network %v \n", key))
		out.WriteString(fmt.Sprintf("object-group network %v \n", key))
		for _, v := range value {
			parts := strings.Split(v, "/")
			if parts[1] == "255.255.255.255" {
				out.WriteString(fmt.Sprintf("network-object host %v  %v\n", parts[0], parts[1]))
			} else {
				out.WriteString(fmt.Sprintf("network-object %v  %v\n", parts[0], parts[1]))
			}
		}

	}
	return out

}
