package gtoken

import (
  "fmt"
  "strings"
  "strconv"

  "github.com/charmbracelet/bubbletea"
  "github.com/charmbracelet/bubbles/textinput"

  "github.com/mcaimi/gtoken/pkg/gtoken"
  "github.com/mcaimi/gtoken/pkg/common"
)

type newAccountModel struct {
  focusIndex int
  inputs     []textinput.Model
  cursorMode textinput.CursorMode
  cliError   error
}

func AddAccountModel() newAccountModel {
  m := newAccountModel{
    inputs: make([]textinput.Model, inputCard),
  }

  m.cursorMode = textinput.CursorBlink;
  m.cliError = nil;

  var t textinput.Model
  for i := range m.inputs {
    t = textinput.New()
    t.CharLimit = inputLength;

    switch i {
      case account_name:
        t.Placeholder = "John Appleseed"
        t.PromptStyle = inputSelectedStyle;
        t.TextStyle = foregroundStyle;
        t.Prompt = "Account Name > ";
        t.CharLimit = inputLength;
        t.Width = t.CharLimit + pad;
        t.Focus()
      case email:
        t.Placeholder = "appleseed@email.tld"
        t.PromptStyle = inputGrayedStyle;
        t.TextStyle = foregroundStyle;
        t.Prompt = "E-Mail Address > ";
        t.CharLimit = 2 * inputLength;
        t.Width = t.CharLimit + pad;
      case totp_key:
        t.Placeholder = "TOTP Key"
        t.PromptStyle = inputGrayedStyle;
        t.TextStyle = foregroundStyle;
        t.Prompt = "TOTP Seed Key > ";
        t.CharLimit = 2 * inputLength;
        t.Width = t.CharLimit + pad;
      case totp_algo:
        t.Placeholder = "SHA1"
        t.PromptStyle = inputGrayedStyle;
        t.TextStyle = foregroundStyle;
        t.Prompt = "Hashing Algo > ";
        t.CharLimit = 2 * inputLength;
        t.Width = t.CharLimit + pad;
      case totp_interval:
        t.Placeholder = "180"
        t.PromptStyle = inputGrayedStyle;
        t.TextStyle = foregroundStyle;
        t.Prompt = "TOTP Validity > ";
        t.CharLimit = 3
        t.Width = t.CharLimit + pad;
        t.Validate = validateIntegerInput
      case totp_flavor:
        t.SetValue("google");
        t.PromptStyle = inputGrayedStyle;
        t.TextStyle = foregroundStyle;
        t.Prompt = "Authenticator > ";
        t.CharLimit = 2 * inputLength;
        t.Width = t.CharLimit + pad;
      }

    m.inputs[i] = t
    m.inputs[i].SetCursorMode(m.cursorMode);
  }

  return m
}

func (m newAccountModel) Init() tea.Cmd {
  return textinput.Blink
}

func (m newAccountModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  cmds := make([]tea.Cmd, len(m.inputs))

  switch msg := msg.(type) {
    case error:
      m.cliError = msg;
      return m, nil
      
    case tea.KeyMsg:
      switch msg.Type {
        case tea.KeyCtrlC, tea.KeyEsc:
          return m, tea.Quit

        // Set focus to next input
        case tea.KeyTab, tea.KeyShiftTab:
          // check for a valid email input
          if m.focusIndex == email {
            if validateEmailInput(m.inputs[email].Value()) != nil {
              m.inputs[email].TextStyle = errorStyle;
              m.inputs[email].SetValue("");
              m.inputs[email].Placeholder = "Invalid Email Format";
            } else {
              m.inputs[email].TextStyle = foregroundStyle;
            }
          }

          if m.focusIndex == totp_algo {
            if validateAlgorithmInput(m.inputs[totp_algo].Value()) != nil {
              m.inputs[totp_algo].TextStyle = errorStyle;
              m.inputs[totp_algo].SetValue("");
              m.inputs[totp_algo].Placeholder = "Invalid Algorithm";
            }
          } else {
            m.inputs[totp_algo].TextStyle = foregroundStyle;
          }

          if m.focusIndex == totp_flavor {
            if err := validateFlavor(m.inputs[totp_flavor].Value()); err != nil {
              m.inputs[totp_flavor].TextStyle = errorStyle;
              m.inputs[totp_flavor].SetValue("google");
              m.cliError = err;
            }
          } else {
            m.inputs[totp_flavor].TextStyle = foregroundStyle;
          }

          if msg.Type == tea.KeyTab {
            m.focusIndex++;
          } else if msg.Type == tea.KeyShiftTab {
            m.focusIndex--;
          }

          if m.focusIndex > len(m.inputs) {
            m.focusIndex = 0
          } else if m.focusIndex < 0 {
            m.focusIndex = len(m.inputs)
          }

        case tea.KeyEnter:
          if m.focusIndex == len(m.inputs) {
            // sanity check: all input fields must be populated
            for i := range m.inputs {
              if m.inputs[i].Value() == "" {
                m.cliError = fmt.Errorf("Please fill in all input fields.");
                return m, nil;
              }
            }

            // Build Account object
            var acct gtoken.Account;
            intv, err := strconv.Atoi(m.inputs[totp_interval].Value());
            if err == nil {
              acct = gtoken.Account{
                Name: m.inputs[account_name].Value(),
                Email: m.inputs[email].Value(),
                Key: m.inputs[totp_key].Value(),
                Hash: m.inputs[totp_algo].Value(), 
                Interval: intv,
                Flavor: m.inputs[totp_flavor].Value(),
                Type: "totp",
              }
            } else {
              m.cliError = err;
              return m, nil;
            }

            // write account data to the database on disk
            accountFile, err := common.GetAccountsDB();
            if err != nil {
              m.cliError = err;
              return m, nil;
            }
            acctDb, err := gtoken.OpenAccountDB(accountFile);
            if err != nil {
              m.cliError = err;
              return m, nil;
            }
            acctDb.InsertAccount(acct);
            acctDb.WriteAccountsDB(accountFile);

            return m, tea.Quit
          }
      }
  }

  for i := range m.inputs {
    if i == m.focusIndex {
      // Update widget Style
      m.inputs[i].PromptStyle = inputSelectedStyle;
      m.inputs[i].TextStyle = inputSelectedStyle;
      // Set focused state
      m.inputs[i].Focus()
      continue
    }

    // Remove focused state
    m.inputs[i].PromptStyle = inputGrayedStyle;
    m.inputs[i].TextStyle = inputGrayedStyle;
    m.inputs[i].Blur()
  }

  // Handle character input and blinking
  for i := range m.inputs {
    m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
  }

  return m, tea.Batch(cmds...);
}

func (m newAccountModel) View() string {
  var b strings.Builder

  header := foregroundStyle.Render("+ Insert a New TOTP Account. Please fill in all required information:\n");

  b.WriteString(header + "\n");

  b.WriteString(fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n",
    m.inputs[account_name].View(), m.inputs[email].View(),
    m.inputs[totp_key].View(), m.inputs[totp_algo].View(),
    m.inputs[totp_interval].View(), m.inputs[totp_flavor].View()));

  if m.focusIndex == len(m.inputs) {
    b.WriteString(inputSelectedStyle.Render("\n[ Save to disk. ]\n"));
  } else {
    b.WriteString(foregroundStyle.Render("\nSave to disk.\n"));
  }

  // display error if needed
  if m.cliError != nil {
    b.WriteString(errorStyle.Render(fmt.Sprintf("Failure: %v\n", m.cliError)));
  }

  return b.String()
}

