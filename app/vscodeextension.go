package app

import (
	"context"
	"path/filepath"
)

type VSCodeExtensionInstaller struct {
	baseInstaller

	name, extensionName, vsixURL string
}

var _ Installer = (*VSCodeExtensionInstaller)(nil)

func NewVSCodeExtensionInstaller(baseDir, name, extensionName, vsixURL string) *VSCodeExtensionInstaller {
	i := VSCodeExtensionInstaller{
		name:          name,
		extensionName: extensionName,
		vsixURL:       vsixURL,
		baseInstaller: newBaseInstaller(filepath.Join(baseDir, name)),
	}
	return &i
}

func (i *VSCodeExtensionInstaller) Name() string {
	return i.name
}

func (i *VSCodeExtensionInstaller) BinName() string {
	return noExecutable
}

func (i *VSCodeExtensionInstaller) Version() string {
	return versionUnSpecified
}

func (i *VSCodeExtensionInstaller) Requires() []string {
	return noRequires
}

func (i *VSCodeExtensionInstaller) Install(ctx context.Context) error {
	return i.FetchWithExtract(ctx, i.vsixURL, filepath.Join(i.Dir(), i.extensionName+".zip"))
}
