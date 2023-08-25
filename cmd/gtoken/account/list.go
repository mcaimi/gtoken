package account

import (
	"github.com/mcaimi/go-totp/rfc6238"
	"github.com/mcaimi/gtoken/pkg/token_io"
)

func GenerateTokens() ([]token_io.Account, error) {
  var acctDb token_io.Database;
  var rows []token_io.Account;

  acctDb, err := token_io.ReadAccountDb();
  if err != nil {
    return nil, err;
  }

  var a token_io.Account;
  var n int = acctDb.Count();
  rows = make([]token_io.Account, n);
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
    rows[i] = token_io.Account{UUID: a.UUID,
      Name: a.Name,
      Email: a.Email,
      Algorithm: a.Algorithm,
      Flavor: a.Flavor,
      Interval: a.Interval,
      Type: a.Type,
      Key: a.Key,
      Token: token};
  }

  return rows, nil;
}

