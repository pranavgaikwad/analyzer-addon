package main

import (
	"os"
	"path"
	"time"

	"github.com/konveyor/tackle2-addon/repository"
	"github.com/konveyor/tackle2-addon/ssh"
	hub "github.com/konveyor/tackle2-hub/addon"
	"github.com/konveyor/tackle2-hub/api"
	"github.com/konveyor/tackle2-hub/nas"
	"k8s.io/apimachinery/pkg/util/errors"
)

const (
	ENV_WORKDIR = "WORKDIR"
)

var (
	addon                = hub.Addon
	WorkDir              = ""
	ConfigDir            = "config"
	ReportDir            = "report"
	RuleDir              = "rules"
	SourceDir            = "source"
	BinDir               = "bin"
	ProviderSettingsPath = "provider-settings.json"
	OutputPath           = "output.yaml"
)

func init() {
	WorkDir = os.Getenv(ENV_WORKDIR)
	ConfigDir = path.Join(WorkDir, ConfigDir)
	ReportDir = path.Join(WorkDir, ReportDir)
	RuleDir = path.Join(WorkDir, RuleDir)
	SourceDir = path.Join(WorkDir, SourceDir)
	BinDir = path.Join(WorkDir, BinDir)
	ProviderSettingsPath = path.Join(ConfigDir, "provider-settings.json")
	OutputPath = path.Join(ReportDir, "output.yaml")
}

type Data struct {
	Rulesets []api.Ref `json:"rulesets"`
	Output   string    `json:"output"`
}

func main() {
	addon.Run(func() (err error) {
		// load data from task
		data := &Data{}
		err = addon.DataWith(data)
		if err != nil {
			err = &hub.SoftError{Reason: err.Error()}
			return
		}
		// create directories
		for _, dir := range []string{
			ConfigDir, ReportDir, RuleDir, SourceDir, BinDir} {
			err = nas.MkDir(dir, 0755)
			if err != nil {
				return
			}
		}

		analyzer := Analyzer{}
		analyzer.Data = data
		// Fetch application.
		addon.Activity("Fetching application.")
		application, err := addon.Task.Application()
		if err == nil {
			analyzer.application = application
		} else {
			return
		}
		// Delete old report
		// mark := time.Now()
		// bucket := addon.Application.Bucket(application.ID)
		// err = bucket.Delete(OutputPath)
		// if err != nil {
		// 	return
		// }
		// addon.Activity(
		// 	"[BUCKET] Report deleted:%s duration:%v.",
		// 	OutputPath,
		// 	time.Since(mark))
		// Setup SSH
		agent := ssh.Agent{}
		err = agent.Start()
		if err != nil {
			return
		}
		// Download source code
		err = ensureSourceCode(application.Repository, application.Identities)
		if err != nil {
			return
		}
		// Run analyzer
		mark := time.Now()
		err = analyzer.Run()
		if err != nil {
			return
		}
		addon.Activity(
			"[BUCKET] Report generation:%s duration:%v.",
			OutputPath,
			time.Since(mark))
		return
	})
}

func ensureSourceCode(repo *api.Repository, ids []api.Ref) (err error) {
	r, err := repository.New(SourceDir, repo, ids)
	if err != nil {
		return
	}
	err = r.Fetch()
	if err == nil {
		addon.Increment()
	} else {
		return
	}
	return
}

func cleanup() error {
	errs := []error{}
	for _, dir := range []string{SourceDir} {
		err := nas.RmDir(dir)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errors.NewAggregate(errs)
}
