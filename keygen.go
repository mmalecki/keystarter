package main

import (
  "errors"
  "fmt"
  "net"
  "os/exec"
)

var openssl string = "openssl"

func privateKey() error {
  cmd := exec.Command(openssl, "genrsa", "-out", "server.key", "4096")
  return cmd.Run()
}

func csr(domains []string) error {
  cmd := exec.Command(openssl, "req", "new", "-key", "server.key", "-out", "server.csr")
  return cmd.Run()
}

func Keygen(domains []string) error {
  err := privateKey()
  if err != nil {
    return errors.New("Generating a private key failed")
  }

  ifaces, err := net.Interfaces()
  if err != nil {
    return errors.New("Getting network interfaces failed")
  }

  for i := range ifaces {
    if ifaces[i].Flags & net.FlagLoopback == net.FlagLoopback {
      continue
    }

    addrs, err := ifaces[i].Addrs()
    if err != nil {
      continue
    }

    for j := range addrs {
      ip, _, err := net.ParseCIDR(addrs[j].String())
      if err == nil {
        domains = append(domains, ip.String())
      }
    }
  }

  fmt.Println("Generating CSR for:", domains)

  err = csr(domains)
  if err != nil {
    return errors.New("Generating a CSR failed")
  }

  return nil
}
