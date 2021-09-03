package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jessevdk/go-flags"
	"io/ioutil"
	"net/http"
)

type CleanCommand struct {
	TokenType string `long:"token-type" description:"Gitlab Token Type" choice:"job" choice:"private" default:"private" env:"GITLAB_TOKEN_TYPE"`
	Token  string `long:"token" description:"API token to access Gitlab ($CI_JOB_TOKEN)" required:"true" env:"GITLAB_TOKEN"`
	ApiV4URL string `long:"api-v4-url" description:"Gitlab API v4 URL ($CI_API_V4_URL)" required:"true" env:"GITLAB_API_V4_URL"`
	ProjectID      string `long:"project-id" description:"Gitlab Project ID ($CI_PROJECT_ID)" required:"true" env:"GITLAB_PROJECT_ID"`
	RepositoryName string `long:"repository-name" description:"Gitlab Repository Name ($CI_COMMIT_REF_SLUG)" required:"true" env:"GITLAB_REPOSITORY_NAME"`
}

func RegisterCleanCommand(parser *flags.Parser) *CleanCommand {
	cmd := &CleanCommand{}
	_, err := parser.AddCommand("clean", "makes an API call to Gitlab to delete a docker repository", "", cmd)
	if err != nil {
		panic(err)
	}
	return cmd
}

type RepositoryData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (cmd *CleanCommand) Execute(_ []string) error {
	var found *RepositoryData

	i := 1
	for ;;i++ {
		u := fmt.Sprintf("%s/projects/%s/registry/repositories?page=%d&per_page=100", cmd.ApiV4URL, cmd.ProjectID, i)
		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			return err
		}
		if cmd.TokenType == "private" {
			req.Header.Add("PRIVATE-TOKEN", cmd.Token)
		} else {
			req.Header.Add("JOB-TOKEN", cmd.Token)
		}

		data, err := (func() ([]byte, error) {
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				return nil, errors.New("status code is not 200")
			}

			// fmt.Println(resp.Header.Values("Link"))
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}

			return data, nil
		})()
		if err != nil {
			return err
		}

		// fmt.Println("Existing repositories: " + string(data))

		var repos []RepositoryData
		err = json.Unmarshal(data, &repos)
		if err != nil {
			return err
		}

		if len(repos) == 0 {
			break
		}

		for _, repo := range repos {
			if repo.Name == cmd.RepositoryName {
				d, err := json.Marshal(repo)
				if err != nil {
					return err
				}
				fmt.Println(string(d))
				found = &repo
				break
			}
		}
	}

	if found == nil {
		return errors.New("repository not found")
	}

	u := fmt.Sprintf("%s/projects/%s/registry/repositories/%d", cmd.ApiV4URL, cmd.ProjectID, found.ID)
	req, err := http.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}
	if cmd.TokenType == "private" {
		req.Header.Add("PRIVATE-TOKEN", cmd.Token)
	} else {
		req.Header.Add("JOB-TOKEN", cmd.Token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 202 {
		return fmt.Errorf("unexpected status code (expected 202): %d", resp.StatusCode)
	}

	return nil
}
