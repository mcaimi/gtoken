package account

import (
	"fmt"

	"github.com/mcaimi/gtoken/pkg/common"
	"github.com/mcaimi/gtoken/pkg/token_io"
	"github.com/mcaimi/gtoken/pkg/database"
)

func ValidateToken(newToken database.TokenEntity) error {
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

