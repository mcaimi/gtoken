package token_io

import (
	"encoding/json"
	"os"

	"github.com/google/uuid"
	"github.com/mcaimi/gtoken/pkg/common"
)

func OpenAccountDB(fileName string) (Database, error) {
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

func (d *Database) Count() int {
  var acctNum int = 0;
  for range d.Accounts {
    acctNum += 1;
  }

  return acctNum;
}

func (d *Database) InsertAccountByFields(name string, email string, key string, hash string, interval int64, acct_type string) {
  account_uuid := uuid.New().String();
  var acct Account = Account{ Name: name, Email: email, Key: key, Algorithm: hash, Interval: interval, Type: acct_type, UUID: account_uuid }

  d.Accounts = append(d.Accounts, acct);
}

func (d *Database) InsertAccount(acct Account) {
  acct.UUID = uuid.New().String();
  d.Accounts = append(d.Accounts, acct);
}

func (d *Database) DeleteAccount(accountUuid string) {
  // make sure that account exists
  var acct *Account = d.SearchAccount(accountUuid);
  if acct != nil {
    var updatedDb []Account;
    for i := range(d.Accounts) {
      if d.Accounts[i].UUID != acct.UUID {
        updatedDb = append(updatedDb, d.Accounts[i]);
      }
    }
    // swap databases
    d.Accounts = updatedDb;
  }
}

func (d *Database) SearchAccount(accountUuid string) *Account {
  if d.Count() == 0 {
    // empty db
    return nil;
  }

  // search for account in the database
  for i := range(d.Accounts) {
    if d.Accounts[i].UUID == accountUuid {
      return &d.Accounts[i]
    }
  }

  // account is not in the DB
  return nil
}

func (d *Database) WriteAccountsDB(fileName string) error {
  var accountDescriptor *os.File;
  var err error;

  // open new file
  accountDescriptor, err = os.OpenFile(fileName, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0700);
  if err != nil {
    return err;
  }
  defer accountDescriptor.Close();

  // marshal json db and write to the disk
  var dbContents []byte;
  dbContents, err = json.Marshal(d);
  if err != nil {
    return err;
  }

  accountDescriptor.Seek(0,0);
  accountDescriptor.Write(dbContents);
  accountDescriptor.Sync();

  return nil;
}

// convenience function: open and returns the local default database
func ReadAccountDb() (Database, error) {
  var dbPath string;
  var db Database;
  var e error;

  // get default db path
  if dbPath, e = common.GetAccountsDB(); e == nil {
    // open and return database
    if db, e = OpenAccountDB(dbPath); e == nil {
      // set db file name
      db.DbFilePath = dbPath;
      // return database object
      return db, nil;
    } else {
      return Database{}, e;
    }
  } else {
    return Database{}, e;
  }
}
