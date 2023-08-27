package account

import (
	"github.com/mcaimi/go-totp/rfc6238"
	"github.com/mcaimi/gtoken/pkg/token_io"
	"github.com/mcaimi/gtoken/pkg/database"
)

func GenerateTokens() ([]database.TokenEntity, error) {
  var acctDb token_io.Database;
  var rows []database.TokenEntity;

  acctDb, err := token_io.ReadAccountDb();
  if err != nil {
    return nil, err;
  }

  var a database.TokenEntity;
  var n int = acctDb.Entries;
  rows = make([]database.TokenEntity, n);
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
    rows[i] = database.TokenEntity{UUID: a.UUID,
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

