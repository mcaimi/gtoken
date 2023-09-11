package token_io

import (
	"encoding/json"
	"os"
  "fmt"
  "errors"

	"github.com/google/uuid"
	"github.com/mcaimi/gtoken/pkg/common"
  "github.com/mcaimi/gtoken/pkg/database"
)

func JsonLoad(filename string) error {
  var accountBackupDescriptor *os.File;
  var err error;

  accountBackupDescriptor, err = os.OpenFile(filename, os.O_RDONLY, 0600);
  if err != nil {
    return err;
  }
  defer accountBackupDescriptor.Close();

  // load accounts from json file
  var fileContents []byte;
  var fileInfo os.FileInfo;
  fileInfo, err = accountBackupDescriptor.Stat();
  if err != nil {
    return err;
  }
  fileContents = make([]byte, fileInfo.Size());

  var accountsDb Database;
  _, err = accountBackupDescriptor.Read(fileContents);

  // unmarshal results
  err = json.Unmarshal(fileContents, &accountsDb.Accounts);
  if err != nil {
    return nil;
  }

  // import data into the database
  // get default db path
  var dbPath string;
  var db *database.SqliteDatabase;
  if dbPath, err = common.GetAccountsDB(); err == nil {
    // open and return database
    if db, err = database.NewDB(dbPath); err == nil {
      defer db.CloseDB();
      // load entries
      for entry := range accountsDb.Accounts {
        fmt.Printf("Importing entry [%s]...", accountsDb.Accounts[entry].UUID);
        if err = db.InsertRow(accountsDb.Accounts[entry]); err != nil {
          return err;
        }
        fmt.Printf("..OK\n");
      }
    }
  }

  // return data
  return nil;
}

func JsonDump(fileName string) error {
  var accountDescriptor *os.File;
  var err error;

  accountDescriptor, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0600);
  if err != nil {
    return err;
  }
  defer accountDescriptor.Close();

  var accountsDb Database;
  if accountsDb, err = ReadAccountDb(); err != nil {
    return err;
  }

  // dump accounts to file
  var marshaledContent []byte;
  if marshaledContent, err = json.Marshal(accountsDb.Accounts); err != nil {
    return err;
  }

  // write to file
  accountDescriptor.Write(marshaledContent);

  // return
  return nil;
}

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
  inputStringOk := common.StringNotZeroLen(newToken.Key) && common.StringNotZeroLen(newToken.Name)
  if !inputStringOk {
    return fmt.Errorf("Token Seed and Account Name cannot be empty");
  }

  return nil;
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
      // initialize DB
      if e = db.DoInit(); e != nil {
        return e;
      }
      return nil;
    }
  }
  return nil;
}

// recompute checksum 
func ValidateDatabase() error {
  var dbPath string;
  var db *database.SqliteDatabase;
  var e error;
  var tokenArray Database;
  // get default db path
  if dbPath, e = common.GetAccountsDB(); e == nil {
    // open and return database
    if db, e = database.NewDB(dbPath); e == nil {
      defer db.CloseDB();
      // load entries
      tokenArray.DbFilePath = dbPath;
      if tokenArray.Accounts, e = db.AllRows(); e != nil {
        return e;
      }
      // Read Metadata
      if tokenArray.Version, tokenArray.Entries, tokenArray.IntegrityChecksum, e = db.ReadMetadata(); e != nil {
        return e;
      }

      // validate entries
      for i := range tokenArray.Accounts {
        fmt.Printf("Validating entry [%s] ...", tokenArray.Accounts[i].UUID);
        if e = ValidateToken(tokenArray.Accounts[i]); e != nil {
          return errors.New(fmt.Sprintf("Malformed or Corrputed entry: UUID %s", tokenArray.Accounts[i].UUID));
        }
        fmt.Printf("OK\n");
      }

      // update checksum
      fmt.Printf("%s\n", "Updating Checksum...");
      if e = db.UpdateChecksum(); e != nil {
        return e;
      }
      fmt.Printf("Done.\n");
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

func InsertAccountByFields(name string, email string, key string, hash string, interval int64, acct_type string, token_len int64) error {
  account_uuid := uuid.New().String();
  var acct database.TokenEntity = database.TokenEntity{ Name: name, Email: email, Key: key, Algorithm: hash, Interval: interval, Type: acct_type, UUID: account_uuid, Length: token_len }

  // insert account
  if e := InsertAccount(acct); e != nil {
    return e;
  }

  return nil;
}

func UpdateAccount(accountUuid string, acct database.TokenEntity) error {
  var dbPath string;
  var db *database.SqliteDatabase;
  var e error;

  // get default db path
  if dbPath, e = common.GetAccountsDB(); e == nil {
    // open and return database
    if db, e = database.NewDB(dbPath); e == nil {
      defer db.CloseDB();
      // search entry
      var updatedAccount database.TokenEntity;
      if updatedAccount, e = db.SearchRow(accountUuid); e != nil {
        return e;
      }

      // patch values
      updatedAccount.Name = common.Ternary(acct.Name != "", acct.Name, updatedAccount.Name);
      updatedAccount.Email = common.Ternary(acct.Email != "", acct.Email, updatedAccount.Email);
      updatedAccount.Algorithm = common.Ternary(acct.Algorithm != "", acct.Algorithm, updatedAccount.Algorithm);
      updatedAccount.Type = common.Ternary(acct.Type != "", acct.Type, updatedAccount.Type);
      updatedAccount.Flavor = common.Ternary(acct.Flavor != "", acct.Flavor, updatedAccount.Flavor);
      updatedAccount.Period = common.Ternary(acct.Period != 0, acct.Period, updatedAccount.Period);
      updatedAccount.Interval = common.Ternary(acct.Interval != 0, acct.Interval, updatedAccount.Interval);
      updatedAccount.Length = common.Ternary(acct.Length != 0, acct.Length, updatedAccount.Length);
      updatedAccount.Key = common.Ternary(acct.Key != "", acct.Key, updatedAccount.Key);

      // Update row
      db.UpdateRow(updatedAccount.UUID, updatedAccount)
    }

    return e;
  }
  return e;
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
        return Database{}, errors.New("Database integrity check failed. Database may be corrupted.");
      } else {
        tokenArray.IsValid = true;
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

