package account

import (
  "github.com/mcaimi/gtoken/pkg/gtoken"
)

func RemoveToken(tokenUuid string) error {
  // Remove entry from the DB
  acctDb, err := gtoken.ReadAccountDb();
  if err != nil {
    return err;
  }
  acctDb.DeleteAccount(tokenUuid);
  acctDb.WriteAccountsDB(acctDb.DbFilePath);

  return nil;
}

