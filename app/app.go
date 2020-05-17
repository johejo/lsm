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

const (
	linux   = "linux"
	darwin  = "darwin"
	windows = "windows"
	freebsd = "freebsd"
	openbsd = "openbsd"
	solaris = "solaris"

	amd64 = "amd64"
	_386  = "386"
	arm   = "arm"

	appName = "lsm"
	servers = "servers"
)

var isWindows bool

func init() {
	if runtime.GOOS == windows {
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

func getBaseDir() (string, error) {
	var baseDir string
	switch runtime.GOOS {
	case linux, darwin:
		xdgDataHome := os.Getenv("XDG_DATA_HOME")
		if xdgDataHome != "" {
			baseDir = filepath.Join(xdgDataHome, appName, servers)
			break
		}
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		baseDir = filepath.Join(home, ".local", "share", appName, servers)
	case windows:
		local := os.Getenv("LOCALAPPDATA")
		if local == "" {
			return "", errors.New("LOCALAPPDATA is not defined")
		}
		baseDir = filepath.Join(local, appName, servers)
	default:
		return "", fmt.Errorf("unsupported operating system: %v", runtime.GOOS)
	}
	return filepath.Abs(baseDir)
}

func New(baseDir string) (*App, error) {
	if baseDir == "" {
		p, err := getBaseDir()
		if err != nil {
			return nil, err
		}
		baseDir = p
	} else {
		p, err := filepath.Abs(baseDir)
		if err != nil {
			return nil, err
		}
		baseDir = p
	}
	installers := map[string]Installer{
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
		"terraform-lsp":                     NewTerraformLSPInstaller(baseDir),
		"terraform-ls":                      NewTerraformLSInstaller(baseDir),
		"eclipse.jdt.ls":                    NewEclipseJDTLSInstaller(baseDir),
	}

	return &App{
		baseDir:    baseDir,
		installers: installers,
		out:        os.Stdout,
	}, nil
}

func (a *App) getInstaller(name string) (Installer, error) {
	i, ok := a.installers[name]
	if !ok {
		return nil, errors.New(name + " not found")
	}
	return i, nil
}

func isSupported(i Installer) error {
	ss := i.Supports()
	if len(ss) == 0 {
		return nil
	}
	for _, s := range ss {
		if runtime.GOARCH == s.arch && runtime.GOOS == s.os {
			return nil
		}
	}
	return fmt.Errorf("installer does not supports %s on %s %s", i.Name(), runtime.GOOS, runtime.GOARCH)
}

func (a *App) Install(ctx context.Context, name string) error {
	i, err := a.getInstaller(name)
	if err != nil {
		return err
	}

	if err := isSupported(i); err != nil {
		return err
	}
	for _, r := range i.Requires() {
		if _, err := exec.LookPath(r); err != nil {
			return err
		}
	}
	if err := i.RequireHook(ctx); err != nil {
		return err
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
		installed := false
		for _, d := range dirs {
			if d.IsDir() && d.Name() == i.Name() {
				bin := filepath.Join(a.baseDir, i.Name(), i.BinName())
				info, err := os.Stat(bin)
				if err != nil {
					installed = false
					continue
				}
				if isExecutable(info.Mode()) {
					installed = true
					break
				}
			}
		}
		list = append(list, languageServer{
			Name:      i.Name(),
			Version:   i.Version(),
			Installed: installed,
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
