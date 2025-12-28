package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/pprof"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"

	"github.com/kpumuk/lazykiq/internal/ui"
)

var (
	Version = "dev"
	Commit  = ""
	Date    = ""
	BuiltBy = ""
)

var rootCmd = &cobra.Command{
	Use:   "lazykiq",
	Short: "A terminal UI for Sidekiq.",
	Long:  "A terminal UI for Sidekiq.",
	Args:  cobra.NoArgs,
}

func Execute() error {
	return fang.Execute(
		context.Background(),
		rootCmd,
		fang.WithVersion(rootCmd.Version),
		fang.WithoutCompletions(),
		fang.WithoutManpage(),
	)
}

func buildVersion(version, commit, date, builtBy string) string {
	result := version
	if commit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, commit)
	}
	if date != "" {
		result = fmt.Sprintf("%s\nbuilt at: %s", result, date)
	}
	if builtBy != "" {
		result = fmt.Sprintf("%s\nbuilt by: %s", result, builtBy)
	}
	result = fmt.Sprintf("%s\ngoos: %s\ngoarch: %s", result, runtime.GOOS, runtime.GOARCH)
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Sum != "" {
		result = fmt.Sprintf("%s\nmodule version: %s, checksum: %s", result, info.Main.Version, info.Main.Sum)
	}

	return result
}

func init() {
	rootCmd.Version = buildVersion(Version, Commit, Date, BuiltBy)
	rootCmd.SetVersionTemplate(`lazykiq {{printf "version %s\n" .Version}}`)

	rootCmd.Flags().String(
		"cpuprofile",
		"",
		"write cpu profile to file",
	)

	rootCmd.Flags().BoolP(
		"help",
		"h",
		false,
		"help for lazykiq",
	)

	rootCmd.RunE = func(cmd *cobra.Command, _ []string) error {
		cpuprofile, err := cmd.Flags().GetString("cpuprofile")
		if err != nil {
			return fmt.Errorf("parse cpuprofile flag: %w", err)
		}

		var profileFile *os.File
		if cpuprofile != "" {
			file, err := os.Create(cpuprofile)
			if err != nil {
				return fmt.Errorf("create cpuprofile file: %w", err)
			}
			profileFile = file
			if err := pprof.StartCPUProfile(profileFile); err != nil {
				_ = profileFile.Close()
				return fmt.Errorf("start cpu profile: %w", err)
			}
			defer func() {
				pprof.StopCPUProfile()
				_ = profileFile.Close()
			}()
		}

		app := ui.New()
		p := tea.NewProgram(app, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("run lazykiq: %w", err)
		}

		return nil
	}
}
