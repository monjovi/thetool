package cmd

import (
	"fmt"

	"github.com/solo-io/thetool/pkg/downloader"
	"github.com/solo-io/thetool/pkg/feature"
	"github.com/spf13/cobra"
)

func AddCmd() *cobra.Command {
	var featureName string
	var featureRepository string
	var featureHash string
	var verbose bool

	cmd := &cobra.Command{
		Use:   "add",
		Short: "add a feature",
		RunE: func(c *cobra.Command, args []string) error {
			return runAdd(featureName, featureRepository, featureHash, verbose)
		},
	}

	cmd.PersistentFlags().StringVarP(&featureName, "name", "n", "", "Feature name")
	cmd.PersistentFlags().StringVarP(&featureRepository, "repository", "r", "", "Repository URL")
	cmd.PersistentFlags().StringVarP(&featureHash, "commit", "c", "", "Commit hash")
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")

	return cmd
}

func runAdd(name, repo, hash string, verbose bool) error {
	if name == "" {
		return fmt.Errorf("feature name can't be empty")
	}
	if repo == "" {
		return fmt.Errorf("feature repository URL can't be empty")
	}
	if hash == "" {
		return fmt.Errorf("feature commit hash can't be empty")
	}

	f := feature.Feature{
		Name:       name,
		Version:    hash,
		Repository: repo,
		Enabled:    true,
	}
	existing, err := feature.LoadFromFile(dataFile)
	if err != nil {
		fmt.Printf("Unable to load feature list: %q\n", err)
		return nil
	}

	// check it isn't already existing feature
	for _, e := range existing {
		if e.Name == name {
			fmt.Printf("Feature %s already exists\n", name)
			return nil
		}
	}
	// let's get the external dependency
	err = downloader.Download(f, repositoryDir, verbose)
	if err != nil {
		fmt.Printf("Unable to download repository %s: %q\n", f.Repository, err)
		return nil
	}

	err = feature.SaveToFile(append(existing, f), "features.json")
	if err != nil {
		fmt.Printf("Unable to save feature %s: %q\n", f.Name, err)
	}
	return nil
}
