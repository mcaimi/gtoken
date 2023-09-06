package database

import (
  "os"
  "fmt"

  "github.com/spf13/cobra"
  "github.com/mcaimi/gtoken/pkg/token_io"
  "github.com/mcaimi/gtoken/pkg/common"
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

  rehashCmd = &cobra.Command{
    Use: "rehash",
    Short: "Validates entries and rehashes the Database.",
    Long: "Checks validity of each database entry and recomputes the DB checksum if everything is ok.",
    Aliases: []string{"r", "rehash", "fix", "verify"},
    Run: func(cmd *cobra.Command, args []string) {
      if err := token_io.ValidateDatabase(); err != nil {
        fmt.Printf("gtoken status error: [%s]\n", err);
        os.Exit(1);
      }
    },
  }

  dumpCmd = &cobra.Command{
    Use: "dump",
    Short: "Dumps DB contents to a JSON file",
    Long: "Dumps all accounts stored into the Database into a JSON file.",
    Aliases: []string{"d", "dump", "export"},
    Run: func(cmd *cobra.Command, args []string) {
      fName, _ := cmd.Flags().GetString("filename");
      dumpOk := common.StringNotZeroLen(fName);

      if ! dumpOk {
        fmt.Println("Invalid output Filename");
        os.Exit(1);
      } else {
        if err := token_io.JsonDump(fName); err != nil {
          fmt.Printf("gtoken status error: [%s]\n", err);
          os.Exit(1);
        }
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
  // parameters
  dumpCmd.Flags().StringP("filename", "f", "", "File name where to dump account data.");

  DatabaseCmd.AddCommand(initCmd);
  DatabaseCmd.AddCommand(rehashCmd);
  DatabaseCmd.AddCommand(dumpCmd);
  DatabaseCmd.AddCommand(statusCmd);
}
