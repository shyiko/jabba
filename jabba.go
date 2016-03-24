package main

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/shyiko/jabba/command"
	"github.com/shyiko/jabba/semver"
	log "github.com/Sirupsen/logrus"
)

var version string

func init() {
	// todo: make it configurable through the command line
	log.SetLevel(log.InfoLevel)
}

func main() {
	var rootCmd = &cobra.Command{
		Use: "jabba",
		Long: "Java Version Manager (https://github.com/shyiko/jabba).",
		RunE: func(cmd *cobra.Command, args []string) error {
			if showVersion, _ := cmd.Flags().GetBool("version"); !showVersion {
				return pflag.ErrHelp
			}
			fmt.Println(version)
			return nil
		},
	}
	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "install [version to install]",
			Short: "Download and install JDK",
			RunE: func(cmd *cobra.Command, args []string) error {
				if (len(args) == 0) {
					return pflag.ErrHelp
				}
				ver, err := command.Install(args[0])
				if err != nil {
					log.Fatal(err)
				}
				return use(ver)
			},
			Example: "  jabba install 1.8\n" +
			"  jabba install ~1.8.73 # same as \">=1.8.73 <1.9.0\"\n" +
			"  jabba install 1.8.73=dmg+http://.../jdk-9-ea+110_osx-x64_bin.dmg",
		},
		&cobra.Command{
			Use:   "uninstall [version to uninstall]",
			Short: "Uninstall JDK",
			RunE: func(cmd *cobra.Command, args []string) error {
				if (len(args) == 0) {
					return pflag.ErrHelp
				}
				err := command.Uninstall(args[0])
				if err != nil {
					log.Fatal(err)
				}
				return nil
			},
			Example: "  jabba uninstall 1.8",
		},
		&cobra.Command{
			Use:   "use [version to use]",
			Short: "Modify PATH & JAVA_HOME to use specific JDK",
			RunE: func(cmd *cobra.Command, args []string) error {
				if (len(args) == 0) {
					return pflag.ErrHelp
				}
				return use(args[0])
			},
			Example: "  jabba use 1.8",
		},
		&cobra.Command{
			Use:   "current",
			Short: "Display currently 'use'ed version",
			Run: func(cmd *cobra.Command, args []string) {
				ver := command.Current()
				if (ver != "") {
					println(ver)
				}
			},
		},
		&cobra.Command{
			Use:   "ls",
			Short: "List installed versions",
			RunE: func(cmd *cobra.Command, args []string) error {
				releases, err := command.Ls()
				if err != nil {
					log.Fatal(err)
				}
				for _, v := range releases {
					fmt.Println(v)
				}
				return nil
			},
		},
		&cobra.Command{
			Use:   "ls-remote",
			Short: "List remote versions available for install",
			RunE: func(cmd *cobra.Command, args []string) error {
				releaseMap, err := command.LsRemote()
				if err != nil {
					log.Fatal(err)
				}
				var vs = make([]string, len(releaseMap))
				var i = 0
				for k := range releaseMap {
					vs[i] = k
					i++
				}
				for _, v := range semver.Sort(vs) {
					fmt.Println(v)
				}
				return nil
			},
		},
	)
	rootCmd.Flags().Bool("version", false, "version of jabba")
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1);
	}
}

func use(ver string) error {
	out, err := command.Use(ver)
	if err != nil {
		log.Fatal(err)
	}
	// writing "export ..." to 3rd fd for shell to eval
	fd3 := os.NewFile(3, "fd3")
	for _, line := range out {
		fmt.Fprintln(fd3, line)
	}
	return nil
}
