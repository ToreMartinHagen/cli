package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/kosli-dev/cli/internal/gitview"
	"github.com/kosli-dev/cli/internal/jira"
	"github.com/kosli-dev/cli/internal/requests"
	"github.com/spf13/cobra"
)

type reportEvidenceCommitJiraOptions struct {
	userDataFilePath string
	evidencePaths    []string
	baseURL          string
	username         string
	apiToken         string
	pat              string
	srcRepoRoot      string
	payload          JiraEvidencePayload
}

const reportEvidenceCommitJiraShortDesc = `Report Jira evidence for a commit in Kosli flows.`

const reportEvidenceCommitJiraLongDesc = reportEvidenceCommitJiraShortDesc + `
Parses the given commit's message or current branch name for Jira issue references of the 
form: 'at least 2 characters long, starting with an uppercase letter project key followed by
dash and one or more digits'.
The found issue references will be checked against Jira to confirm their existence.
The evidence is reported in all cases, and its compliance status depends on referencing
existing Jira issues.
`

const reportEvidenceCommitJiraExample = `
# report Jira evidence for a commit related to one Kosli flow (with Jira Cloud):
kosli report evidence commit jira \
	--commit yourGitCommitSha1 \
	--name yourEvidenceName \
	--jira-base-url https://kosli.atlassian.net \
	--jira-username user@domain.com \
	--jira-api-token yourJiraAPIToken \
	--flows yourFlowName \
	--build-url https://exampleci.com \
	--api-token yourAPIToken \
	--org yourOrgName

# report Jira evidence for a commit related to one Kosli flow (with self-hosted Jira):
kosli report evidence commit jira \
	--commit yourGitCommitSha1 \
	--name yourEvidenceName \
	--jira-base-url https://kosli.atlassian.net \
	--jira-pat yourJiraPATToken \
	--flows yourFlowName \
	--build-url https://exampleci.com \
	--api-token yourAPIToken \
	--org yourOrgName

# report Jira  evidence for a commit related to multiple Kosli flows with user-data (with Jira Cloud):
kosli report evidence commit jira \
	--commit yourGitCommitSha1 \
	--name yourEvidenceName \
	--jira-base-url https://kosli.atlassian.net \
	--jira-username user@domain.com \
	--jira-api-token yourJiraAPIToken \
	--flows yourFlowName1,yourFlowName2 \
	--build-url https://exampleci.com \
	--api-token yourAPIToken \
	--org yourOrgName \
	--user-data /path/to/json/file.json
`

func newReportEvidenceCommitJiraCmd(out io.Writer) *cobra.Command {
	o := new(reportEvidenceCommitJiraOptions)
	cmd := &cobra.Command{
		Use:     "jira",
		Short:   reportEvidenceCommitJiraShortDesc,
		Long:    reportEvidenceCommitJiraLongDesc,
		Example: reportEvidenceCommitJiraExample,
		Args:    cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := RequireGlobalFlags(global, []string{"Org", "ApiToken"})
			if err != nil {
				return ErrorBeforePrintingUsage(cmd, err.Error())
			}

			err = MuXRequiredFlags(cmd, []string{"jira-pat", "jira-api-token"}, true)
			if err != nil {
				return err
			}

			err = MuXRequiredFlags(cmd, []string{"jira-pat", "jira-username"}, true)
			if err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.run(args)
		},
	}

	ci := WhichCI()
	addCommitEvidenceFlags(cmd, &o.payload.TypedEvidencePayload, ci)
	cmd.Flags().StringVar(&o.baseURL, "jira-base-url", "", jiraBaseUrlFlag)
	cmd.Flags().StringVar(&o.username, "jira-username", "", jiraUsernameFlag)
	cmd.Flags().StringVar(&o.apiToken, "jira-api-token", "", jiraAPITokenFlag)
	cmd.Flags().StringVar(&o.pat, "jira-pat", "", jiraPATFlag)
	cmd.Flags().StringVar(&o.srcRepoRoot, "repo-root", ".", repoRootFlag)
	cmd.Flags().StringVarP(&o.userDataFilePath, "user-data", "u", "", evidenceUserDataFlag)
	cmd.Flags().StringSliceVarP(&o.evidencePaths, "evidence-paths", "e", []string{}, evidencePathsFlag)
	addDryRunFlag(cmd)

	err := RequireFlags(cmd, []string{"commit", "build-url", "name", "jira-base-url"})
	if err != nil {
		logger.Error("failed to configure required flags: %v", err)
	}

	return cmd
}

func (o *reportEvidenceCommitJiraOptions) run(args []string) error {
	var err error

	o.baseURL = strings.TrimSuffix(o.baseURL, "/")

	jc := jira.NewJiraConfig(o.baseURL, o.username, o.apiToken, o.pat)

	url := fmt.Sprintf("%s/api/v2/evidence/%s/commit/jira", global.Host, global.Org)
	o.payload.UserData, err = LoadJsonData(o.userDataFilePath)
	if err != nil {
		return err
	}

	gv, err := gitview.New(o.srcRepoRoot)
	if err != nil {
		return err
	}

	o.payload.CommitSHA, err = gv.ResolveRevision(o.payload.CommitSHA)
	if err != nil {
		return err
	}

	o.payload.JiraResults = []*jira.JiraIssueInfo{}

	// Jira issue keys consist of [project-key]-[sequential-number]
	// project key must be at least 2 characters long and start with an uppercase letter
	// more info: https://support.atlassian.com/jira-software-cloud/docs/what-is-an-issue/#Workingwithissues-Projectandissuekeys
	jiraIssueKeyPattern := `[A-Z][A-Za-z]+-[0-9]+`

	issueIDs, err := gv.MatchPatternInCommitMessageORBranchName(jiraIssueKeyPattern, o.payload.CommitSHA)
	if err != nil {
		return err
	}

	logger.Debug("the following Jira references are found in commit message or branch name: %v", issueIDs)

	for _, issueID := range issueIDs {
		result, err := jc.GetJiraIssueInfo(issueID)
		if err != nil {
			return err
		}
		o.payload.JiraResults = append(o.payload.JiraResults, result)
	}

	form, cleanupNeeded, evidencePath, err := newEvidenceForm(o.payload, o.evidencePaths)
	// if we created a tar package, remove it after uploading it
	if cleanupNeeded {
		defer os.Remove(evidencePath)
	}

	if err != nil {
		return err
	}

	reqParams := &requests.RequestParams{
		Method:   http.MethodPost,
		URL:      url,
		Form:     form,
		DryRun:   global.DryRun,
		Password: global.ApiToken,
	}

	_, err = kosliClient.Do(reqParams)
	if err == nil && !global.DryRun {
		logger.Info("jira evidence is reported to commit: %s", o.payload.CommitSHA)
	}
	return err
}
