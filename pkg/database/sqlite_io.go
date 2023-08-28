package database

import (
  "fmt"
  "errors"
  "crypto/sha1"
  "encoding/json"
  "database/sql"
  "github.com/mattn/go-sqlite3"
)

type SqliteDatabase struct {
  db *sql.DB;
}

// convenience function that interprets a sql error
func raise(e error) error {
  var innerError sqlite3.Error;
  if errors.As(e, &innerError) {
    s := fmt.Sprintf("Code: [%s] SystemErrno [%s] ExtendedCode [%s]", innerError.Code.Error(), innerError.SystemErrno.Error(), innerError.ExtendedCode.Error())
    return errors.New(s);
  } else {
    return e;
  }
}

// get a new DB descriptor
func NewDB(dbPath string) (*SqliteDatabase, error) {
  if fileDescriptor, err := sql.Open("sqlite3", dbPath); err != nil {
    return nil, raise(err);
  } else {
    return &SqliteDatabase{db: fileDescriptor}, nil;
  }
}

// count entries
func (d *SqliteDatabase) Count(table string) int {
  // count entries in the token table
  result := d.db.QueryRow(fmt.Sprintf(COUNT, table));

  var t_num int 
  result.Scan(&t_num); 

  return t_num;
}

func (d *SqliteDatabase) TokenCount() (int) {
  return d.Count(TOKEN_ENTITIES_TABLE);
}

// compute checksum of a serialized version of the token account table
func (d *SqliteDatabase) Checksum() (string, error) {
    // compute db checksum signature
    var err error;
    var t []TokenEntity;

    if  t, err = d.AllRows(); err != nil {
      return "", raise(err);
    }

    var serializedDbContents []byte;
    if serializedDbContents, err = json.Marshal(t); err != nil {
      return "", raise(err);
    }

    // return token db as a byte array
    return fmt.Sprintf("%x", sha1.Sum(serializedDbContents)), nil;
}

// perform database initial creation or schema migration
func (d *SqliteDatabase) DoInit() error {
  // init token and metadata tables
  if _, err := d.db.Exec(DB_INIT); err != nil {
    return raise(err);
  }
  if _, err := d.db.Exec(METADATA_INIT); err != nil {
    return raise(err);
  }

  // look for entries in the metadata table
  if d.Count(METADATA_TABLE) == 0 {
    // ok, first init needs to be performed
    fmt.Printf("Initializing Metadata Table...");

    if checksum, err := d.Checksum(); err != nil {
      return raise(err);
    } else {
      if _, err := d.db.Exec(METADATA_INSERT, SCHEMA_VERSION, d.TokenCount(), checksum); err != nil {
        return raise(err);
      }
      fmt.Printf("Done.\n");
    }
  } else {
    return errors.New("Already Initialized.");
  }

  // ok
  return nil;
}

// insert a new token
func (d *SqliteDatabase) InsertRow(t TokenEntity) error {
  if _, err := d.db.Exec(ROW_INSERT, t.UUID, t.Name, t.Email, t.Algorithm, t.Flavor, t.Interval, t.Type, t.Key, t.Period); err != nil {
    return raise(err);
  }

  // update Metadata...
  if checksum, err := d.Checksum(); err != nil {
    return raise(err);
  } else {
    if _, err := d.db.Exec(METADATA_UPDATE, SCHEMA_VERSION, d.TokenCount(), checksum, SCHEMA_VERSION); err != nil {
      return raise(err);
    }
  }

  // ok
  return nil;
}

// return an array with all token accounts
func (d *SqliteDatabase) AllRows() ([]TokenEntity, error) {
  // perform a full scan
  if rows, err := d.db.Query(SELECT_ALL); err != nil {
    return nil, raise(err);
  } else {
    defer rows.Close();

    // build token array
    var tokens []TokenEntity;
    for rows.Next() {
      var t TokenEntity;
      if err := rows.Scan(&t.UUID, &t.Name, &t.Email, &t.Algorithm, &t.Flavor, &t.Interval, &t.Type, &t.Key, &t.Period); err != nil {
        return nil, raise(err);
      }
      tokens = append(tokens, t);
    }

    // ok
    return tokens, nil;
  }
}

// read metadata
func (d *SqliteDatabase) ReadMetadata() (string, int, string, error) {
  // make a query
  var row *sql.Row;
  row = d.db.QueryRow(METADATA_ALL);

  // parse into a Token Struct
  var n_entries int;
  var version, checksum string;
  if err := row.Scan(&version, &n_entries, &checksum); err != nil {
    return "", 0,"", err;
  }

  // ok
  return version, n_entries, checksum, nil;
}

// validate integrity
func (d *SqliteDatabase) IntegrityCheck() (bool, error) {
  var metadata_checksum, computed_checksum string;
  var e error;

  // read metadata
  if _, _, metadata_checksum, e = d.ReadMetadata(); e != nil {
    return false, raise(e);
  }

  // compute checksum
  if computed_checksum, e = d.Checksum(); e != nil {
    return false, raise(e);
  }

  // validate
  return computed_checksum == metadata_checksum, nil;
}

// search for a named account by UUID
func (d *SqliteDatabase) SearchRow(uuid string) (TokenEntity, error) {
  // make a query
  var row *sql.Row;
  row = d.db.QueryRow(SEARCH_BY_UUID, uuid);

  // parse into a Token Struct
  var t TokenEntity;
  if err := row.Scan(&t.UUID, &t.Name, &t.Email, &t.Algorithm, &t.Flavor, &t.Interval, &t.Type, &t.Key, &t.Period); err != nil {
    return TokenEntity{}, raise(err);
  }

  // ok
  return t, nil;
}

// update a named token account
func (d *SqliteDatabase) UpdateRow(uuid string, t TokenEntity) error {
  // update a row
  if r, err := d.db.Exec(ROW_UPDATE, t.Name, t.Email, t.Algorithm, t.Flavor, t.Interval, t.Type, t.Key, t.Period, uuid); err != nil {
    return raise(err);
  } else {
    if n, err := r.RowsAffected(); err != nil {
      return err;
    } else {
      // update Metadata...
      if checksum, err := d.Checksum(); err != nil {
        return raise(err);
      } else {
        if _, err := d.db.Exec(METADATA_UPDATE, SCHEMA_VERSION, d.TokenCount(), checksum, SCHEMA_VERSION); err != nil {
          return raise(err);
        }
      }
      fmt.Printf("Updated %d Rows.\n", n);
      return nil;
    }
  }
}

// delete a named token account
func (d *SqliteDatabase) DeleteRow(UUID string) error {
  // delete a row
  if r, err := d.db.Exec(ROW_DELETE, UUID); err != nil {
    return raise(err);
  } else {
    if n, err := r.RowsAffected(); err != nil {
      return err;
    } else {
      // update Metadata...
      if checksum, err := d.Checksum(); err != nil {
        return raise(err);
      } else {
        if _, err := d.db.Exec(METADATA_UPDATE, SCHEMA_VERSION, d.TokenCount(), checksum, SCHEMA_VERSION); err != nil {
          return raise(err);
        }
      }
      fmt.Printf("Deleted %d Rows.\n", n);
      return nil;
    }
  }
}

// close the database connection
func (d *SqliteDatabase) CloseDB() {
  if e := d.db.Ping(); e == nil {
    d.db.Close();
  }
}
