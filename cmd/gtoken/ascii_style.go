package gtoken

import "github.com/charmbracelet/lipgloss"

var (
  // foreground colors - normal
  black = lipgloss.Color("0");
  red = lipgloss.Color("1"); 
  green = lipgloss.Color("2");
  yellow = lipgloss.Color("3");
  blue = lipgloss.Color("4");
  magenta = lipgloss.Color("5");
  cyan = lipgloss.Color("6");
  white = lipgloss.Color("7");

  // foreground colors - bright
  bright_black = lipgloss.Color("8");
  bright_red = lipgloss.Color("9"); 
  bright_green = lipgloss.Color("10");
  bright_yellow = lipgloss.Color("11");
  bright_blue = lipgloss.Color("12");
  bright_magenta = lipgloss.Color("13");
  bright_cyan = lipgloss.Color("14");
  bright_white = lipgloss.Color("15");

  // Styles for foreground and background text inputs
  textFgStyle = lipgloss.NewStyle().Foreground(white);
  textBgStyle = lipgloss.NewStyle().Background(black);
  
  // styles for input text fields
  inputSelectedStyle = lipgloss.NewStyle().Foreground(yellow);
  inputGrayedStyle = lipgloss.NewStyle().Foreground(blue);

  // common colors and styles
  foregroundStyle = lipgloss.NewStyle().Foreground(bright_white);
  backgroundStyle = lipgloss.NewStyle().Background(black);
  errorStyle = lipgloss.NewStyle().Foreground(red);
  normalBorder = lipgloss.NormalBorder();
)

