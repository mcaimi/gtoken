package account

import (
  "fmt"
  "net/url"
  "errors"
	"github.com/mcaimi/gtoken/pkg/token_io"
	"github.com/mcaimi/gtoken/pkg/database"
)

func LoadTokens() ([]database.TokenEntity, error) {
  var acctDb token_io.Database;
  var rows []database.TokenEntity;

  acctDb, err := token_io.ReadAccountDb();
  if err != nil {
    return nil, err;
  }

  var a database.TokenEntity;
  var n int = acctDb.Entries;
  rows = make([]database.TokenEntity, n);
  for i := range acctDb.Accounts {
    a = acctDb.Accounts[i];

    // update table data
    rows[i] = database.TokenEntity{UUID: a.UUID,
      Name: a.Name,
      Email: a.Email,
      Algorithm: a.Algorithm,
      Flavor: a.Flavor,
      Interval: a.Interval,
      Type: a.Type,
      Key: a.Key,
      Length: a.Length};
  }

  return rows, nil;
}

func LoadToken(uuid string) (database.TokenEntity, error) {
  var acctDb token_io.Database;

  acctDb, err := token_io.ReadAccountDb();
  if err != nil {
    return database.TokenEntity{}, err;
  }

  var a database.TokenEntity;
  for i := range acctDb.Accounts {
    a = acctDb.Accounts[i];
    if a.UUID == uuid {
      return a, nil;
    }
  }

  return database.TokenEntity{}, errors.New("Entry Not Found");
}

func OtpUrl(k database.TokenEntity) string {
  urlTemplate := "otpauth://%s/%s?%s";
  labelTemplate := "%s:%s";

  var parmsString string;
  var labelString string;

  labelString = fmt.Sprintf(labelTemplate, k.Name, k.Email);
  labelString = url.QueryEscape(labelString);

  if k.Type == "totp" {
    parmsTemplate := "secret=%s&issuer=%s&digits=%d&period=%s";

    // fill in parameters
    parmsString = fmt.Sprintf(parmsTemplate, k.Key, k.Name, k.Length, k.Interval);
    parmsString = url.QueryEscape(parmsString);
  }

  // render template
  renderedUrl := fmt.Sprintf(urlTemplate, k.Type, labelString, parmsString);

  return renderedUrl;
}
