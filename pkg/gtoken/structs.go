package gtoken

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

type Account struct {
  Uuid string `json:"uuid"`;
  Name string `json:"name"`;
  Email string `json:"email"`;
  Algorithm string `json:"hash"`;
  Flavor string `json:"flavor"`;
  Interval int `json:"interval"`;
  Type string `json:"type"`;
  Key string `json:"key"`;
  Token string;
}

type Database struct {
  DbFilePath string;
  Accounts []Account `json:"accounts"`;
}
