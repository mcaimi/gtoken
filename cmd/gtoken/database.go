package gtoken

import (
  "os"
  "fmt"

  "github.com/spf13/cobra"
  "github.com/mcaimi/gtoken/pkg/common"
  "github.com/mcaimi/gtoken/pkg/gtoken"
)

var (
  databaseCommand = &cobra.Command{
    Use: "database",
    Short: "Manage the Account Database",
    Long: "Family of commands that allows for Database maintenance and other operations",
    Aliases: []string{"d", "db"},
  }

  statusCommand = &cobra.Command{
    Use: "status",
    Short: "Show status",
    Long: "Show information about the Account Database",
    Aliases: []string{"s", "stat", "info", "show"},
    Run: func(cmd *cobra.Command, args []string) {
      if dbPath, err := common.GetAccountsDB(); err == nil {
        fmt.Printf("Database Location: [%s]\n", dbPath);

        // count entries
        if d, err := gtoken.OpenAccountDB(dbPath); err == nil {
          fmt.Printf("\t-> Currently, the Database contains [%d] configured entries.\n", d.Count());
        }
      } else {
        fmt.Printf("gtoken status error: [%s]\n", err);
        os.Exit(1);
      }

      // ok
      os.Exit(0);
    },
  }
)

func init() {
  databaseCommand.AddCommand(statusCommand);
  rootCommand.AddCommand(databaseCommand);
}
