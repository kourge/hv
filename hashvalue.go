package main

import (
	"strings"
	"fmt"
	"crypto"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha512"
)
var _ = md5.New
var _ = sha1.New
var _ = sha512.New

type HashValue struct {
	crypto.Hash
}

func (h *HashValue) String() string {
	switch h.Hash {
	case crypto.MD4: return "MD4"
	case crypto.MD5: return "MD5"
	case crypto.SHA1: return "SHA1"
	case crypto.SHA224: return "SHA224"
	case crypto.SHA256: return "SHA256"
	case crypto.SHA384: return "SHA384"
	case crypto.SHA512: return "SHA512"
	case crypto.MD5SHA1: return "MD5SHA1"
	case crypto.RIPEMD160: return "RIPEMD160"
	default: return ""
	}
}

func (h *HashValue) Set(s string) (err error) {
	switch strings.ToUpper(s) {
	case "MD5": h.Hash = crypto.MD5
	case "SHA1": h.Hash = crypto.SHA1
	case "SHA512": h.Hash = crypto.SHA512
	default: err = HashUnavailableError{h.Hash}
	}
	return
}

func (h *HashValue) Values() []string {
	return []string{"MD5", "SHA1", "SHA512"}
}

type HashUnavailableError struct {
	h crypto.Hash
}

func (e HashUnavailableError) Error() string {
	h := &HashValue{e.h}
	hashName := h.String()
	if hashName == "" {
		return fmt.Sprintf("%#v is not a known hash function", e.h)
	} else {
		return fmt.Sprintf("%s is not supported or not linked into the binary", hashName)
	}
}

