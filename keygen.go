package main

import (
  "os/exec"
  "errors"
)

var openssl string = "openssl"

func privateKey() (error) {
  cmd := exec.Command(openssl, "genrsa", "-out", "server.key", "4096")
  return cmd.Run()
}

func Keygen() (error) {
  err := privateKey()
  if err != nil {
    return errors.New("Generating a private key failed")
  }

  return nil
}
