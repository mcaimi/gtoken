package gtoken

import (
  "fmt"
  "strings"

  "github.com/charmbracelet/bubbletea"
  "github.com/charmbracelet/bubbles/table"

  "github.com/mcaimi/gtoken/pkg/gtoken"
  "github.com/mcaimi/gtoken/pkg/common"

  "github.com/mcaimi/go-totp/rfc6238"
)

var columnLabels map[int]string = map[int]string {
  account_name: "Account Name",
  email: "E-Mail Address",
  totp_algo: "Hashing Algorithm",
  totp_flavor: "Totp Flavor",
  totp_interval: "Token Validity",
  totp_type: "Token Type",
  totp_uuid: "Unique ID",
}

type totptable struct {
  table table.Model
  cliError error
}

func TableInitialModel() totptable {
  var (
    t totptable;
    acctDb gtoken.Database;
    columns []table.Column;
    rows []table.Row;
  )
  t.cliError = nil;

  columns = []table.Column {
    {Title: columnLabels[account_name], Width: len(columnLabels[account_name]) + pad},
    {Title: columnLabels[email], Width: len(columnLabels[email]) + pad},
    {Title: columnLabels[totp_algo], Width: len(columnLabels[totp_algo]) + pad},
    {Title: columnLabels[totp_flavor], Width: len(columnLabels[totp_flavor]) + pad},
    {Title: columnLabels[totp_interval], Width: len(columnLabels[totp_interval]) + pad},
    {Title: columnLabels[totp_type], Width: len(columnLabels[totp_type]) + pad},
    {Title: columnLabels[totp_uuid], Width: len(columnLabels[totp_uuid]) + pad},
    {Title: "Token", Width: len("Token") + pad},
  }

  // Read database contents from disk
  accountFile, err := common.GetAccountsDB();
  if err != nil {
    t.cliError = err;
    return t;
  }
  acctDb, err = gtoken.OpenAccountDB(accountFile);
  if err != nil {
    t.cliError = err;
    return t;
  }

  var a gtoken.Account;
  var n int = acctDb.Count();
  rows = make([]table.Row, n);
  for i := range acctDb.Accounts {
    a = acctDb.Accounts[i];

    // compute totp token
    var token string;
    if a.Flavor == "google" {
      token = rfc6238.GoogleAuth([]byte(a.Key), 6);
    } else {
      token = "Not Implemented";
    }

    // update table data
    rows[i] = table.Row{a.Name, a.Email, a.Hash, a.Flavor, fmt.Sprintf("%d", a.Interval), a.Type, a.Uuid, token};
  }

  // build table object
  t.table = table.New(
    table.WithColumns(columns),
    table.WithRows(rows),
    table.WithFocused(true),
    table.WithHeight(n),
  );

  style := table.DefaultStyles();
  style.Header = style.Header.
		BorderStyle(normalBorder).
		BorderForeground(white).
		BorderBottom(true).
		Bold(false)
	style.Selected = style.Selected.
		Foreground(yellow).
		Background(blue).
		Bold(false)
	t.table.SetStyles(style)

  return t;
}

func (t totptable) Init() tea.Cmd { return nil; }

func (t totptable) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  var cmd tea.Cmd;

  switch msg := msg.(type) {
    case tea.KeyMsg:
      switch msg.Type {
        case tea.KeyEsc, tea.KeyCtrlC:
          return t, tea.Quit

        case tea.KeyCtrlD:
          var acctDb gtoken.Database;
          // Read database contents from disk
          accountFile, err := common.GetAccountsDB();
          if err != nil {
            t.cliError = err;
            return t, nil;
          }
          acctDb, err = gtoken.OpenAccountDB(accountFile);
          if err != nil {
            t.cliError = err;
            return t, nil;
          }

          // delete entry from database
          acctDb.DeleteAccount(t.table.SelectedRow()[6]);
          acctDb.WriteAccountsDB(accountFile);
          return t, tea.Quit;
      }
  }

  t.table, cmd = t.table.Update(msg);
  return t, cmd;
}

func (t totptable) View() string {
  var b strings.Builder;

  b.WriteString(t.table.View());
  b.WriteString(fmt.Sprintf("\n\n%s - %s\n", "Use ESC to Quit", "Press Ctrl-D to delete the selected entry"));

  if t.cliError != nil {
    b.WriteString(fmt.Sprintf("%v", t.cliError));
  }

  return b.String();
}
