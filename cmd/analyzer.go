package main

import (
	"encoding/json"
	"fmt"
	"os"
	pathlib "path"
	"strconv"

	"github.com/konveyor/analyzer-lsp/provider/builtin"
	"github.com/konveyor/analyzer-lsp/provider/java"
	provider "github.com/konveyor/analyzer-lsp/provider/lib"
	"github.com/konveyor/tackle2-addon/command"
	"github.com/konveyor/tackle2-addon/repository"
	"github.com/konveyor/tackle2-hub/api"
	"github.com/konveyor/tackle2-hub/nas"
)

type Analyzer struct {
	application *api.Application
	*Data
}

const (
	Bin      = "/home/pranav/Projects/analyzer-lsp/analyzer-lsp"
	TagsFile = "tags-file.yaml"
)

func (a *Analyzer) Run() error {
	err := a.ensureProviderSettings()
	if err != nil {
		return err
	}

	cmd := command.Command{Path: Bin}
	cmd.Options, err = a.options()
	if err != nil {
		return err
	}

	err = cmd.Run()
	if err != nil {
		fmt.Printf("%v", cmd.Output)
		return err
	}
	return nil
}

func (a *Analyzer) options() (opts command.Options, err error) {
	err = a.AddOptions(&opts)
	if err != nil {
		return
	}
	opts.Add("--provider-settings", ProviderSettingsPath)
	opts.Add("--output-file", OutputPath)
	return
}

func (d *Data) AddOptions(opts *command.Options) error {
	for _, ruleset := range d.Rulesets {
		var bundle *api.RuleBundle
		bundle, err := addon.RuleBundle.Get(ruleset.ID)
		if err != nil {
			return err
		}
		path, err := fetchBundleRepo(bundle)
		if err != nil {
			return err
		}
		opts.Add("--rules", path)
	}
	return nil
}

func fetchBundleRepo(bundle *api.RuleBundle) (string, error) {
	if bundle.Repository == nil {
		return "", nil
	}
	rootDir := pathlib.Join(
		RuleDir,
		"bundles",
		strconv.Itoa(int(bundle.ID)),
		"repository")
	err := nas.MkDir(rootDir, 0755)
	if err != nil {
		return "", err
	}
	identities := []api.Ref{}
	if bundle.Identity != nil {
		identities = []api.Ref{*bundle.Identity}
	}
	rp, err := repository.New(
		rootDir,
		bundle.Repository,
		identities)
	if err != nil {
		return "", err
	}
	err = rp.Fetch()
	if err != nil {
		return "", err
	}
	return pathlib.Join(rootDir, bundle.Repository.Path), nil
}

func (a *Analyzer) ensureProviderSettings() error {
	configs := []provider.Config{
		{
			Name:     "builtin",
			Location: SourceDir,
			ProviderSpecificConfig: map[string]string{
				builtin.TAGS_FILE_INIT_OPTION: pathlib.Join(ConfigDir, TagsFile),
			},
		},
		{
			Name:     "java",
			Location: SourceDir,
			BinaryLocation: pathlib.Join(
				BinDir,
				"jdtls",
			),
			ProviderSpecificConfig: map[string]string{
				java.BUNDLES_INIT_OPTION: pathlib.Join(
					BinDir,
					"java-analyzer-bundle.core-1.0.0-SNAPSHOT.jar",
				),
				java.WORKSPACE_INIT_OPTION: BinDir,
			},
		},
	}
	configContent, err := json.Marshal(configs)
	if err != nil {
		return err
	}
	err = os.WriteFile(
		ProviderSettingsPath,
		configContent,
		0755,
	)
	if err != nil {
		return err
	}
	return nil
}
