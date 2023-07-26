package gtoken

import (
  "os"
  "fmt"
  "github.com/spf13/cobra"

  "github.com/charmbracelet/bubbletea"
  "github.com/jedib0t/go-pretty/v6/table"
)

var (
  name string
  key string
  hash string
  interval int
  token_type string
)

var (
  accountCmd = &cobra.Command {
    Use: "account [verb]",
    Short: "Manage OTP Accounts",
    Long: "Manage the OTP Account database. Add, Remove or Update OTP Account Entries.",
    Aliases: []string{"a", "acct"},
    Run: func(cmd *cobra.Command, args []string) {},
  }

  listCmd = &cobra.Command {
    Use: "list",
    Short: "List available token accounts",
    Long: "Displays a list of all OTP accounts currently configured in the database.",
    Aliases: []string{"l", "ls"},
    Run: func(cmd *cobra.Command, args []string) {
      // load tokens
      var tokens []tokenObject;
      var tokenError error;

      tokens, tokenError = GenerateTokens();
      if tokenError != nil {
        fmt.Printf("Cannot Generate Tokens. Err: %s\n", tokenError);
        os.Exit(1);
      }
      // Display Tokens
      tbl := table.NewWriter();
      tbl.SetOutputMirror(os.Stdout);
      tbl.AppendHeader(table.Row{"Account Name", "E-Mail Address", "Type", "Flavor", "Validity", "Token"});
      for item := range tokens {
        tbl.AppendRow(table.Row{
          tokens[item].account_name,
          tokens[item].email,
          tokens[item].totp_type,
          tokens[item].totp_flavor,
          tokens[item].totp_interval,
          tokens[item].token,
        });
      }
      tbl.Render();

      os.Exit(0);
    },
  }

  insertCmd = &cobra.Command {
    Use: "insert",
    Short: "Add a new account",
    Long: "Add a new account to the Account DB. If the DB is empty, create a new one.",
    Aliases: []string{"i", "a", "ins", "add"},
    Run: func(cmd *cobra.Command, args []string) {
      if _, err := tea.NewProgram(AddAccountModel()).Run(); err != nil {
        fmt.Printf("Could Not Start Program. Err: %s\n", err);
        os.Exit(1);
      }
    },
  }

  removeCmd = &cobra.Command {
    Use: "remove",
    Short: "Removes an account from the DB",
    Long: "Removes an account from the DB. After successful removal, the DB is automatically saved.",
    Aliases: []string{"d", "r", "rm", "del"},
    Run: func(cmd *cobra.Command, args []string) {},
  }
)

func init() {
  // build command
  accountCmd.AddCommand(listCmd);
  accountCmd.AddCommand(insertCmd);
  accountCmd.AddCommand(removeCmd);
  rootCommand.AddCommand(accountCmd);
}

