/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// urlsCmd represents the urls command
var urlsCmd = &cobra.Command{
	Use:   "urls",
	Short: "list urls/FQDNs",
	Long:  `get the list of urls / FQDNs  from the office 365 ip list.`,
	Run: func(cmd *cobra.Command, args []string) {
		format, _ := cmd.Flags().GetString("format")
		switch format {
		case "raw":
			fqdn := GetOffice(officeRef)
			out := FormatRaw(fqdn)
			fmt.Println(out.String())
			break
		case "proxy":
			fqdn := GetOffice(officeRef)
			out := FormatProxy(fqdn)
			fmt.Println(out.String())
			break
		default:
			fmt.Println(`unknown flag please use --help to see available options`)

		}
		//fmt.Println("urls called")
	},
}

func init() {
	officeCmd.AddCommand(urlsCmd)
	officeCmd.PersistentFlags().StringP("format", "f", "raw", "output format")
	officeCmd.MarkPersistentFlagRequired("format")

}

func FormatRaw(eps []OfficeEndpoint) bytes.Buffer {
	var out bytes.Buffer
	for _, v := range eps {
		if v.Urls != nil {
			for _, str := range v.Urls {
				out.WriteString(str + "\n")
			}
		}
	}
	return out
}

func FormatProxy(eps []OfficeEndpoint) bytes.Buffer {
	var out bytes.Buffer
	var list []string
	for _, v := range eps {
		if v.Urls != nil {
			for _, str := range v.Urls {
				list = append(list, str)
			}
		}
	}
	for key, val := range list {
		if strings.Contains(val, "*") {
			if key == (len(list) - 1) {
				out.WriteString(fmt.Sprintf("\tshExpMatch(url, '%v') \n", val))
			} else {
				out.WriteString(fmt.Sprintf("\tshExpMatch(url, '%v') || \n", val))
			}
		} else {
			if key == (len(list) - 1) {
				out.WriteString(fmt.Sprintf("\tdnsDomainIs(host,%v) \n", val))
			} else {
				out.WriteString(fmt.Sprintf("\tdnsDomainIs(host,%v) || \n", val))
			}
		}
	}
	return out

}
