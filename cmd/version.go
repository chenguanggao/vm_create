/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var shortPrint bool

// versionCmd represents the version command
func VersionCmd() *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version of the vm tool",
		Long:  `vm tool create .`,
		RunE: func(cmd *cobra.Command, args []string) error {
			marshalled, err := json.Marshal(Get())
			if err != nil {
				return err
			}
			if shortPrint {
				fmt.Println(Get().String())
			} else {
				fmt.Println(string(marshalled))
			}
			return nil
		},
	}
	return versionCmd
}

// component-base/version/base.go
var (
	gitMajor string // major version, always numeric
	gitMinor string // minor version, numeric possibly followed by "+"

	gitVersion = "unknown"
	gitCommit  = "" // sha1 from git, output of $(git rev-parse HEAD)

	buildDate = "1970-01-01T00:00:00Z" // build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
)

type Info struct {
	Major      string `json:"major,omitempty"`
	Minor      string `json:"minor,omitempty"`
	GitVersion string `json:"gitVersion"`
	GitCommit  string `json:"gitCommit,omitempty"`
	BuildDate  string `json:"buildDate"`
	GoVersion  string `json:"goVersion"`
	Compiler   string `json:"compiler"`
	Platform   string `json:"platform"`
}

// String returns info as a human-friendly version string.
func (info Info) String() string {
	return info.GitVersion
}

func Get() Info {
	// These variables typically come from -ldflags settings and in
	// their absence fallback to the settings in ./base.go
	return Info{
		Major:      gitMajor,
		Minor:      gitMinor,
		GitVersion: gitVersion,
		GitCommit:  gitCommit,
		BuildDate:  buildDate,
		GoVersion:  runtime.Version(),
		Compiler:   runtime.Compiler,
		Platform:   fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// GetSingleVersion returns single version of sealer
func GetSingleVersion() string {
	return gitVersion
}
