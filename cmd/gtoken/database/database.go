package database

import (
  "os"
  "fmt"

  "github.com/spf13/cobra"
  "github.com/mcaimi/gtoken/pkg/token_io"
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

  initCmd = &cobra.Command{
    Use: "init",
    Short: "Initialize a new database",
    Long: "First time database setup, in which the relevant tables are created in a new db file.",
    Aliases: []string{"i", "in", "new"},
    Run: func(cmd *cobra.Command, args []string) {
      if err := token_io.InitDB(); err != nil {
        fmt.Printf("gtoken status error: [%s]\n", err);
        os.Exit(1);
      }
    },
  }

  statusCmd = &cobra.Command{
    Use: "status",
    Short: "Show status",
    Long: "Show information about the Account Database",
    Aliases: []string{"s", "stat", "info", "show"},
    Run: func(cmd *cobra.Command, args []string) {
      var d token_io.Database;
      var err error;

      if d, err = token_io.ReadAccountDb(); err != nil {
        fmt.Printf("gtoken status error: [%s]\n", err);
        os.Exit(1);
      }

      // format and print database status
      tbl := table.NewWriter();
      tbl.SetOutputMirror(os.Stdout);
      tbl.SetStyle(styles.GetStyle());
      tbl.AppendHeader(table.Row{"DB Path", "Entries", "Schema Version", "DB Checksum", "Valid"});
      tbl.AppendRow(table.Row{d.DbFilePath, d.Entries, d.Version, d.IntegrityChecksum, d.IsValid});
      tbl.Render();

      // ok
      os.Exit(0);
    },
  }
)

func init() {
  DatabaseCmd.AddCommand(initCmd);
  DatabaseCmd.AddCommand(statusCmd);
}
