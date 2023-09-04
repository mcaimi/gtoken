package account

import (
  "os"
  "fmt"
  "github.com/spf13/cobra"

  "github.com/mcaimi/gtoken/pkg/database"
  "github.com/mcaimi/gtoken/pkg/token_io"
  "github.com/mcaimi/gtoken/pkg/common"
  "github.com/mcaimi/gtoken/cmd/gtoken/styles"
  "github.com/jedib0t/go-pretty/v6/table"
)

var (
  name string
  key string
  hash string
  interval int64
  token_type string
)

var (
  AccountCmd = &cobra.Command {
    Use: "account [verb]",
    Short: "Manage OTP Accounts",
    Long: "Manage the OTP Account database. Add, Remove or Update OTP Account Entries.",
    Aliases: []string{"a", "acct"},
    Run: func(cmd *cobra.Command, args []string) { generateCmd.Run(cmd, args); },
  }

  listCmd = &cobra.Command {
    Use: "list",
    Short: "List available token accounts",
    Long: "Displays a list of all OTP accounts currently configured in the database.",
    Aliases: []string{"l", "ls"},
    Run: func(cmd *cobra.Command, args []string) {
      // flags
      var accountUUID string;
      var showSeed bool;
      accountUUID, _ = cmd.Flags().GetString("show-seed");
      showSeed = common.StringNotZeroLen(accountUUID);

      if ! showSeed {
        tokens, tokenError := LoadTokens();
        if tokenError != nil {
          fmt.Printf("Cannot Load Tokens. Err: %s\n", tokenError);
          os.Exit(1);
        }
        // Display Tokens
        tbl := table.NewWriter();
        tbl.SetOutputMirror(os.Stdout);
        tbl.SetStyle(styles.GetStyle());

        tbl.AppendHeader(table.Row{"UUID", "Account Name", "E-Mail Address", "Type", "Flavor"});
        for item := range tokens {
          tbl.AppendRow(table.Row{
            tokens[item].UUID,
            tokens[item].Name,
            tokens[item].Email,
            tokens[item].Type,
            tokens[item].Flavor,
          });
        }
        tbl.Render();
      } else {
        token, tokenError := LoadToken(accountUUID);
        if tokenError != nil {
          fmt.Println(tokenError);
          os.Exit(1)
        } else {
          // Display Tokens
          tbl := table.NewWriter();
          tbl.SetOutputMirror(os.Stdout);
          tbl.SetStyle(styles.GetStyle());

          tbl.AppendHeader(table.Row{"UUID", "Account Name", "E-Mail Address", "Seed"});
          otpUrl := OtpUrl(token);
          tbl.AppendRow(table.Row{
              token.UUID,
              token.Name,
              token.Email,
              token.Key,
            });
          tbl.Render();

          fmt.Printf("\nOtp Provisioning Url: \n[%s]\n", otpUrl);
        }
      }

      // ok
      os.Exit(0);
    },
  }

  generateCmd = &cobra.Command {
    Use: "generate",
    Short: "Compute token values",
    Long: "Generates tokens for all registered accounts",
    Aliases: []string{"g", "gen"},
    Run: func(cmd *cobra.Command, args []string) {
      // load tokens
      var tokens []database.TokenEntity;
      var tokenError error;

      tokens, tokenError = GenerateTokens();
      if tokenError != nil {
        fmt.Printf("Cannot Generate Tokens. Err: %s\n", tokenError);
        os.Exit(1);
      }
      // Display Tokens
      tbl := table.NewWriter();
      tbl.SetOutputMirror(os.Stdout);
      tbl.SetStyle(styles.GetStyle());

      tbl.AppendHeader(table.Row{"Account Name", "E-Mail Address", "Interval", "Token Value"});
      for item := range tokens {
        tbl.AppendRow(table.Row{
          tokens[item].Name,
          tokens[item].Email,
          tokens[item].Interval,
          tokens[item].Token,
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
      interval, _ := cmd.Flags().GetInt64("interval");

      // fill data into the new token object
      newAccount := database.TokenEntity{Name: name, Email: email, Type: tok_type, Flavor: tok_flavor, Algorithm: tok_algo, Token: tok_seed, Interval: interval};
      if err := token_io.ValidateToken(newAccount); err != nil {
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
      tbl.SetStyle(styles.GetStyle());

      tbl.AppendHeader(table.Row{"Account Name", "E-Mail Address", "Type", "Flavor", "Validity", "Token"});
      tbl.AppendRow(table.Row{
          newAccount.Name,
          newAccount.Email,
          newAccount.Type,
          newAccount.Flavor,
          newAccount.Interval,
          newAccount.Token,
        });
      tbl.Render();

      // ok
      os.Exit(0);
    },
  }

  updateCmd = &cobra.Command {
    Use: "update",
    Short: "Update an exisiting account",
    Long: "Update data fields related to an existing account in the DB.",
    Aliases: []string{"u", "set"},
    Run: func(cmd *cobra.Command, args []string) {
      // parse command flags
      uuid, _ := cmd.Flags().GetString("uuid");
      name, _ := cmd.Flags().GetString("name");
      email, _ := cmd.Flags().GetString("email");
      tok_type, _ := cmd.Flags().GetString("type");
      tok_flavor, _ := cmd.Flags().GetString("flavor");
      tok_seed, _ := cmd.Flags().GetString("seed");
      tok_algo, _ := cmd.Flags().GetString("algorithm");
      interval, _ := cmd.Flags().GetInt64("interval");

      // flags
      uuidProvided := common.StringNotZeroLen(uuid);

      if ! uuidProvided {
        fmt.Printf("UUID Parameter is mandatory\n");
        os.Exit(1);
      }

      // fill data into the new token object
      newAccount := database.TokenEntity{Name: name, Email: email, Type: tok_type, Flavor: tok_flavor, Algorithm: tok_algo, Token: tok_seed, Interval: interval};
      //
      // update token
      if err := UpdateToken(uuid, newAccount); err != nil {
        fmt.Printf("Token Update Error: %s\n", err);
        os.Exit(1);
      }

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
      if ! common.StringNotZeroLen(uuid) {
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
  insertCmd.Flags().Int64P("interval", "i", 30, "Token Refresh Interval (seconds)");
  insertCmd.Flags().StringP("seed", "s", "", "Token Secret Seed");

  // define flags (update command)
  updateCmd.Flags().StringP("uuid", "u", "", "Update the account identified by the specified UUID");
  updateCmd.Flags().StringP("name", "n", "", "Specify The Account Name");
  updateCmd.Flags().StringP("email", "e", "", "E-Mail address associated with the Account");
  updateCmd.Flags().StringP("type", "t", "totp", "Specify The Account Name");
  updateCmd.Flags().StringP("flavor", "f", "google", "Two Factor Flavor (Google or RFC)");
  updateCmd.Flags().StringP("algorithm", "a", "sha1", "Token Hash Function");
  updateCmd.Flags().Int64P("interval", "i", 30, "Token Refresh Interval (seconds)");
  updateCmd.Flags().StringP("seed", "s", "", "Token Secret Seed");

  // define flags (list cmd)
  listCmd.Flags().StringP("show-seed", "s", "", "Show the Token Seed associated with the specific Account UUID");

  // define flags (delete command)
  removeCmd.Flags().StringP("uuid", "u", "", "The UUID of the account that needs to be deleted from the Account DB");

  // build command
  AccountCmd.AddCommand(listCmd);
  AccountCmd.AddCommand(insertCmd);
  AccountCmd.AddCommand(updateCmd);
  AccountCmd.AddCommand(generateCmd);
  AccountCmd.AddCommand(removeCmd);
}

