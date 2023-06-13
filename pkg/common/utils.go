package common

import (
  "os"
  "errors"
  "strings"
)

func GetConfigPath() (string, error) {
  var configDir string;
  var err error;

  configDir, err = os.UserConfigDir(); 
  if err != nil {
    return "", err;
  }

  return strings.Join([]string{ configDir, CONFIG_PATH_NAME }, "/"), nil;
}

func GetAccountsDB() (string, error) {
  var configPath string;
  var err error;

  configPath, err = GetConfigPath();
  if err != nil {
    configPath, err = os.UserHomeDir();
    if err != nil {
      return "", err;
    }
  }

  // sanity check: make sure config directory exists
  if _, err = os.Stat(configPath); err != nil && errors.Is(err, os.ErrNotExist) {
    // create directory
    if err = os.MkdirAll(configPath, os.ModePerm); err != nil {
      return "", err;
    }
  } 

  return strings.Join([]string{ configPath, ACCOUNTS_DB_NAME }, "/"), nil;
}

