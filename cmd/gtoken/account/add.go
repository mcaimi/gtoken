package account

import (
	"github.com/mcaimi/gtoken/pkg/token_io"
	"github.com/mcaimi/gtoken/pkg/database"
)

func InsertToken(newToken database.TokenEntity) error {
  // validate
  if err := token_io.ValidateToken(newToken); err != nil {
    return err;
  } else {
    // insert token
    if e := token_io.InsertAccount(newToken); e != nil {
      return e;
    }
  }
  return nil;
}

