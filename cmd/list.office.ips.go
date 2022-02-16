/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

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
var ipsCmd = &cobra.Command{
	Use:   "ips",
	Short: "ipv4 listing for azure",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		format, _ := cmd.Flags().GetString("format")
		switch format {
		case "raw":
			eps := GetOffice(officeRef)
			out := FormatIpRaw(eps)
			fmt.Println(out.String())
			break
		case "asa":
			eps := GetOffice(officeRef)
			out := FormatAsa(eps)
			fmt.Println(out.String())
			break
		default:
			fmt.Println(`unknown flag please use --help to see available options`)

		}

	},
}

func init() {
	officeCmd.AddCommand(ipsCmd)
	ipsCmd.PersistentFlags().StringP("format", "f", "raw", "output format")
	ipsCmd.MarkPersistentFlagRequired("format")
}

func FormatIpRaw(eps []OfficeEndpoint) bytes.Buffer {
	var out bytes.Buffer
	for _, v := range eps {
		if v.Ips != nil {
			for _, str := range v.Ips {
				if !strings.Contains(str, ":") {
					//ips = append(ips, str)
					out.WriteString(fmt.Sprintln(str))
				}

			}
		}
	}
	return out
}

func FormatAsa(eps []OfficeEndpoint) bytes.Buffer {
	var out bytes.Buffer
	var ips []string
	var list = make(map[string][]string)
	for _, v := range eps {
		name := fmt.Sprintf("O365.%v.%v", v.Id, v.ServiceArea)
		if v.Ips != nil {
			for _, str := range v.Ips {
				if !strings.Contains(str, ":") { // filter ipv6
					_, ipNet, _ := net.ParseCIDR(str)
					mask := ipv4MaskString(ipNet.Mask)
					ips = append(ips, fmt.Sprintf("%v/%v", ipNet.IP.String(), mask))
				}
			}
		}
		if len(ips) > 0 {
			list[name] = ips
		}
		fmt.Println((list))
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
