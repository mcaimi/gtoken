package gtoken

import (
  "github.com/mcaimi/gtoken/pkg/gtoken"
  "github.com/mcaimi/gtoken/pkg/common"

  "github.com/mcaimi/go-totp/rfc6238"
)

func GenerateTokens() ([]tokenObject, error) {
  var acctDb gtoken.Database;
  var rows []tokenObject;

  // Read database contents from disk
  accountFile, err := common.GetAccountsDB();
  if err != nil {
    return nil, err;
  }
  acctDb, err = gtoken.OpenAccountDB(accountFile);
  if err != nil {
    return nil, err;
  }

  var a gtoken.Account;
  var n int = acctDb.Count();
  rows = make([]tokenObject, n);
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
    rows[i] = tokenObject{a.Uuid, a.Name, a.Email, a.Hash, a.Flavor, a.Interval, a.Type, a.Uuid, token};
  }

  return rows, nil;
}

