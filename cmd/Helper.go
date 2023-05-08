package cmd

import (
	"log"

	"github.com/gateixeira/ghas-org-enablement/cmd/github"
)

func ManageGhasFeatures(enterprise, organization, token, url string, activate bool) error {
	if organization != "" {
		log.Println("[🔄] Changing GHAS settings for organization: " + organization)
		err := github.ChangeGHASOrgSettings(organization, activate, token, url)

		if err != nil {
			log.Println("[❌] Error changing GHAS settings for organization: " + organization)
			return err
		}

		log.Println("[✅] Done")
		return nil
	}

	log.Println("[🔄] Fetching organizations from enterprise...")
	organizations, err := github.GetOrganizationsInEnterprise(enterprise, token, url)
	log.Println("[✅] Done")

	if err != nil {
		log.Println("[❌] Error fetching organizations from enterprise")
		return err
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
