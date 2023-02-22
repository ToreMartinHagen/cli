---
title: "kosli list approvals"
---

## kosli list approvals

List approvals in a flow.

### Synopsis

List approvals in a flow.
The results are paginated and ordered from latests to oldest. 
By default, the page limit is 15 approvals per page.  


```shell
kosli list approvals FLOW-NAME [flags]
```

### Flags
| Flag | Description |
| :--- | :--- |
|    -h, --help  |  help for approvals  |
|    -o, --output string  |  [defaulted] The format of the output. Valid formats are: [table, json]. (default "table")  |
|        --page int  |  [defaulted] The page number of a response. (default 1)  |
|    -n, --page-limit int  |  [defaulted] The number of elements per page. (default 15)  |


### Options inherited from parent commands
| Flag | Description |
| :--- | :--- |
|    -a, --api-token string  |  The Kosli API token.  |
|    -c, --config-file string  |  [optional] The Kosli config file path. (default "kosli")  |
|        --debug  |  [optional] Print debug logs to stdout.  |
|    -H, --host string  |  [defaulted] The Kosli endpoint. (default "https://app.kosli.com")  |
|    -r, --max-api-retries int  |  [defaulted] How many times should API calls be retried when the API host is not reachable. (default 3)  |
|        --owner string  |  The Kosli user or organization.  |


### Examples

```shell

# list the last 15 approvals for a flow:
kosli list approvals yourFlowName \
	--api-token yourAPIToken \
	--owner yourOrgName

# list the last 30 approvals for a flow:
kosli list approvals yourFlowName \
	--page-limit 30 \
	--api-token yourAPIToken \
	--owner yourOrgName

# list the last 30 approvals for a flow (in JSON):
kosli list approvals yourFlowName \
	--page-limit 30 \
	--api-token yourAPIToken \
	--owner yourOrgName \
	--output json

```

