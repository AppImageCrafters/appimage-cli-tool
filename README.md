# appimage-cli-tool

A CLI app to manage your AppImage collection.

## Installation 

```shell script
sudo wget https://github.com/AppImageCrafters/appimage-cli-tool/releases/latest/download/appimage-cli-tool -O /usr/local/bin/appimage-cli-tool; 
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
