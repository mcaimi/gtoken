package gtoken

import (
  "fmt"
  "os"

  "github.com/spf13/cobra"
  "github.com/mcaimi/gtoken/cmd/gtoken/account"
  "github.com/mcaimi/gtoken/cmd/gtoken/database"
)

var (
  configFile string
  rootCommand = &cobra.Command{
    Use: "gtoken",
    Short: "Auth Token Generation Utility",
    Long: "CLI utility that implements various OTP generation algorithms based on RFC4226 and RFC6238.",
  }
)

func init() {
  rootCommand.AddCommand(account.AccountCmd);
  rootCommand.AddCommand(database.DatabaseCmd);
}

func Execute() {
  if err := rootCommand.Execute(); err != nil {
    fmt.Fprintf(os.Stderr, err.Error());
    os.Exit(1);
  }
}
