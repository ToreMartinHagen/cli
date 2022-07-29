package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kosli-dev/cli/internal/requests"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const environmentGetDesc = `Get an environment metadata.`

type environmentGetOptions struct {
	json bool
}

func newEnvironmentGetCmd(out io.Writer) *cobra.Command {
	o := new(environmentGetOptions)
	cmd := &cobra.Command{
		Use:   "get [ENVIRONMENT-NAME]",
		Short: environmentGetDesc,
		Long:  environmentGetDesc,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := RequireGlobalFlags(global, []string{"Owner", "ApiToken"})
			if err != nil {
				return ErrorAfterPrintingHelp(cmd, err.Error())
			}
			if len(args) < 1 {
				return ErrorAfterPrintingHelp(cmd, "environment name argument is required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.run(out, args)
		},
	}

	cmd.Flags().BoolVarP(&o.json, "json", "j", false, jsonOutputFlag)

	return cmd
}

func (o *environmentGetOptions) run(out io.Writer, args []string) error {
	url := fmt.Sprintf("%s/api/v1/environments/%s/%s", global.Host, global.Owner, args[0])
	response, err := requests.DoBasicAuthRequest([]byte{}, url, "", global.ApiToken,
		global.MaxAPIRetries, http.MethodGet, map[string]string{}, logrus.New())
	if err != nil {
		return err
	}

	if o.json {
		pj, err := prettyJson(response.Body)
		if err != nil {
			return err
		}
		fmt.Println(pj)
		return nil
	}

	var env map[string]interface{}
	err = json.Unmarshal([]byte(response.Body), &env)
	if err != nil {
		return err
	}

	// last_reported_str := ""
	// last_reported_at := env["last_reported_at"]
	// if last_reported_at != nil {
	// 	last_reported_str = time.Unix(int64(last_reported_at.(float64)), 0).Format(time.RFC3339)
	// }
	// last_modified_str := ""
	// last_modified_at := env["last_modified_at"]
	// if last_modified_at != nil {
	// 	last_modified_str = time.Unix(int64(last_modified_at.(float64)), 0).Format(time.RFC3339)
	// }

	lastReportedAt, err := formattedTimestamp(env["last_reported_at"], false)
	if err != nil {
		return err
	}

	state := "N/A"
	if env["state"] != nil && env["state"].(bool) {
		state = "COMPLIANT"
	} else if env["state"] != nil {
		state = "INCOMPLIANT"
	}

	header := []string{}
	rows := []string{}
	rows = append(rows, fmt.Sprintf("Name:\t%s", env["name"]))
	rows = append(rows, fmt.Sprintf("Type:\t%s", env["type"]))
	rows = append(rows, fmt.Sprintf("Description:\t%s", env["description"]))
	rows = append(rows, fmt.Sprintf("State:\t%s", state))
	rows = append(rows, fmt.Sprintf("Last Reported At:\t%s", lastReportedAt))

	printTable(out, header, rows)

	return nil
}
