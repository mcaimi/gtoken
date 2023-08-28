package account

import (
	"github.com/mcaimi/gtoken/pkg/token_io"
	"github.com/mcaimi/gtoken/pkg/database"
)

func UpdateToken(uuid string, newToken database.TokenEntity) error {
  // update token
  if e := token_io.UpdateAccount(uuid, newToken); e != nil {
    return e;
  }

  return nil;
}

