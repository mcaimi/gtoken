package gtoken

import (
  "os"
  "fmt"
  "github.com/spf13/cobra"

  "github.com/charmbracelet/bubbletea"
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
      if _, err := tea.NewProgram(TableInitialModel()).Run(); err != nil {
        fmt.Printf("Could Not Start Program. Err: %s\n", err);
        os.Exit(1);
      }
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

