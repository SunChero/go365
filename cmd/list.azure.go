/*
Copyright © 2022 Adil han

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var azureRef = "https://download.microsoft.com/download/7/1/D/71D86715-5596-4529-9B13-DA13A5DE5B63/ServiceTags_Public_20220207.json"

type AzureEndpoint struct {
	Values []struct {
		Id         string `json:"id"`
		Name       string `json:"name"`
		Properties struct {
			Region          string   `json:"region"`
			RegionId        int      `json:"regionId"`
			AddressPrefixes []string `json:"addressPrefixes"`
		} `json:"properties"`
	} `json:"values"`
}

// azureCmd represents the azure command
var azureCmd = &cobra.Command{
	Use:   "azure",
	Short: "list azure ips ",
	Long: `list azure ip addresses
	ips are filtered to only show ipv4 listing
	
	command also support different formats
	for more info type --help.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`unknown flag please use --help to see available options`)
	},
}

func init() {
	listCmd.AddCommand(azureCmd)
	azureCmd.PersistentFlags().StringP("format", "f", "raw", "output format")
	azureCmd.MarkPersistentFlagRequired("format")

}

func GetAzure(ref string) AzureEndpoint {
	var azep AzureEndpoint
	res, err := http.Get(ref)
	if err != nil {
		fmt.Errorf("Could not get the link from microsoft: %v", err)
	}
	defer res.Body.Close()
	json.NewDecoder(res.Body).Decode(&azep)
	return azep
}
