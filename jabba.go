package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"sort"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/Jabba-Team/jabba/command"
	"github.com/Jabba-Team/jabba/semver"
	log "github.com/Sirupsen/logrus"
	rootcerts "github.com/hashicorp/go-rootcerts"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var version string
var rootCmd *cobra.Command

func init() {
	log.SetFormatter(&simpleFormatter{})
	// todo: make it configurable through the command line
	log.SetLevel(log.InfoLevel)

	tlsConfig := &tls.Config{}
	err := rootcerts.ConfigureTLS(tlsConfig, &rootcerts.Config{
		CAFile: os.Getenv("JABBA_CAFILE"),
		CAPath: os.Getenv("JABBA_CAPATH"),
	})
	if err != nil {
		log.Fatal(err)
	}
	defTransport := http.DefaultTransport.(*http.Transport)
	defTransport.TLSClientConfig = tlsConfig
}

type simpleFormatter struct{}

func (f *simpleFormatter) Format(entry *log.Entry) ([]byte, error) {
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
		Use:  "jabba",
		Long: "Java Version Manager (https://github.com/Jabba-Team/jabba).",
		RunE: func(cmd *cobra.Command, args []string) error {
			if showVersion, _ := cmd.Flags().GetBool("version"); !showVersion {
				return pflag.ErrHelp
			}
			fmt.Println(version)
			return nil
		},
	}
	var whichHome bool
	whichCmd := &cobra.Command{
		Use:   "which [version]",
		Short: "Display path to installed JDK",
		RunE: func(cmd *cobra.Command, args []string) error {
			var ver string
			if len(args) == 0 {
				ver = rc().JDK
				if ver == "" {
					return pflag.ErrHelp
				}
			} else {
				ver = args[0]
			}
			dir, _ := command.Which(ver, whichHome)
			if dir != "" {
				fmt.Println(dir)
			}
			return nil
		},
	}
	whichCmd.Flags().BoolVar(&whichHome, "home", false,
		"Account for platform differences so that value could be used as JAVA_HOME (e.g. append \"/Contents/Home\" on macOS)")
	var customInstallDestination string
	installCmd := &cobra.Command{
		Use:   "install [version to install]",
		Short: "Download and install JDK",
		RunE: func(cmd *cobra.Command, args []string) error {
			var ver string
			if len(args) == 0 {
				ver = rc().JDK
				if ver == "" {
					return pflag.ErrHelp
				}
			} else {
				ver = args[0]
			}
			ver, err := command.Install(ver, customInstallDestination)
			if err != nil {
				log.Fatal(err)
			}
			if customInstallDestination == "" {
				if err := command.LinkLatest(); err != nil {
					log.Fatal(err)
				}
				return use(ver)
			} else {
				return nil
			}
		},
		Example: "  jabba install 1.8\n" +
			"  jabba install ~1.8.73 # same as \">=1.8.73 <1.9.0\"\n" +
			"  jabba install 1.8.73=dmg+http://.../jdk-9-ea+110_osx-x64_bin.dmg",
	}
	installCmd.Flags().StringVarP(&customInstallDestination, "output", "o", "",
		"Custom destination (any JDK outside of $JABBA_HOME/jdk is considered to be unmanaged, i.e. not available to jabba ls, use, etc. (unless `jabba link`ed))")
	var trimTo string
	lsCmd := &cobra.Command{
		Use:   "ls",
		Short: "List installed versions",
		RunE: func(cmd *cobra.Command, args []string) error {
			var r *semver.Range
			if len(args) > 0 {
				var err error
				r, err = semver.ParseRange(args[0])
				if err != nil {
					log.Fatal(err)
				}
			}
			vs, err := command.Ls()
			current := command.Current()
			if err != nil {
				log.Fatal(err)
			}
			if trimTo != "" {
				vs = semver.VersionSlice(vs).TrimTo(parseTrimTo(trimTo))
			}
			for _, v := range vs {
				if r != nil && !r.Contains(v) {
					continue
				}
				if current == v.String() {
					fmt.Print("\033[32m")
					fmt.Println(v.String() + "*")
					fmt.Print("\033[0m")
				} else {
					fmt.Println(v)
				}
			}
			return nil
		},
	}
	lsRemoteCmd := &cobra.Command{
		Use:   "ls-remote",
		Short: "List remote versions available for install",
		RunE: func(cmd *cobra.Command, args []string) error {
			var r *semver.Range
			if len(args) > 0 {
				var err error
				r, err = semver.ParseRange(args[0])
				if err != nil {
					log.Fatal(err)
				}
			}
			os, _ := cmd.Flags().GetString("os")
			arch, _ := cmd.Flags().GetString("arch")
			releaseMap, err := command.LsRemote(os, arch)
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
			if trimTo != "" {
				vs = semver.VersionSlice(vs).TrimTo(parseTrimTo(trimTo))
			}
			for _, v := range vs {
				if r != nil && !r.Contains(v) {
					continue
				}
				fmt.Println(v)
			}
			return nil
		},
	}
	lsRemoteCmd.Flags().String("os", runtime.GOOS, "Operating System (darwin, linux, windows)")
	lsRemoteCmd.Flags().String("arch", runtime.GOARCH, "Architecture (amd64, 386)")
	for _, cmd := range []*cobra.Command{lsCmd, lsRemoteCmd} {
		cmd.Flags().StringVar(&trimTo, "latest", "",
			"Part of the version to trim to (\"major\", \"minor\" or \"patch\")")
	}
	rootCmd.AddCommand(
		installCmd,
		&cobra.Command{
			Use:   "uninstall [version to uninstall]",
			Short: "Uninstall JDK",
			RunE: func(cmd *cobra.Command, args []string) error {
				if len(args) == 0 {
					return pflag.ErrHelp
				}
				if strings.HasPrefix(args[0], "system@") {
					log.Fatal("Link to system JDK can only be removed with 'unlink'" +
						" (e.g. 'jabba unlink " + args[0] + "')")
				}
				err := command.Uninstall(args[0])
				if err != nil {
					log.Fatal(err)
				}
				if err := command.LinkLatest(); err != nil {
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
					if err := command.LinkLatest(); err != nil {
						log.Fatal(err)
					}
					return nil
				}
				if len(args) == 1 {
					if value := command.GetLink(args[0]); value != "" {
						fmt.Println(value)
					}
				} else if err := command.Link(args[0], args[1]); err != nil {
					log.Fatal(err)
				}
				return nil
			},
			Example: "  jabba link system@1.8.20 /Library/Java/JavaVirtualMachines/jdk1.8.0_20.jdk\n" +
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
			Example: "  jabba unlink system@1.8.20",
		},
		&cobra.Command{
			Use:   "use [version to use]",
			Short: "Modify PATH & JAVA_HOME to use specific JDK",
			RunE: func(cmd *cobra.Command, args []string) error {
				var ver string
				if len(args) == 0 {
					ver = rc().JDK
					if ver == "" {
						return pflag.ErrHelp
					}
				} else {
					ver = args[0]
				}
				return use(ver)
			},
			Example: "  jabba use 1.8\n" +
				"  jabba use ~1.8.73 # same as \">=1.8.73 <1.9.0\"",
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
		lsCmd,
		lsRemoteCmd,
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
				name := args[0]
				if len(args) == 1 {
					if value := command.GetAlias(name); value != "" {
						fmt.Println(value)
					}
					return nil
				}
				if err := command.SetAlias(name, args[1]); err != nil {
					log.Fatal(err)
				}
				if err := command.LinkAlias(name); err != nil {
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
		whichCmd,
	)
	rootCmd.Flags().Bool("version", false, "version of jabba")
	rootCmd.PersistentFlags().String("fd3", "", "")
	rootCmd.PersistentFlags().MarkHidden("fd3")
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func parseTrimTo(value string) semver.VersionPart {
	switch strings.ToLower(value) {
	case "major":
		return semver.VPMajor
	case "minor":
		return semver.VPMinor
	case "patch":
		return semver.VPPatch
	default:
		log.Fatal("Unexpected value of --latest (must be either \"major\", \"minor\" or \"patch\")")
		return -1
	}
}

type jabbarc struct {
	JDK string
}

func rc() (rc jabbarc) {
	b, err := ioutil.ReadFile(".jabbarc")
	if err != nil {
		return
	}
	// content can be a string (jdk version)
	err = yaml.Unmarshal(b, &rc.JDK)
	if err != nil {
		// or a struct
		err = yaml.Unmarshal(b, &rc)
		if err != nil {
			log.Fatal(".jabbarc is not valid")
		}
	}
	return
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
