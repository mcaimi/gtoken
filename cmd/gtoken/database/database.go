package database

import (
  "os"
  "fmt"

  "github.com/spf13/cobra"
  "github.com/mcaimi/gtoken/pkg/gtoken"
  "github.com/mcaimi/gtoken/cmd/gtoken/styles"
  "github.com/jedib0t/go-pretty/v6/table"
)

var (
  DatabaseCmd = &cobra.Command{
    Use: "database",
    Short: "Manage the Account Database",
    Long: "Family of commands that allows for Database maintenance and other operations",
    Aliases: []string{"d", "db"},
  }

  statusCmd = &cobra.Command{
    Use: "status",
    Short: "Show status",
    Long: "Show information about the Account Database",
    Aliases: []string{"s", "stat", "info", "show"},
    Run: func(cmd *cobra.Command, args []string) {
      var d gtoken.Database;
      var err error;

      if d, err = gtoken.ReadAccountDb(); err != nil {
        fmt.Printf("gtoken status error: [%s]\n", err);
        os.Exit(1);
      }

      // format and print database status
      tbl := table.NewWriter();
      tbl.SetOutputMirror(os.Stdout);
      tbl.SetStyle(styles.GetStyle());
      tbl.AppendHeader(table.Row{"DB Path", "Entries"});
      tbl.AppendRow(table.Row{d.DbFilePath, d.Count()});
      tbl.Render();

      // ok
      os.Exit(0);
    },
  }
)

func init() {
  DatabaseCmd.AddCommand(statusCmd);
}
