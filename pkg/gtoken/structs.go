package gtoken

type Account struct {
  Name string `json:"name"`;
  Email string `json:"email"`;
  Key string `json:"key"`;
  Hash string `json:"hash"`;
  Interval int `json:"interval"`;
  Flavor string `json:"flavor"`;
  Type string `json:"type"`;
  Uuid string `json:"uuid"`;

}

type Database struct {
  Accounts []Account `json:"accounts"`;
}
