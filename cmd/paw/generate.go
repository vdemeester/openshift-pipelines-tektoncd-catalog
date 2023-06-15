package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	// "github.com/Masterminds/semver"
	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	"github.com/openshift-pipelines/tektoncd-catalog/cmd/paw/config"
	"github.com/spf13/cobra"
)

type generateOptions struct {
	externalTaskConfig string
}

func Generate() *cobra.Command {
	opts := &generateOptions{}
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a folder-based catalog from a folder",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("Require 2 argument, got : %d", len(args))
			}
			source := args[0]
			dest := args[1]
			return generateDirectoryCatalog(source, dest, opts.externalTaskConfig)
		},
	}
	cmd.Flags().StringVarP(&opts.externalTaskConfig, "external-tasks", "e", "external-tasks.yaml", "Configuration file for external tasks")
	return cmd
}

type resources struct {
	resources map[string]resourceByVersion
}

type resourceByVersion struct {
	versions map[string]resource
}

type resource struct {
	name string
	url  string
}

func generateDirectoryCatalog(source, dest, externalTaskConfig string) error {
	c, err := config.Load(filepath.Join(source, externalTaskConfig))
	if err != nil {
		return err
	}
	ghclient, err := gh.RESTClient(nil)
	if err != nil {
		return err
	}
	// TODO: prepare the "dest" workspace
	//       (fetch the repository's `p` branch, …)
	// TODO: fetch release assets
	//       - fetch yamls
	//       - fetch tests (yamls with kttl)
	//       - fetch bundles, sbom, …
	// TODO: extract in destination folder
	// TODO: copy source/ to dest/ as well
	//       (warn if there is conflicts)
	// TODO: create a PR
	resources, err := mapReleasedTasks(ghclient, c)
	if err != nil {
		return err
	}
	for r, rbv := range resources.resources {
		if err := os.MkdirAll(filepath.Join(dest, r), 0755); err != nil {
			return err
		}
		for v, res := range rbv.versions {
			if err := os.MkdirAll(filepath.Join(dest, r, v), 0755); err != nil {
				return nil
			}
			fmt.Println(v, res)
		}
	}
	return nil
}

func mapReleasedTasks(ghclient api.RESTClient, c config.Config) (resources, error) {
	resources := resources{
		resources: map[string]resourceByVersion{},
	}
	for _, t := range c.Tasks {
		if !strings.HasPrefix(t.Repository, "https://github.com/") {
			return resources, fmt.Errorf("Non-github repository not supported, provided: %s", t.Repository)
		}
		repo := strings.TrimPrefix(t.Repository, "https://github.com/")
		versions, err := fetchVersions(ghclient, repo)
		if err != nil {
			return resources, err
		}
		if len(versions) == 0 {
			fmt.Fprintf(os.Stderr, "%s has no release, ignoring\n", repo)
			continue
		}
		for _, v := range versions {
			// FIXME: update this..
			// For now, we start with just fetching the application/x-yaml at the right place
			// Long term, this is going to be way more involved
			// - Find yamls
			// - Find tests (required)
			// - Find bundle(s)
			// - Find attestation / sbom / signatures
			for _, a := range v.Assets {
				if a.ContentType != "application/x-yaml" {
					fmt.Fprintf(os.Stderr, "%s's asset %s ignored, not a yaml file\n", repo, a.Name)
					continue
				}
				// FIXME: basic assumption here, the task name is the file name (minus the extension)
				name := strings.TrimSuffix(a.Name, filepath.Ext(a.Name))
				if _, ok := resources.resources[name]; !ok {
					resources.resources[name] = resourceByVersion{
						versions: map[string]resource{},
					}
				}
				resources.resources[name].versions[v.TagName] = resource{
					name: name,
					url:  a.DownloadURL,
				}
			}
		}
	}
	return resources, nil
}

type Version struct {
	Name       string
	TagName    string `json:"tag_name"`
	Id         int
	Draft      bool
	PreRelease bool
	Assets     []Assets
}

type Assets struct {
	Id          int
	URL         string `json:"url"`
	Name        string
	Label       string
	ContentType string `json:"content_type"`
	State       string
	DownloadURL string `json:"browser_download_url"`
}

func fetchVersions(client api.RESTClient, github string) ([]Version, error) {
	versions := []Version{}
	err := client.Get(fmt.Sprintf("repos/%s/releases", github), &versions)
	if err != nil {
		return nil, err
	}
	return versions, nil
}
