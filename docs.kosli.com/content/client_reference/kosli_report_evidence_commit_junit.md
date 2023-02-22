---
title: "kosli report evidence commit junit"
---

## kosli report evidence commit junit

Report JUnit test evidence for a commit in a Kosli flow.

### Synopsis

Report JUnit test evidence for an artifact in a Kosli flow.

```shell
kosli report evidence commit junit [flags]
```

### Flags
| Flag | Description |
| :--- | :--- |
|    -b, --build-url string  |  The url of CI pipeline that generated the evidence. (defaulted in some CIs: https://docs.kosli.com/ci-defaults ).  |
|        --commit string  |  The git commit SHA1 for which the evidence belongs. (defaulted in some CIs: https://docs.kosli.com/ci-defaults ).  |
|    -D, --dry-run  |  [optional] Run in dry-run mode. When enabled, no data is sent to Kosli and the CLI exits with 0 exit code regardless of any errors.  |
|    -f, --flow strings  |  The comma separated list of pipelines for which a commit evidence belongs.  |
|    -h, --help  |  help for junit  |
|    -n, --name string  |  The name of the evidence.  |
|    -R, --results-dir string  |  [defaulted] The path to a folder with JUnit test results. (default ".")  |
|    -u, --user-data string  |  [optional] The path to a JSON file containing additional data you would like to attach to this evidence.  |


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

# report JUnit test evidence for a commit related to one Kosli flow:
kosli report evidence commit junit \
	--commit yourGitCommitSha1 \
	--name yourEvidenceName \
	--flow yourFlowName \
	--build-url https://exampleci.com \
	--api-token yourAPIToken \
	--owner yourOrgName	\
	--results-dir yourFolderWithJUnitResults

# report JUnit test evidence for a commit related to multiple Kosli flows:
kosli report evidence commit junit \
	--commit yourGitCommitSha1 \
	--name yourEvidenceName \
	--flow yourFlowName1,yourFlowName2 \
	--build-url https://exampleci.com \
	--api-token yourAPIToken \
	--owner yourOrgName	\
	--results-dir yourFolderWithJUnitResults

```

