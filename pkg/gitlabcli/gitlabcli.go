package gitlabcli

import (
	"github.com/lokalise/go-lokalise-api/v3"
)

type Cli struct {
	client *lokalise.Api
}

func New(client *lokalise.Api) *Cli {
	return &Cli{
		client: client,
	}
}

// PullTranslations makes an API call to Lokalise to generate an archive with translations.
// It returns a link to the newly generated archive.
// If triggerGitlab is set to true, it will also force Lokalise to trigger gitlab webhook,
// as a result there will be a new Merge Request in the corresponding project in Gitlab.
func (c *Cli) PullTranslations(projectID string, triggerGitlab bool) (bundleURL string, err error) {
	f := c.client.Files()
	directoryPrefix := ""

	var triggers []string
	if triggerGitlab {
		triggers = append(triggers, "gitlab")
	}

	r, err := f.Download(projectID, lokalise.FileDownload{
		Format:            "yml",
		DirectoryPrefix:   &directoryPrefix,
		BundleStructure:   "translations/messages+intl-icu.%LANG_ISO%.%FORMAT%",
		AllPlatforms:      true,
		Triggers:          triggers,
		PluralFormat:      "icu",
		PlaceholderFormat: "icu",
		YAMLIncludeRoot:   false,
	})
	if err != nil {
		return "", err
	}

	return r.BundleURL, nil
}
