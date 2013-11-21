package main

import (
  "fmt"
  "github.com/codegangsta/cli"
  "os"
)

func main() {
  app := cli.NewApp()
  app.Name = "keystarter"
  app.Usage = "Create a TLS key and upload it to the key server"

  app.Flags = []cli.Flag{
    cli.StringFlag{"key-server", "localhost", "Key server host"},
  }

  app.Commands = []cli.Command{
    {
      Name:  "add",
      Usage: "Add key to the key server",
      Action: func(c *cli.Context) {
        fmt.Println("Generating a self-signed certificate")
        Keygen([]string{})
      },
    },
    {
      Name:  "remove",
      Usage: "Remove key from the key server",
      Action: func(c *cli.Context) {
      },
    },
  }

  app.Run(os.Args)
}
