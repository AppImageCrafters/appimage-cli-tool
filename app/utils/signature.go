package utils

import (
	"bytes"
	"crypto/sha256"
	"debug/elf"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/openpgp"
	"io"
	"os"
	"strings"
)

func VerifySignature(target string) (result *openpgp.Entity, err error) {
	key, err := readElfSection(target, ".sig_key")
	if err != nil {
		return nil, err
	}

	signature, err := readElfSection(target, ".sha256_sig")
	if err != nil {
		return nil, err
	}

	file, err := newAppImagePreSignatureReader(target)
	if err != nil {
		return
	}

	sha256Hash := sha256.New()
	_, err = io.Copy(sha256Hash, file)

	if err != nil {
		return nil, err
	}

	verification_target := hex.EncodeToString(sha256Hash.Sum(nil))

	keyring, err := openpgp.ReadArmoredKeyRing(bytes.NewReader(key))
	if err != nil {
		return nil, err
	}

	entity, err := openpgp.CheckArmoredDetachedSignature(keyring, strings.NewReader(verification_target), bytes.NewReader(signature))
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func readElfSection(appImagePath string, sectionName string) ([]byte, error) {
	elfFile, err := elf.Open(appImagePath)
	if err != nil {
		panic("Unable to open target: \"" + appImagePath + "\"." + err.Error())
	}

	section := elfFile.Section(sectionName)
	if section == nil {
		return nil, fmt.Errorf("missing " + sectionName + " section on target elf")
	}
	sectionData, err := section.Data()

	if err != nil {
		return nil, fmt.Errorf("Unable to parse " + sectionName + " section: " + err.Error())
	}

	str_end := bytes.Index(sectionData, []byte("\000"))
	if str_end == -1 || str_end == 0 {
		return nil, nil
	}

	return sectionData[:str_end], nil
}

func ReadSignature(appImagePath string) ([]byte, error) {
	elfFile, err := elf.Open(appImagePath)
	if err != nil {
		panic("Unable to open target: \"" + appImagePath + "\"." + err.Error())
	}

	updInfo := elfFile.Section(".sha256_sig")
	if updInfo == nil {
		panic("Missing .sha256_sig section on target elf ")
	}
	sectionData, err := updInfo.Data()

	if err != nil {
		panic("Unable to parse .sha256_sig section: " + err.Error())
	}

	str_end := bytes.Index(sectionData, []byte("\000"))
	if str_end == -1 || str_end == 0 {
		return nil, fmt.Errorf("No update information found in: " + appImagePath)
	}

	return sectionData[:str_end], nil
}

type appImagePreSignatureReader struct {
	keySectionOffset uint64
	keySectionSize   uint64

	sigSectionOffset uint64
	sigSectionSize   uint64

	offset uint64
	file   *os.File
}

func newAppImagePreSignatureReader(target string) (*appImagePreSignatureReader, error) {
	elfFile, err := elf.Open(target)
	if err != nil {
		return nil, err
	}

	key := elfFile.Section(".sig_key")
	if key == nil {
		return nil, fmt.Errorf("Missing .sig_key section")
	}

	signature := elfFile.Section(".sha256_sig")
	if signature == nil {
		return nil, fmt.Errorf("Missing .sha256_sig section")
	}

	file, err := os.Open(target)
	if err != nil {
		return nil, err
	}

	return &appImagePreSignatureReader{
		offset:           0,
		file:             file,
		keySectionOffset: key.Offset,
		keySectionSize:   key.Size,
		sigSectionOffset: signature.Offset,
		sigSectionSize:   signature.Size,
	}, nil
}

func (reader *appImagePreSignatureReader) Read(p []byte) (n int, err error) {
	n, err = reader.file.Read(p)
	if err != nil {
		return
	}

	oldOffset := reader.offset
	reader.offset += uint64(n)

	if reader.keySectionOffset >= oldOffset && reader.keySectionOffset < reader.offset {
		start := reader.keySectionOffset - oldOffset
		for i := start; i < uint64(n) && (i-start) < reader.keySectionSize; i++ {
			p[i] = 0
		}
	}

	if reader.sigSectionOffset >= oldOffset && reader.sigSectionOffset < reader.offset {
		start := reader.sigSectionOffset - oldOffset
		for i := start; i < uint64(n) && (i-start) < reader.sigSectionSize; i++ {
			p[i] = 0
		}
	}

	return n, err
}
