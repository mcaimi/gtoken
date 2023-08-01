package gtoken

import (
	"github.com/mcaimi/gtoken/pkg/common"
	"github.com/mcaimi/gtoken/pkg/gtoken"
)

func RemoveToken(tokenUuid string) error {
  // Remove entry from the DB
  accountFile, err := common.GetAccountsDB();
  if err != nil {
    return err;
  }
  acctDb, err := gtoken.OpenAccountDB(accountFile);
  if err != nil {
    return err;
  }
  acctDb.DeleteAccount(tokenUuid);
  acctDb.WriteAccountsDB(accountFile);

  return nil;
}

