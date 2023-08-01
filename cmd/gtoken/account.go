package gtoken

import (
  "os"
  "fmt"
  "github.com/spf13/cobra"

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
      tbl.AppendHeader(table.Row{"UUID", "Account Name", "E-Mail Address", "Type", "Flavor", "Validity", "Token"});
      for item := range tokens {
        tbl.AppendRow(table.Row{
          tokens[item].entry_uuid,
          tokens[item].account_name,
          tokens[item].email,
          tokens[item].totp_type,
          tokens[item].totp_flavor,
          tokens[item].totp_interval,
          tokens[item].token,
        });
      }
      tbl.Render();

      // ok
      os.Exit(0);
    },
  }

  insertCmd = &cobra.Command {
    Use: "insert",
    Short: "Add a new account",
    Long: "Add a new account to the Account DB. If the DB is empty, create a new one.",
    Aliases: []string{"i", "a", "ins", "add"},
    Run: func(cmd *cobra.Command, args []string) {
      // parse command flags
      name, _ := cmd.Flags().GetString("name");
      email, _ := cmd.Flags().GetString("email");
      tok_type, _ := cmd.Flags().GetString("type");
      tok_flavor, _ := cmd.Flags().GetString("flavor");
      tok_seed, _ := cmd.Flags().GetString("seed");
      tok_algo, _ := cmd.Flags().GetString("algorithm");
      interval, _ := cmd.Flags().GetInt("interval");

      // fill data into the new token object
      newAccount := tokenObject{account_name: name, email: email, totp_type: tok_type, totp_flavor: tok_flavor, totp_algo: tok_algo, token: tok_seed, totp_interval: interval};
      if err := ValidateToken(newAccount); err != nil {
        fmt.Printf("Token Insert Error: %s\n", err);
        os.Exit(1);
      }

      // insert new token into the database
      if err := InsertToken(newAccount); err != nil {
        fmt.Printf("Token Insert Error: %s\n", err);
        os.Exit(1);
      }

      // display newly inserted token
      tbl := table.NewWriter();
      tbl.SetOutputMirror(os.Stdout);
      tbl.AppendHeader(table.Row{"Account Name", "E-Mail Address", "Type", "Flavor", "Validity", "Token"});
      tbl.AppendRow(table.Row{
          newAccount.account_name,
          newAccount.email,
          newAccount.totp_type,
          newAccount.totp_flavor,
          newAccount.totp_interval,
          newAccount.token,
        });
      tbl.Render();

      // ok
      os.Exit(0);
    },
  }

  removeCmd = &cobra.Command {
    Use: "remove",
    Short: "Removes an account from the DB",
    Long: "Removes an account from the DB. After successful removal, the DB is automatically saved.",
    Aliases: []string{"d", "r", "rm", "del"},
    Run: func(cmd *cobra.Command, args []string) {
      // parse flags
      uuid, _ := cmd.Flags().GetString("uuid");
      if ! stringNotZeroLen(uuid) {
        fmt.Println("Please specify a valid UUID");
        os.Exit(1);
      }

      // remove account from the Database
      if err := RemoveToken(uuid); err != nil {
        fmt.Printf("Token Removal Error: %s\n", err);
        os.Exit(1);
      }

      // ok
      os.Exit(0);
    },
  }
)

func init() {
  // define flags (insert command)
  insertCmd.Flags().StringP("name", "n", "", "Specify The Account Name");
  insertCmd.Flags().StringP("email", "e", "", "E-Mail address associated with the Account");
  insertCmd.Flags().StringP("type", "t", "totp", "Specify The Account Name");
  insertCmd.Flags().StringP("flavor", "f", "google", "Two Factor Flavor (Google or RFC)");
  insertCmd.Flags().StringP("algorithm", "a", "sha1", "Token Hash Function");
  insertCmd.Flags().IntP("interval", "i", 30, "Token Refresh Interval (seconds)");
  insertCmd.Flags().StringP("seed", "s", "", "Token Secret Seed");

  // define flags (delete command)
  removeCmd.Flags().StringP("uuid", "u", "", "The UUID of the account that needs to be deleted from the Account DB");

  // build command
  accountCmd.AddCommand(listCmd);
  accountCmd.AddCommand(insertCmd);
  accountCmd.AddCommand(removeCmd);
  rootCommand.AddCommand(accountCmd);
}

