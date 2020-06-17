package install

import "errors"

var InvalidTargetFormat = errors.New("invalid target format")

var NoAppImageBinariesFound = errors.New("no AppImage found in releases")

var ApplicationInstalled = errors.New("the application is installed already")
