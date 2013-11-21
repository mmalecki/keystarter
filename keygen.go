package main

import (
  "os/exec"
  "errors"
)

var openssl string = "openssl"

func Keygen() (error) {
  cmd := exec.Command(openssl, "genrsa", "-out", "server.key", "4096")
  err := cmd.Run()
  if err != nil {
    return errors.New("Generating a private key failed")
  }

  return nil
}
