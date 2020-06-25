package utils

// #include <stdio.h>
// #include <stdlib.h>
import "C"

import (
	"fmt"
	"github.com/rainycape/dl"
)

type libAppImageBind struct {
	lib *dl.DL

	appimage_shall_not_be_integrated func(path *C.char) int
	appimage_register_in_system      func(path *C.char, verbose int) int
	appimage_unregister_in_system    func(path *C.char, verbose int) int
}

type LibAppImage interface {
	Register(filePath string) error
	Unregister(filePath string) error
	ShallNotBeIntegrated(filePath string) bool
	Close()
}

func NewLibAppImageBindings() (LibAppImage, error) {
	bindings := libAppImageBind{}
	var err error
	bindings.lib, err = dl.Open("libappimage.so", 0)
	if err != nil {
		return nil, fmt.Errorf("desktop integration not available")
	}

	err = bindings.lib.Sym("appimage_shall_not_be_integrated", &bindings.appimage_shall_not_be_integrated)
	if err != nil {
		return nil, err
	}

	err = bindings.lib.Sym("appimage_unregister_in_system", &bindings.appimage_unregister_in_system)
	if err != nil {
		return nil, err
	}

	err = bindings.lib.Sym("appimage_register_in_system", &bindings.appimage_register_in_system)
	if err != nil {
		return nil, err
	}

	return &bindings, nil
}

func (bind *libAppImageBind) Register(filePath string) error {
	if bind.appimage_register_in_system(C.CString(filePath), 1) != 0 {
		return fmt.Errorf("unregister failed")
	}

	return nil
}

func (bind *libAppImageBind) ShallNotBeIntegrated(filePath string) bool {
	return bind.appimage_shall_not_be_integrated(C.CString(filePath)) != 0
}

func (bind *libAppImageBind) Unregister(filePath string) error {
	if bind.appimage_unregister_in_system(C.CString(filePath), 1) != 0 {
		return fmt.Errorf("unregister failed")
	}

	return nil
}

func (bind *libAppImageBind) Close() {
	_ = bind.lib.Close()
}
