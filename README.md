## Linkkeeper

This is a simple bookmark manager/reader TUI app I am building to play around with Go and the Charm ecosystem.

### Getting Started (Dev mode)

```
devcontainer up --workspace-folder .
code --folder-uri vscode-remote://attached-container+$(printf "$(docker ps -lq)" | xxd -p)/code
```