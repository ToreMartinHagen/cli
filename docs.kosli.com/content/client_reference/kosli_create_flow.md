---
title: "kosli create flow"
---

# kosli create flow

## Synopsis

Create or update a Kosli flow.
You can provide a JSON pipefile or specify flow parameters in flags. 
The pipefile contains the flow metadata and compliance policy (template).

```shell
kosli create flow [FLOW-NAME] [flags]
```

## Flags
| Flag | Description |
| :--- | :--- |
|        --description string  |  [optional] The Kosli pipeline description.  |
|    -D, --dry-run  |  [optional] Run in dry-run mode. When enabled, no data is sent to Kosli and the CLI exits with 0 exit code regardless of any errors.  |
|    -h, --help  |  help for flow  |
|        --pipefile string  |  [deprecated] The path to the JSON pipefile.  |
|    -t, --template strings  |  [defaulted] The comma-separated list of required compliance controls names. (default [artifact])  |
|        --visibility string  |  [defaulted] The visibility of the Kosli pipeline. Valid visibilities are [public, private]. (default "private")  |


## Options inherited from parent commands
| Flag | Description |
| :--- | :--- |
|    -a, --api-token string  |  The Kosli API token.  |
|    -c, --config-file string  |  [optional] The Kosli config file path. (default "kosli")  |
|        --debug  |  [optional] Print debug logs to stdout.  |
|    -H, --host string  |  [defaulted] The Kosli endpoint. (default "https://app.kosli.com")  |
|    -r, --max-api-retries int  |  [defaulted] How many times should API calls be retried when the API host is not reachable. (default 3)  |
|        --owner string  |  The Kosli user or organization.  |


## Examples

```shell

# create/update a Kosli flow without a pipefile:
kosli create flow yourFlowName \
	--description yourFlowDescription \
    --visibility private OR public \
	--template artifact,evidence-type1,evidence-type2 \
	--api-token yourAPIToken \
	--owner yourOrgName

# create/update a Kosli flow with a pipefile (this is a legacy way which will be removed in the future):
kosli create flow yourFlowName \
	--pipefile /path/to/pipefile.json \
	--api-token yourAPIToken \
	--owner yourOrgName

The pipefile format is:
{
    "name": "yourFlowName",
    "description": "yourFlowDescription",
    "visibility": "public or private",
    "template": [
        "artifact",
        "evidence-type1",
        "evidence-type2"
    ]
}

```

