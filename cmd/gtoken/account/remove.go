package account

import (
  "github.com/mcaimi/gtoken/pkg/token_io"
)

func RemoveToken(tokenUuid string) error {
  // Remove entry from the DB
  if e := token_io.DeleteAccount(tokenUuid); e != nil {
    return e;
  }

  return nil;
}

