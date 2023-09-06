package database

// Single Token Generation Account Entity
type TokenEntity struct {
  UUID string `json:"uuid"`;
  Name string `json:"name"`;
  Email string `json:"email"`;
  Algorithm string `json:"hash"`;
  Flavor string `json:"flavor"`;
  Interval int64 `json:"interval"`;
  Type string `json:"type"`;
  Key string `json:"key"`;
  Period int64 `json:"period"`;
  Length int64 `json:"token_length"`;
  Token string;
}

// database metadata
type DatabaseMetadataEntity struct {
  Version string;
  Entries int64;
  Checksum string;
}
