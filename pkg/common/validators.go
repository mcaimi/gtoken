package common

import (
  "fmt"
  "regexp"
  "strconv"
  "strings"
)

const (
  pad int = 4
  inputCard int = 6;
  inputLength int = 32;
  emailRegex string = "^(\\w\\.?)+@[\\w\\.-]+\\.\\w{2,}$"
)

func StringNotZeroLen(s string) bool {
  if len(s) == 0 {
    return false;
  }

  return true;
}

func ValidateIntegerInput(s string) error {
  _, err := strconv.ParseInt(s, 10, inputLength);
  return err;
}

func ValidateEmailInput(s string) error { 
  emailMatcher := regexp.MustCompile(emailRegex);

  if ! emailMatcher.MatchString(s) {
    return fmt.Errorf("Input is not an email address");
  }

  return nil;
}

func ValidateAlgorithmInput(s string) error {
  var supportedAlgos []string = []string{ "md5", "sha1", "sha256", "sha512" }

  for i := range supportedAlgos {
    if strings.ToLower(s) == supportedAlgos[i] {
      return nil;
    }
  }

  return fmt.Errorf("Hashing Algorithm is not currently supported.");
}

func ValidateFlavor(s string) error {
  var supportedFlavors []string = []string{ "google", "rfc" }

  for i := range supportedFlavors {
    if strings.ToLower(s) == supportedFlavors[i] {
      return nil 
    }
  }

  return fmt.Errorf("Unsupported Authenticator");
}
