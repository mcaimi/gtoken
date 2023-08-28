package token_io

import "github.com/mcaimi/gtoken/pkg/database"

const (
  totp_uuid int = iota;
  account_name;
  email;
  totp_algo;
  totp_flavor;
  totp_interval;
  totp_type;
  totp_key;
  totp_computed_value;
)

type Database struct {
  Version string;
  Entries int;
  IntegrityChecksum string;
  IsValid bool;
  DbFilePath string;
  Accounts []database.TokenEntity `json:"accounts"`;
}
