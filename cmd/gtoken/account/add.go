package account

import (
	"fmt"

	"github.com/mcaimi/gtoken/pkg/common"
	"github.com/mcaimi/gtoken/pkg/gtoken"
)

func ValidateToken(newToken gtoken.Account) error {
  var err error;
  // validate input
  if err = common.ValidateEmailInput(newToken.Email); err != nil {
    return err;
  }
  if err = common.ValidateAlgorithmInput(newToken.Algorithm); err != nil {
    return err;
  }
  if err = common.ValidateFlavor(newToken.Flavor); err != nil {
    return err;
  }
  // validate string lengths
  inputStringOk := common.StringNotZeroLen(newToken.Token) && common.StringNotZeroLen(newToken.Name)
  if !inputStringOk {
    return fmt.Errorf("Token Seed and Account Name cannot be empty");
  }

  return nil;
}

func InsertToken(newToken gtoken.Account) error {
  // Build Account object
  var acct gtoken.Account;
  acct = gtoken.Account{
    Name: newToken.Name,
    Email: newToken.Email,
    Key: newToken.Token,
    Algorithm: newToken.Algorithm, 
    Interval: newToken.Interval,
    Flavor: newToken.Flavor,
    Type: newToken.Type,
  }

  // write account data to the database on disk
  acctDb, err := gtoken.ReadAccountDb();
  if err != nil {
    return err;
  }
  acctDb.InsertAccount(acct);
  acctDb.WriteAccountsDB(acctDb.DbFilePath);

  return nil;
}

