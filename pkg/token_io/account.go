package token_io

import (
	"encoding/json"
	"os"
  "errors"

	"github.com/google/uuid"
	"github.com/mcaimi/gtoken/pkg/common"
  "github.com/mcaimi/gtoken/pkg/database"
)

func JsonOpen(fileName string) (Database, error) {
  var accountDescriptor *os.File;
  var err error;

  accountDescriptor, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0700);
  if err != nil {
    return Database{}, err;
  }
  defer accountDescriptor.Close();

  // load accounts from json file
  var fileContents []byte;
  var fileInfo os.FileInfo;
  fileInfo, err = accountDescriptor.Stat();
  if err != nil {
    return Database{}, err;
  }
  fileContents = make([]byte, fileInfo.Size());

  var accountsDb Database;
  _, err = accountDescriptor.Read(fileContents);

  // unmarshal results
  err = json.Unmarshal(fileContents, &accountsDb);
  if err != nil {
    return Database{}, nil;
  }

  // return data
  return accountsDb, nil;
}

func InitDB() error {
  var dbPath string;
  var db *database.SqliteDatabase;
  var e error;

  // get default db path
  if dbPath, e = common.GetAccountsDB(); e == nil {
    // open and return database
    if db, e = database.NewDB(dbPath); e == nil {
      defer db.CloseDB();
      // search entries
      if e = db.DoInit(); e != nil {
        return e;
      }
      // return database object
      return nil;
    }
  }
  return nil;
}

func InsertAccount(acct database.TokenEntity) error {
  // compute new UUID
  acct.UUID = uuid.New().String();
  var dbPath string;
  var db *database.SqliteDatabase;
  var e error;

  // get default db path
  if dbPath, e = common.GetAccountsDB(); e == nil {
    // open and return database
    if db, e = database.NewDB(dbPath); e == nil {
      defer db.CloseDB();
      // search entries
      if e = db.InsertRow(acct); e != nil {
        return e;
      }
      // return database object
      return nil;
    }
  }
  return nil;
}

func InsertAccountByFields(name string, email string, key string, hash string, interval int64, acct_type string) error {
  account_uuid := uuid.New().String();
  var acct database.TokenEntity = database.TokenEntity{ Name: name, Email: email, Key: key, Algorithm: hash, Interval: interval, Type: acct_type, UUID: account_uuid }

  // insert account
  if e := InsertAccount(acct); e != nil {
    return e;
  }

  return nil;
}

func DeleteAccount(accountUuid string) error {
  var dbPath string;
  var db *database.SqliteDatabase;
  var e error;

  // get default db path
  if dbPath, e = common.GetAccountsDB(); e == nil {
    // open and return database
    if db, e = database.NewDB(dbPath); e == nil {
      defer db.CloseDB();
      // search entries
      db.DeleteRow(accountUuid)
    }
    return e;
  }
  return e;
}

func SearchAccount(accountUuid string) *database.TokenEntity {
  var dbPath string;
  var db *database.SqliteDatabase;
  var t database.TokenEntity;
  var e error;

  // get default db path
  if dbPath, e = common.GetAccountsDB(); e == nil {
    // open and return database
    if db, e = database.NewDB(dbPath); e == nil {
      defer db.CloseDB();
      // search entries
      if t, e = db.SearchRow(accountUuid); e != nil {
        return nil;
      }
      // return database object
      return &t;
    } else {
      return nil;
    }
  } else {
    return nil;
  }
}

// convenience function: open and returns the local default database
func ReadAccountDb() (Database, error) {
  var dbPath string;
  var db *database.SqliteDatabase;
  var tokenArray Database;
  var e error;

  // get default db path
  if dbPath, e = common.GetAccountsDB(); e == nil {
    // open and return database
    if db, e = database.NewDB(dbPath); e == nil {
      defer db.CloseDB();
      tokenArray.DbFilePath = dbPath;
      // load entries
      if tokenArray.Accounts, e = db.AllRows(); e != nil {
        return Database{}, e;
      }

      // Read Metadata
      if tokenArray.Version, tokenArray.Entries, tokenArray.IntegrityChecksum, e = db.ReadMetadata(); e != nil {
        return Database{}, e;
      }

      // Validate integrity
      var isValid bool;
      if isValid, e = db.IntegrityCheck(); e != nil {
        return Database{}, e;
      }

      if !isValid {
        return Database{}, errors.New("Database integrity check failed. Db may be corrupted\n");
      }

      // return database object
      return tokenArray, nil;
    } else {
      return Database{}, e;
    }
  } else {
    return Database{}, e;
  }
}

