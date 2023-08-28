package account

import (
	"github.com/mcaimi/gtoken/pkg/token_io"
	"github.com/mcaimi/gtoken/pkg/database"
)

func InsertToken(newToken database.TokenEntity) error {
  // Build Account object
  var acct database.TokenEntity;
  acct = database.TokenEntity{
    Name: newToken.Name,
    Email: newToken.Email,
    Key: newToken.Token,
    Algorithm: newToken.Algorithm, 
    Interval: newToken.Interval,
    Flavor: newToken.Flavor,
    Type: newToken.Type,
  }

  // insert token
  if e := token_io.InsertAccount(acct); e != nil {
    return e;
  }

  return nil;
}

