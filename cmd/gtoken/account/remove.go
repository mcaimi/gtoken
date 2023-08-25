package account

import (
  "github.com/mcaimi/gtoken/pkg/token_io"
)

func RemoveToken(tokenUuid string) error {
  // Remove entry from the DB
  acctDb, err := token_io.ReadAccountDb();
  if err != nil {
    return err;
  }
  acctDb.DeleteAccount(tokenUuid);
  acctDb.WriteAccountsDB(acctDb.DbFilePath);

  return nil;
}

