package terminal

import (
	"fmt"
	"io"
	"log"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
)

// UI is a terminal UI
type UI interface {
	AutoConfirm() bool
	Ask(answer interface{}, questions ...*survey.Question) error
	AskOne(answer interface{}, prompt survey.Prompt) error
	Input(answer interface{}, askOptions AskOptions) error
	Select(answer interface{}, askOptions AskOptions) error
	MultiSelect(answer interface{}, askOptions AskOptions) error
	Password(answer interface{}, askOptions AskOptions) error
	Confirm(format string, args ...interface{}) (bool, error)
	Print(logs ...Log)
	Spinner(message string, opts SpinnerOptions) Spinner
}

// NewUI creates a new terminal UI
func NewUI(config UIConfig, in io.Reader, out, err io.Writer) UI {
	minimalUI := config.DisableColors
	if config.OutputFormat == OutputFormatJSON {
		minimalUI = true
	}
	color.NoColor = minimalUI

	return &ui{
		config,
		minimalUI,
		fdReader{in},
		fdWriter{out},
		err,
	}
}

// AskOptions mirror the config details for survey prompts
type AskOptions struct {
	Message string
	Default string
	Options []string
}

type ui struct {
	config  UIConfig
	minimal bool
	in      fdReader
	out     fdWriter
	err     io.Writer
}

func (ui *ui) AutoConfirm() bool {
	return ui.config.AutoConfirm
}

func (ui *ui) Ask(answer interface{}, questions ...*survey.Question) error {
	return survey.Ask(
		questions,
		answer,
		survey.WithStdio(ui.in, ui.out, ui.err),
	)
}

func (ui *ui) AskOne(answer interface{}, prompt survey.Prompt) error {
	return survey.AskOne(
		prompt,
		answer,
		survey.WithStdio(ui.in, ui.out, ui.err),
	)
}

func (ui *ui) Input(answer interface{}, askOptions AskOptions) error {
	prompt := &survey.Input{
		Message: askOptions.Message,
		Default: askOptions.Default,
	}
	return survey.AskOne(
		prompt,
		answer,
		survey.WithStdio(ui.in, ui.out, ui.err),
	)
}

func (ui *ui) Password(answer interface{}, askOptions AskOptions) error {
	return survey.AskOne(
		&survey.Password{Message: askOptions.Message},
		answer,
		survey.WithStdio(ui.in, ui.out, ui.err),
	)
}

func (ui *ui) Select(answer interface{}, askOptions AskOptions) error {
	prompt := &survey.Select{
		Message: askOptions.Message,
		Options: askOptions.Options,
		Default: askOptions.Default,
	}
	return survey.AskOne(
		prompt,
		answer,
		survey.WithStdio(ui.in, ui.out, ui.err),
	)
}

func (ui *ui) MultiSelect(answer interface{}, askOptions AskOptions) error {
	prompt := &survey.MultiSelect{
		Message: askOptions.Message,
		Options: askOptions.Options,
	}
	return survey.AskOne(
		prompt,
		answer,
		survey.WithStdio(ui.in, ui.out, ui.err),
	)
}

func (ui *ui) Confirm(format string, args ...interface{}) (bool, error) {
	if ui.AutoConfirm() {
		return true, nil
	}

	var proceed bool
	return proceed, ui.AskOne(
		&proceed,
		&survey.Confirm{Message: fmt.Sprintf(format, args...)},
	)
}

func (ui *ui) Print(logs ...Log) {
	for _, l := range logs {
		output, err := l.Print(ui.config.OutputFormat)
		if err != nil {
			ui.Print(NewErrorLog(err))
			return
		}

		var writer io.Writer
		switch l.Level {
		case LogLevelError:
			writer = ui.err
		default:
			writer = ui.out
		}

		if _, err := fmt.Fprintln(writer, output); err != nil {
			log.Print(output) // log the original output
		}
	}
}

func (ui *ui) Spinner(message string, opts SpinnerOptions) Spinner {
	if ui.minimal {
		return noopSpinner{}
	}
	return newUISpinner(message, opts)
}

// UIConfig holds the global config for the CLI ui
type UIConfig struct {
	AutoConfirm   bool
	DisableColors bool
	OutputFormat  OutputFormat
	OutputTarget  string
}

// FileDescriptor is a file descriptor
type FileDescriptor interface {
	Fd() uintptr
}

// fdReader wraps an io.Reader and exposes the FileDesriptor interface on it
// the underlying io.Reader's Fd() implementation will be used if it exists
type fdReader struct {
	io.Reader
}

func (r fdReader) Fd() uintptr {
	if fd, ok := r.Reader.(FileDescriptor); ok {
		return fd.Fd()
	}
	return 0
}

// fdWriter wraps an io.Writer and exposes the FileDescriptor interface on it
// the underlying io.Writer's Fd() implementation will be used if it exists
type fdWriter struct {
	io.Writer
}

func (w fdWriter) Fd() uintptr {
	if fd, ok := w.Writer.(FileDescriptor); ok {
		return fd.Fd()
	}
	return 0
}
