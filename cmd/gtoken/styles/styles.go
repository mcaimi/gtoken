package styles

import (
  "github.com/jedib0t/go-pretty/v6/table"
  "github.com/jedib0t/go-pretty/v6/text"
)

func GetStyle() table.Style {
  var tokenStyle table.Style = table.StyleLight;

  // customize default style
  tokenStyle.Name = "GToken Table Style";
  tokenStyle.Color.IndexColumn = text.Colors{text.BgBlack, text.FgHiBlue};
  tokenStyle.Color.Footer = text.Colors{text.BgBlack, text.FgHiBlue};
  tokenStyle.Color.Header = text.Colors{text.BgBlack, text.FgHiRed};
  tokenStyle.Color.Row = text.Colors{text.BgBlack, text.FgWhite};
  tokenStyle.Color.RowAlternate = text.Colors{text.BgBlack, text.FgHiWhite};

  // formats
  tokenStyle.Format.Row = text.FormatDefault;
  tokenStyle.Format.Header = text.FormatDefault;

  // return object
  return tokenStyle;
}

