package main

import (
  "strconv"
  "errors"
  "fmt"
  "net"
  "os"
  "os/exec"
)

var openssl string = "openssl"

func privateKey() error {
  cmd := exec.Command(openssl, "genrsa", "-out", "server.key", "4096")
  return cmd.Run()
}

func csr(dns []string, ip[]string) error {
  cfg, err := os.Create("openssl.cfg")
  if err != nil {
    return err
  }

  hostname, err := os.Hostname()
  if err != nil {
    hostname = "localhost"
  }

  cfg.WriteString("[req]\n")
  cfg.WriteString("distinguished_name = req_distinguished_name\n")
  cfg.WriteString("req_extensions = v3_req\n")
  cfg.WriteString("[req_distinguished_name]\n")
  cfg.WriteString("countryName = Country Name (2 letter code)\n")
  cfg.WriteString("countryName_default = .\n")
  cfg.WriteString("stateOrProvinceName = State or Province Name (full name)\n")
  cfg.WriteString("stateOrProvinceName_default = .\n")
  cfg.WriteString("localityName = Locality Name (eg, city)\n")
  cfg.WriteString("localityName_default = .\n")
  cfg.WriteString("organizationalUnitName = Organizational Unit Name (eg, section)\n")
  cfg.WriteString("organizationalUnitName_default = .\n")
  cfg.WriteString("commonName = " + hostname + "\n")
  cfg.WriteString("commonName_max = 64\n")
  cfg.WriteString("[v3_req]\n")
  cfg.WriteString("basicConstraints = CA:FALSE\n")
  cfg.WriteString("keyUsage = nonRepudiation, digitalSignature, keyEncipherment\n")
  cfg.WriteString("subjectAltName = @alt_names\n")
  cfg.WriteString("[alt_names]\n")

  for i := range(dns) {
    cfg.WriteString("DNS." + strconv.Itoa(i + 1) + " = " + dns[i] + "\n")
  }

  for i := range(ip) {
    cfg.WriteString("IP." + strconv.Itoa(i + 1) + " = " + ip[i] + "\n")
  }

  err = cfg.Close()
  if err != nil {
    return err
  }

  // MULTILINEEE, jelly much ( ≖‿≖)?
  cmd := exec.Command(openssl, "req",
    "-new", "-key", "server.key", "-out",
    "server.csr","-config", "openssl.cfg",
    "-subj", "/CN=" + hostname)
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

  ips := []string{}

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
        ips = append(ips, ip.String())
      }
    }
  }

  fmt.Println("Generating CSR for:", domains, ips)

  err = csr(domains, ips)
  if err != nil {
    return errors.New("Generating a CSR failed")
  }

  return nil
}
