package app

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
			"gopls":                             NewGoInstaller(baseDir, "golang.org/x/tools/gopls", "gopls"),
			"metals":                            NewMetalsInstaller(baseDir),
			"kotlin-language-server":            NewKotlinLSInstaller(baseDir),
			"rust-analyzer":                     NewRustAnalyzerInstaller(baseDir),
			"efm-langserver":                    NewEfmLSInstaller(baseDir),
			"vscode-css-languageserver":         NewNpmInstaller(baseDir, "vscode-css-languageserver-bin", "css-languageserver"),
			"vscode-html-languageserver":        NewNpmInstaller(baseDir, "vscode-html-languageserver-bin", "html-languageserver"),
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

func (a *App) List(ctx context.Context) error {
	if err := os.MkdirAll(a.baseDir, 0777); err != nil {
		return err
	}
	dirs, err := ioutil.ReadDir(a.baseDir)
	if err != nil {
		return err
	}
	buf := bufio.NewWriter(a.out)
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
		var msg string
		if found {
			msg = fmt.Sprintf("%s is installed\n", i.Name())
		} else {
			msg = fmt.Sprintf("%s is not installed\n", i.Name())
		}
		if _, err := buf.WriteString(msg); err != nil {
			return err
		}
	}
	if err := buf.Flush(); err != nil {
		return err
	}
	return nil
}

func isExecutable(mode os.FileMode) bool {
	// FIXME
	if isWindows {
		return true
	}
	return mode&0100 != 0
}
