package cmd

import (
	"io/ioutil"
	"log"

	"github.com/gateixeira/ghas-org-enablement/cmd/github"
)

func ManageGhasFeatures(enterprise, organization, token, url string, activate bool) error {

	var organizations []string
	if organization != "" {
		organizations = []string{organization}
	} else {
		log.Println("[🔄] Fetching organizations from enterprise...")
		var err error
		organizations, err = github.GetOrganizationsInEnterprise(enterprise, token, url)

		if err != nil {
			log.Println("[❌] Error fetching organizations from enterprise")
			return err
		}

		log.Println("[✅] Done")
	}

	for _, organization := range organizations {
		log.Println("[🔄] Changing GHAS settings for organization: " + organization)
		err := github.ChangeGHASOrgSettings(organization, activate, token, url)

		if err != nil {
			log.Println("[❌] Error changing GHAS settings for organization: " + organization)
			continue
		}
		log.Println("[✅] Done")
	}

	return nil
}

func ListSecretScanning(enterprise, organization, token, url string, activate bool) error {

	var organizations []string
	if organization != "" {
		organizations = []string{organization}
	} else {
		log.Println("[🔄] Fetching organizations from enterprise...")
		var err error
		organizations, err = github.GetOrganizationsInEnterprise(enterprise, token, url)

		if err != nil {
			log.Println("[❌] Error fetching organizations from enterprise")
			return err
		}

		log.Println("[✅] Done")
	}

	mapOrgAndRepos := map[string][]string{}
	totalRepos := 0
	for _, organization := range organizations {
		log.Println("[🔄] Listing secret scanning for organization: " + organization)

		repos, err := github.GetRepositories(organization, token, url)
		totalRepos += len(repos)
		log.Println("Total repositories found in organization: ", len(repos))

		if err != nil {
			log.Println("[❌] Error listing repositories for organization: " + organization)
			continue
		}

		enabledRepos := []string{}
		for _, repo := range repos {
			//print repo

			var repository github.Repository
			repository, err = github.GetRepository(*repo.Name, organization, token, url)

			if err != nil {
				log.Println("[❌] Error getting repository: ", err)
				return err
			}

			if repository.SecurityAndAnalysis == nil || repository.SecurityAndAnalysis.SecretScanning == nil || *repository.SecurityAndAnalysis.SecretScanning.Status != "enabled" {
				continue
			}

			log.Println("Secret scanning enabled for repository: " + *repo.Name)
			enabledRepos = append(enabledRepos, *repo.Name)
		}

		mapOrgAndRepos[organization] = enabledRepos
		log.Println("[✅] Done")
	}

	//convert mapOrgAndRepos to string
	result := ""
	count := 0
	for org, repos := range mapOrgAndRepos {
		result += org + "\n"
		for _, repo := range repos {
			result += "  " + repo + "\n"
			count++
		}
	}

	log.Printf("[✅] Done. %d repositories with secret scanning enabled out of %d", count, totalRepos)

	// save to file using ioutil.WriteFile
	ioutil.WriteFile("report.txt", []byte(result), 0644)

	return nil
}
