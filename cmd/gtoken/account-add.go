package gtoken

import (
	"fmt"

	"github.com/mcaimi/gtoken/pkg/common"
	"github.com/mcaimi/gtoken/pkg/gtoken"
)

func ValidateToken(newToken tokenObject) error {
  var err error;
  // validate input
  if err = validateEmailInput(newToken.email); err != nil {
    return err;
  }
  if err = validateAlgorithmInput(newToken.totp_algo); err != nil {
    return err;
  }
  if err = validateFlavor(newToken.totp_flavor); err != nil {
    return err;
  }
  // validate string lengths
  inputStringOk := stringNotZeroLen(newToken.token) && stringNotZeroLen(newToken.account_name)
  if !inputStringOk {
    return fmt.Errorf("Token Seed and Account Name cannot be empty");
  }

  return nil;
}

func InsertToken(newToken tokenObject) error {
  // Build Account object
  var acct gtoken.Account;
  acct = gtoken.Account{
    Name: newToken.account_name,
    Email: newToken.email,
    Key: newToken.token,
    Hash: newToken.totp_algo, 
    Interval: newToken.totp_interval,
    Flavor: newToken.totp_flavor,
    Type: newToken.totp_type,
  }

  // write account data to the database on disk
  accountFile, err := common.GetAccountsDB();
  if err != nil {
    return err;
  }
  acctDb, err := gtoken.OpenAccountDB(accountFile);
  if err != nil {
    return err;
  }
  acctDb.InsertAccount(acct);
  acctDb.WriteAccountsDB(accountFile);

  return nil;
}

