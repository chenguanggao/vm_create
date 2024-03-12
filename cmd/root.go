/*
Copyright © 2024 gcg <hpu_gcg@163.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

type rootOpts struct {
	cfgFile string
	// configFile string
}

var (
	rootOpt      rootOpts
	log          = logrus.New()
	Verbose      bool
	MarkdownDocs bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vm",
	Short: "Create a kvm VM",
	Long:  `Run the virtual machine creation process according to excel.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(RunCmd(), DeleteCmd(), VersionCmd(), ConvertCmd())

	rootCmd.PersistentFlags().StringVar(&rootOpt.cfgFile, "excel", "./parameter.xlsx", "excel file of vm tool")
	// rootCmd.PersistentFlags().StringVar(&rootOpt.configFile, "config", "./Sheet1Config.yaml", "yaml file of vm tool")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolVarP(&MarkdownDocs, "md-docs", "m", false, "gen Markdown docs")
}

func initConfig() {
	if rootOpt.cfgFile == "" {
		// Find home directory.
		rootOpt.cfgFile = filepath.Join("parameter.xlsx")
	}
	// if rootOpt.configFile != "" {
	// 	viper.SetConfigFile(rootOpt.configFile)
	// } else {
	// 	viper.SetConfigFile("./Sheet1Config.yaml")
	// }
	// viper.ReadInConfig()
	if file, err := os.OpenFile("./create.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err != nil {
		log.Info("failed to log to file,use default stderr")
	} else {
		log.SetLevel(logrus.DebugLevel)
		log.SetReportCaller(true)
		// log.SetFormatter(&logrus.JSONFormatter{})
		multiWriter := io.MultiWriter(os.Stdout, file)
		log.Out = multiWriter

	}

}

func GenDocs() {
	if MarkdownDocs {
		if err := doc.GenMarkdownTree(rootCmd, "./docs/md"); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
