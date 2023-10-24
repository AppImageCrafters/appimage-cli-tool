CLI AppImage Management Tool
============================

Search, install, update and remove AppImage from the comfort of your CLI.

Features:
- Search/Install from the appimagehub.com catalog
- Install from github.com
- Update using the appimage-update
- Manage your local AppImage collection

## Installation 

```shell script
sudo wget https://github.com/AppImageCrafters/appimage-cli-tool/releases/download/continuous/appimage-cli-tool -O /usr/local/bin/appimage-cli-tool; 
sudo chmod +x /usr/local/bin/appimage-cli-tool
```

## Usage
```shell script
Usage: appimage-cli-tool <command>

Flags:
  --help     Show context-sensitive help.
  --debug    Enable debug mode.

Commands:
  search <query>
    Search applications in the store.

  install <target>
    Install an application.

  list
    List installed applications.

  remove <id>
    Remove an application.

  update [<targets> ...]
    Update an application.

Run "app <command> --help" for more information on a command.
```
