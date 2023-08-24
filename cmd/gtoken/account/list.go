package account

import (
  "github.com/mcaimi/gtoken/pkg/gtoken"
  "github.com/mcaimi/go-totp/rfc6238"
)

func GenerateTokens() ([]gtoken.Account, error) {
  var acctDb gtoken.Database;
  var rows []gtoken.Account;

  acctDb, err := gtoken.ReadAccountDb();
  if err != nil {
    return nil, err;
  }

  var a gtoken.Account;
  var n int = acctDb.Count();
  rows = make([]gtoken.Account, n);
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
    rows[i] = gtoken.Account{a.Uuid, a.Name, a.Email, a.Algorithm, a.Flavor, a.Interval, a.Type, a.Key, token};
  }

  return rows, nil;
}

