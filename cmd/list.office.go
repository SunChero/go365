/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

type OfficeEndpoint struct {
	Id                     int      `json:"id"`
	ServiceArea            string   `json:"serviceArea"`
	ServiceAreaDisplayName string   `json:"serviceAreaDisplayName"`
	Urls                   []string `json:"urls"`
	TcpPorts               string   `json:"tcpPorts"`
	Ips                    []string `json:"ips"`
	Notes                  string   `json:"notes"`
}

var officeRef = "https://endpoints.office.com/endpoints/worldwide?clientrequestid=b10c5ed1-bad1-445f-b386-b919946339a7"

// officeCmd represents the office command
var officeCmd = &cobra.Command{
	Use:   "office",
	Short: "list office 365 endpoints",
	Long: `list ip addresses and FQDNs that are part 
	of the office365.`,
	Run: func(cmd *cobra.Command, args []string) {

		//fmt.Println("office called")
	},
}

func GetOffice(ref string) []OfficeEndpoint {
	var oep []OfficeEndpoint
	res, err := http.Get(ref)
	if err != nil {
		fmt.Errorf("Could not get the link from microsoft: %v", err)
	}
	defer res.Body.Close()
	json.NewDecoder(res.Body).Decode(&oep)
	return oep
}

func init() {
	listCmd.AddCommand(officeCmd)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// officeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
