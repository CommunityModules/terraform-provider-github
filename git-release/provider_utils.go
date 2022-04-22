package git_release

import (
	"log"
	"os"
	"testing"
)

var testOrganization string = testOrganizationFunc()
var testToken string = os.Getenv("GITHUB_TOKEN")
var testBaseURLGHES string = os.Getenv("GHES_BASE_URL")

func skipUnlessMode(t *testing.T, providerMode string) {
	switch providerMode {
	case anonymous:
		if os.Getenv("GITHUB_BASE_URL") != "" &&
			os.Getenv("GITHUB_BASE_URL") != "https://api.github.com/" {
			t.Log("anonymous mode not supported for GHES deployments")
			break
		}

		if os.Getenv("GITHUB_TOKEN") == "" {
			return
		} else {
			t.Log("GITHUB_TOKEN environment variable should be empty")
		}
	case individual:
		if os.Getenv("GITHUB_TOKEN") != "" && os.Getenv("GITHUB_OWNER") != "" {
			return
		} else {
			t.Log("GITHUB_TOKEN and GITHUB_OWNER environment variables should be set")
		}
	case organization:
		if os.Getenv("GITHUB_TOKEN") != "" && os.Getenv("GITHUB_ORGANIZATION") != "" {
			return
		} else {
			t.Log("GITHUB_TOKEN and GITHUB_ORGANIZATION environment variables should be set")
		}
	}

	t.Skipf("Skipping %s which requires %s mode", t.Name(), providerMode)
}

func OwnerOrOrgEnvDefaultFunc() (interface{}, error) {
	if organization := os.Getenv("GITHUB_ORGANIZATION"); organization != "" {
		log.Printf("[INFO] Selecting owner %s from GITHUB_ORGANIZATION environment variable", organization)
		return organization, nil
	}
	owner := os.Getenv("GITHUB_OWNER")
	log.Printf("[INFO] Selecting owner %s from GITHUB_OWNER environment variable", owner)
	return owner, nil
}

func testOrganizationFunc() string {
	organization := os.Getenv("GITHUB_ORGANIZATION")
	if organization == "" {
		organization = os.Getenv("GITHUB_TEST_ORGANIZATION")
	}
	return organization
}

func testOwnerFunc() string {
	owner := os.Getenv("GITHUB_OWNER")
	if owner == "" {
		owner = os.Getenv("GITHUB_TEST_OWNER")
	}
	return owner
}

const anonymous = "anonymous"
const individual = "individual"
const organization = "organization"
