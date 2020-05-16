# Language Server Manager (LSM)

[![ci](https://github.com/johejo/lsm/workflows/ci/badge.svg)](https://github.com/johejo/lsm/actions?query=workflow%3Aci)
[![codecov](https://codecov.io/gh/johejo/lsm/branch/master/graph/badge.svg)](https://codecov.io/gh/johejo/lsm)

## Motivation

Language Server is difficult to set up and is tightly coupled with each editor plugin.<br>
https://github.com/neovim/nvim-lsp/issues/200#issuecomment-618807759

There are also lots of scripts to support different platforms and Language Servers, which is becoming difficult to maintain.<br>
https://github.com/mattn/vim-lsp-settings/tree/master/installer

LSM aims to be a simple and cross platform installation manager for Language Server.

## Notice
Windows support is still incomplete.

## Language Server Destination

### macOS, Linux

```
$HOME/.local/share/lsm/servers
```

If you defined `$XDG_DATA_HOME`

```
$XDG_DATA_HOME/lsm/servers
```

### Windows

```
%LOCALAPPDATA%\lsm\servers
```

## Install

go get
```
go get github.com/johejo/lsm
```

pre-built binary
```
# COMMING SOON
```

## Usage

For example, when gopls is the target

install
```
lsm install gopls
```

uninstall
```
lsm uninstall gopls
```

list (Installation status of Language Servers)

```
lsm list
```

## Supported Language Servers

- [gopls](https://github.com/golang/tools/tree/master/gopls)
- [vim-language-server](https://github.com/iamcco/vim-language-server)
- [bash-language-server](https://github.com/bash-lsp/bash-language-server)
- [typescript-language-server](https://github.com/theia-ide/typescript-language-server)
- [kotlin-language-server](https://github.com/fwcd/kotlin-language-server)
- [yaml-language-sever](https://github.com/redhat-developer/yaml-language-server/)
- [vscode-json-languagesever](https://github.com/vscode-langservers/vscode-json-languageserver)
- [vscode-css-languagesever](https://github.com/vscode-langservers/vscode-css-languageserver)
- [vscode-html-languagesever](https://github.com/vscode-langservers/vscode-html-languageserver)
- [dockerfile-language-server-nodejs](https://github.com/rcjsuen/dockerfile-language-server-nodejs)
- [vls](https://github.com/vuejs/vetur/tree/master/server)
- [metals](https://scalameta.org/metals/)
- [rust-analyzer](https://rust-analyzer.github.io/)
- [efm-langserver](https://github.com/mattn/efm-langserver)
- [python-language-server](https://github.com/palantir/python-language-server)
- [fortran-language-server](https://github.com/hansec/fortran-language-server)
- [sqls](https://github.com/lighttiger2505/sqls)

We plan to add more Language Servers.<br>
If you want to add a new Language Server, please create a PR.<br>

## Thanks

Language Server Manager inherits the concept of [vim-lsp-settings](https://github.com/mattn/vim-lsp-settings).
