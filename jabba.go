package main

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/shyiko/jabba/command"
	"github.com/shyiko/jabba/semver"
	log "github.com/Sirupsen/logrus"
	"sort"
	"io/ioutil"
	"strings"
	"bytes"
)

var version string
var rootCmd *cobra.Command

func init() {
	log.SetFormatter(&SimpleFormatter{})
	// todo: make it configurable through the command line
	log.SetLevel(log.InfoLevel)
}

type SimpleFormatter struct {}

func (f *SimpleFormatter) Format(entry *log.Entry) ([]byte, error) {
	b := &bytes.Buffer{}
	fmt.Fprintf(b, "%s ", entry.Message)
	for k, v := range entry.Data {
		fmt.Fprintf(b, "%s=%+v ", k, v)
	}
	b.WriteByte('\n')
	return b.Bytes(), nil
}

func main() {
	rootCmd = &cobra.Command{
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
				if len(args) == 0 {
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
				if len(args) == 0 {
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
			Use:   "link [name] [path]",
			Short: "Resolve or update a link",
			RunE: func(cmd *cobra.Command, args []string) error {
				if len(args) == 0 {
					return pflag.ErrHelp
				}
				if len(args) == 1 {
					if value := command.GetLink(args[0]); value != "" {
						fmt.Println(value)
					}
				} else
				if err := command.Link(args[0], args[1]); err != nil {
					log.Fatal(err)
				}
				return nil
			},
			Example: "  jabba link system@1.8.20 /Library/Java/JavaVirtualMachines/jdk1.8.0_20.jdk" +
			"  jabba link system@1.8.20 # show link target",
		},
		&cobra.Command{
			Use:   "unlink [name]",
			Short: "Delete a link",
			RunE: func(cmd *cobra.Command, args []string) error {
				if len(args) == 0 {
					return pflag.ErrHelp
				}
				if err := command.Link(args[0], ""); err != nil {
					log.Fatal(err)
				}
				return nil
			},
		},
		&cobra.Command{
			Use:   "use [version to use]",
			Short: "Modify PATH & JAVA_HOME to use specific JDK",
			RunE: func(cmd *cobra.Command, args []string) error {
				if len(args) == 0 {
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
				if ver != "" {
					fmt.Println(ver)
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
				var vs = make([]*semver.Version, len(releaseMap))
				var i = 0
				for k := range releaseMap {
					vs[i] = k
					i++
				}
				sort.Sort(sort.Reverse(semver.VersionSlice(vs)))
				for _, v := range vs {
					fmt.Println(v)
				}
				return nil
			},
		},
		&cobra.Command{
			Use:   "deactivate",
			Short: "Undo effects of `jabba` on current shell",
			RunE: func(cmd *cobra.Command, args []string) error {
				out, err := command.Deactivate()
				if err != nil {
					log.Fatal(err)
				}
				printForShellToEval(out)
				return nil
			},
		},
		&cobra.Command{
			Use:   "alias [name] [version]",
			Short: "Resolve or update an alias",
			RunE: func(cmd *cobra.Command, args []string) error {
				if len(args) == 0 {
					return pflag.ErrHelp
				}
				if len(args) == 1 {
					if value := command.GetAlias(args[0]); value != "" {
						fmt.Println(value)
					}
				} else
				if err := command.SetAlias(args[0], args[1]); err != nil {
					log.Fatal(err)
				}
				return nil
			},
			Example: "  jabba alias default 1.8\n" +
			"  jabba alias default # show value bound to an alias",
		},
		&cobra.Command{
			Use:   "unalias [name]",
			Short: "Delete an alias",
			RunE: func(cmd *cobra.Command, args []string) error {
				if len(args) == 0 {
					return pflag.ErrHelp
				}
				if err := command.SetAlias(args[0], ""); err != nil {
					log.Fatal(err)
				}
				return nil
			},
		},
		&cobra.Command{
			Use:   "which [version]",
			Short: "Display path to installed JDK",
			RunE: func(cmd *cobra.Command, args []string) error {
				if len(args) == 0 {
					return pflag.ErrHelp
				}
				dir, _ := command.Which(args[0])
				if dir != "" {
					fmt.Println(dir)
				}
				return nil
			},
		},
	)
	rootCmd.Flags().Bool("version", false, "version of jabba")
	rootCmd.PersistentFlags().String("fd3", "", "")
	rootCmd.PersistentFlags().MarkHidden("fd3")
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1);
	}
}

func use(ver string) error {
	out, err := command.Use(ver)
	if err != nil {
		log.Fatal(err)
	}
	printForShellToEval(out)
	return nil
}

func printForShellToEval(out []string) {
	fd3, _ := rootCmd.Flags().GetString("fd3")
	if fd3 != "" {
		ioutil.WriteFile(fd3, []byte(strings.Join(out, "\n")), 0666)
	} else {
		fd3 := os.NewFile(3, "fd3")
		for _, line := range out {
			fmt.Fprintln(fd3, line)
		}
	}
}
