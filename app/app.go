package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/olekukonko/tablewriter"
)

var isWindows bool

func init() {
	if runtime.GOOS == "windows" {
		isWindows = true
	} else {
		isWindows = false
	}
}

type App struct {
	installers map[string]Installer
	baseDir    string
	out        io.Writer
}

func New(baseDir string) (*App, error) {
	if baseDir == "" {
		switch runtime.GOOS {
		case "linux", "darwin":
			xdgDataHome := os.Getenv("XDG_DATA_HOME")
			if xdgDataHome != "" {
				baseDir = filepath.Join(xdgDataHome, "lsm", "servers")
				break
			}
			home, err := os.UserHomeDir()
			if err != nil {
				return nil, err
			}
			baseDir = filepath.Join(home, ".local", "share", "lsm", "servers")
		case "windows":
			baseDir = filepath.Clean(`%LOCALAPPDATA%\lsm\servers`)
		}
	} else {
		var err error
		baseDir, err = filepath.Abs(filepath.Clean(baseDir))
		if err != nil {
			return nil, err
		}
	}
	return &App{
		baseDir: baseDir,
		installers: map[string]Installer{
			"vim-language-server":               NewNpmInstaller(baseDir, "vim-language-server", "vim-language-server"),
			"typescript-language-server":        NewNpmInstaller(baseDir, "typescript-language-server", "typescript-language-server"),
			"dockerfile-language-server-nodejs": NewNpmInstaller(baseDir, "dockerfile-language-server-nodejs", "docker-langserver"),
			"bash-language-server":              NewNpmInstaller(baseDir, "bash-language-server", "bash-language-server"),
			"yaml-language-server":              NewNpmInstaller(baseDir, "yaml-language-server", "yaml-language-server"),
			"vscode-json-languageserver":        NewNpmInstaller(baseDir, "vscode-json-languageserver", "vscode-json-languageserver"),
			"vls":                               NewNpmInstaller(baseDir, "vls", "vls"),
			"gopls":                             NewGoInstaller(baseDir, "golang.org/x/tools/gopls", "gopls", false),
			"sqls":                              NewGoInstaller(baseDir, "github.com/lighttiger2505/sqls", "sqls", true),
			"metals":                            NewMetalsInstaller(baseDir),
			"kotlin-language-server":            NewKotlinLSInstaller(baseDir),
			"rust-analyzer":                     NewRustAnalyzerInstaller(baseDir),
			"efm-langserver":                    NewEfmLSInstaller(baseDir),
			"vscode-css-languageserver":         NewNpmInstaller(baseDir, "vscode-css-languageserver-bin", "css-languageserver"),
			"vscode-html-languageserver":        NewNpmInstaller(baseDir, "vscode-html-languageserver-bin", "html-languageserver"),
			"python-language-server":            NewPipInstaller(baseDir, "python-language-server", "pyls"),
			"fortran-language-server":           NewPipInstaller(baseDir, "fortran-language-server", "fortls"),
		},
		out: os.Stdout,
	}, nil
}

func (a *App) getInstaller(name string) (Installer, error) {
	i, ok := a.installers[name]
	if !ok {
		return nil, errors.New(name + " not found")
	}
	return i, nil
}

func (a *App) Install(ctx context.Context, name string) error {
	i, err := a.getInstaller(name)
	if err != nil {
		return err
	}

	for _, r := range i.Requires() {
		if _, err := exec.LookPath(r); err != nil {
			return err
		}
	}
	if hook := i.RequiresHook(); hook != nil {
		if err := hook(ctx); err != nil {
			return err
		}
	}

	if err := os.RemoveAll(i.Dir()); err != nil {
		return err
	}
	if err := os.MkdirAll(i.Dir(), 0777); err != nil {
		return err
	}
	if err := i.Install(ctx); err != nil {
		return err
	}
	log.Printf("%s %s installed into %s", name, i.Version(), i.Dir())
	return nil
}

func (a *App) Uninstall(ctx context.Context, name string) error {
	i, err := a.getInstaller(name)
	if err != nil {
		return err
	}
	if err := os.RemoveAll(i.Dir()); err != nil {
		return err
	}
	log.Printf("%s uninstalled from %s", name, i.Dir())
	return nil
}

type ListStyle string

const (
	ListStyleUndefined ListStyle = ""
	ListStyleJSON      ListStyle = "json"
	ListStyleTable     ListStyle = "table"
)

func (a *App) List(ctx context.Context, style ListStyle) error {
	if err := os.MkdirAll(a.baseDir, 0777); err != nil {
		return err
	}
	dirs, err := ioutil.ReadDir(a.baseDir)
	if err != nil {
		return err
	}
	list := make([]languageServer, 0, len(a.installers))
	for _, i := range a.installers {
		found := false
		for _, d := range dirs {
			if d.IsDir() && d.Name() == i.Name() {
				bin := filepath.Join(a.baseDir, i.Name(), i.BinName())
				info, err := os.Stat(bin)
				if err != nil {
					found = false
					continue
				}
				if isExecutable(info.Mode()) {
					found = true
					break
				}
			}
		}
		list = append(list, languageServer{
			Name:      i.Name(),
			Version:   i.Version(),
			Installed: found,
		})
	}
	switch style {
	case ListStyleJSON:
		return a.renderJSON(list)
	case ListStyleTable, ListStyleUndefined:
		return a.renderTable(list)
	default:
		return fmt.Errorf("unsupported list style: %v", style)
	}
}

func (a *App) renderJSON(list []languageServer) error {
	b, err := json.MarshalIndent(list, "", strings.Repeat(" ", 2))
	if err != nil {
		return err
	}
	if _, err := a.out.Write(b); err != nil {
		return err
	}
	return nil
}

func (a *App) renderTable(list []languageServer) error {
	table := tablewriter.NewWriter(a.out)
	t := reflect.TypeOf(languageServer{})
	headers := make([]string, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		headers = append(headers, f.Name)
	}
	table.SetHeader(headers)
	for _, item := range list {
		v := reflect.ValueOf(item)
		t := reflect.TypeOf(item)
		rows := make([]string, 0, t.NumField())
		for i := 0; i < t.NumField(); i++ {
			rows = append(rows, fmt.Sprint(v.Field(i)))
		}
		table.Append(rows)
	}
	table.Render()
	return nil
}

func isExecutable(mode os.FileMode) bool {
	// FIXME
	if isWindows {
		return true
	}
	return mode&0100 != 0
}

type languageServer struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	Installed bool   `json:"installed"`
}
