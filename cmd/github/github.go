package github

import (
	"context"
	"fmt"
	"log"

	"github.com/gofri/go-github-ratelimit/github_ratelimit"
	"github.com/google/go-github/v50/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type Repository github.Repository

type Workflow github.Workflow

//const githubDelay = 720 * time.Millisecond

type BranchProtectionRule struct {
	Nodes []struct {
		Id string
	}
	PageInfo struct {
		EndCursor   githubv4.String
		HasNextPage bool
	}
}

// type structure for the graphql query GetOrganizationsInEnterprise
type Organization struct {
	Nodes []struct {
		Login string
	}
	PageInfo struct {
		EndCursor   githubv4.String
		HasNextPage bool
	}
}

var (
	ctx         context.Context
	clientV3    *github.Client
	clientV4    *githubv4.Client
	accessToken string
)

func checkClients(token string, url string) error {

	// Sleep to avoid hitting the API rate limit.
	//time.Sleep(githubDelay)

	if clientV3 == nil || clientV4 == nil || token != accessToken {
		accessToken = token
		ctx = context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)
		rateLimiter, err := github_ratelimit.NewRateLimitWaiterClient(tc.Transport)

		if err != nil {
			return err
		}

		clientV3, err = github.NewEnterpriseClient(url, url, rateLimiter)
		if err != nil {
			return err
		}

		var query string
		if url == "https://api.github.com" {
			query = fmt.Sprintf("%s/graphql", url)
		} else {
			query = fmt.Sprintf("%s/api/graphql", url)
		}

		clientV4 = githubv4.NewEnterpriseClient(query, rateLimiter)
	}

	return nil
}

func logTokenRateLimit(response *github.Response) {
	log.Printf("Quota remaining: %d, Limit: %d, Reset: %s", response.Rate.Remaining, response.Rate.Limit, response.Rate.Reset)
}

func ChangeGHASOrgSettings(organization string, activate bool, token string, url string) error {
	checkClients(token, url)

	//create new organization object
	newOrgSettings := github.Organization{
		AdvancedSecurityEnabledForNewRepos:             &activate,
		SecretScanningPushProtectionEnabledForNewRepos: &activate,
		SecretScanningEnabledForNewRepos:               &activate,
	}

	// Update the organization
	_, response, err := clientV3.Organizations.Edit(ctx, organization, &newOrgSettings)

	logTokenRateLimit(response)

	if err != nil {
		log.Println("Error updating organization settings: ", err)
	}

	return err
}

func GetRepositories(org, token, url string) ([]Repository, error) {
	checkClients(token, url)

	// list all repositories for the organization
	opt := &github.RepositoryListByOrgOptions{Type: "all", ListOptions: github.ListOptions{PerPage: 10}}
	var allRepos []*github.Repository
	for {
		repos, response, err := clientV3.Repositories.ListByOrg(ctx, org, opt)

		logTokenRateLimit(response)

		if err != nil {
			log.Println("Error getting repositories: ", err)
			return nil, err
		}
		allRepos = append(allRepos, repos...)
		if response.NextPage == 0 {
			break
		}
		opt.Page = response.NextPage
	}

	var allReposStruct []Repository
	for _, repo := range allRepos {
		allReposStruct = append(allReposStruct, Repository(*repo))
	}

	return allReposStruct, nil
}

func GetRepository(repoName, org, token, url string) (Repository, error) {
	checkClients(token, url)

	repo, _, err := clientV3.Repositories.Get(ctx, org, repoName)
	if err != nil {
		log.Println("Error getting repository: ", err)
		return Repository{}, err
	}

	return Repository(*repo), nil
}

func GetOrganizationsInEnterprise(enterprise string, token string, url string) ([]string, error) {
	checkClients(token, url)

	var query struct {
		Enterprise struct {
			Organizations Organization `graphql:"organizations(first: 100, after: $cursor)"`
		} `graphql:"enterprise(slug: $enterprise)"`
	}

	variables := map[string]interface{}{
		"enterprise": githubv4.String(enterprise),
		"cursor":     (*githubv4.String)(nil),
	}

	results := make([]string, 0)
	for {
		err := clientV4.Query(ctx, &query, variables)
		if err != nil {
			log.Println("Error querying organizations: ", err)
			return nil, err
		}
		for _, organization := range query.Enterprise.Organizations.Nodes {
			results = append(results, organization.Login)
		}

		variables["cursor"] = query.Enterprise.Organizations.PageInfo.EndCursor

		if !query.Enterprise.Organizations.PageInfo.HasNextPage {
			break
		}
	}

	return results, nil
}
