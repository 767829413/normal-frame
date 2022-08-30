package app

import (
	"fmt"
	"log"
	"os"

	cliflag "github.com/767829413/normal-frame/fork/component-base/cli/flag"
	"github.com/767829413/normal-frame/fork/component-base/cli/globalflag"
	"github.com/767829413/normal-frame/fork/component-base/cli/term"
	optionsCli "github.com/767829413/normal-frame/pkg/options"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	progressMessage = color.GreenString("==>")

	usageTemplate = fmt.Sprintf(`%s{{if .Runnable}}
  %s{{end}}{{if .HasAvailableSubCommands}}
  %s{{end}}{{if gt (len .Aliases) 0}}

%s
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

%s
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

%s{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  %s {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

%s
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

%s
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

%s{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "%s --help" for more information about a command.{{end}}
`,
		color.CyanString("Usage:"),
		color.GreenString("{{.UseLine}}"),
		color.GreenString("{{.CommandPath}} [command]"),
		color.CyanString("Aliases:"),
		color.CyanString("Examples:"),
		color.CyanString("Available Commands:"),
		color.GreenString("{{rpad .Name .NamePadding }}"),
		color.CyanString("Flags:"),
		color.CyanString("Global Flags:"),
		color.CyanString("Additional help topics:"),
		color.GreenString("{{.CommandPath}} [command]"),
	)
)

// App is the main structure of a cli application.
// It is recommended that an app be created with the app.NewApp() function.
type App struct {
	basename    string
	confName    string
	name        string
	description string
	options     optionsCli.CliOptions
	runFunc     RunFunc
	noConfig    bool
	commands    []*Command
	args        cobra.PositionalArgs
	cmd         *cobra.Command
}

// NewApp creates a new application instance based on the given application name,
// binary name, and other options.
func NewApp(name string, basename, confName string, opts ...Option) *App {
	a := &App{
		name:     name,
		basename: basename,
		confName: confName,
	}

	for _, o := range opts {
		o(a)
	}

	a.buildCommand()

	return a
}

func (a *App) buildCommand() {
	cmd := cobra.Command{
		Use:   FormatBaseName(a.basename),
		Short: a.name,
		Long:  a.description,
		// stop printing usage when the command errors
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          a.args,
	}
	// cmd.SetUsageTemplate(usageTemplate)
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.Flags().SortFlags = true
	cliflag.InitFlags(cmd.Flags())

	if len(a.commands) > 0 {
		for _, command := range a.commands {
			cmd.AddCommand(command.cobraCommand())
		}
		cmd.SetHelpCommand(helpCommand(FormatBaseName(a.basename)))
	}
	if a.runFunc != nil {
		cmd.RunE = a.runCommand
	}

	var namedFlagSets cliflag.NamedFlagSets
	if a.options != nil {
		namedFlagSets = a.options.Flags()
		fs := cmd.Flags()
		for _, f := range namedFlagSets.FlagSets {
			fs.AddFlagSet(f)
		}
	}

	if !a.noConfig {
		addConfigFlag(a.confName, namedFlagSets.FlagSet("global"))
	}
	globalflag.AddGlobalFlags(namedFlagSets.FlagSet("global"), cmd.Name())
	// add new global flagset to cmd FlagSet
	cmd.Flags().AddFlagSet(namedFlagSets.FlagSet("global"))

	addCmdTemplate(&cmd, namedFlagSets)
	a.cmd = &cmd
}

// Option defines optional parameters for initializing the application
// structure.
type Option func(*App)

// RunFunc defines the application's startup callback function.
type RunFunc func(basename string) error

// WithOptions to open the application's function to read from the command line
// or read parameters from the configuration file.
func WithOptions(opt optionsCli.CliOptions) Option {
	return func(a *App) {
		a.options = opt
	}
}

// WithRunFunc is used to set the application startup callback function option.
func WithRunFunc(run RunFunc) Option {
	return func(a *App) {
		a.runFunc = run
	}
}

// WithDescription is used to set the description of the application.
func WithDescription(desc string) Option {
	return func(a *App) {
		a.description = desc
	}
}

// WithDefaultValidArgs set default validation function to valid non-flag arguments.
func WithDefaultValidArgs() Option {
	return func(a *App) {
		a.args = func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}

			return nil
		}
	}
}

func addCmdTemplate(cmd *cobra.Command, namedFlagSets cliflag.NamedFlagSets) {
	usageFmt := "Usage:\n  %s\n"
	cols, _, _ := term.TerminalSize(cmd.OutOrStdout())
	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine())
		cliflag.PrintSections(cmd.OutOrStderr(), namedFlagSets, cols)

		return nil
	})
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
		cliflag.PrintSections(cmd.OutOrStdout(), namedFlagSets, cols)
	})
}

func (a *App) runCommand(cmd *cobra.Command, args []string) error {
	printWorkingDir()
	cliflag.PrintFlags(cmd.Flags())
	if !a.noConfig {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}

		if err := viper.Unmarshal(a.options); err != nil {
			return err
		}
	}

	// run application
	if a.runFunc != nil {
		return a.runFunc(a.basename)
	}

	return nil
}

func printWorkingDir() {
	wd, _ := os.Getwd()
	log.Printf("%v WorkingDir: %s", progressMessage, wd)
}

// Run is used to launch the application.
func (a *App) Run() {
	if err := a.cmd.Execute(); err != nil {
		log.Printf("%v %v\n", color.RedString("Error:"), err)
		os.Exit(1)
	}
}
