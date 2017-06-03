package robologger

import (
  "bufio"
  "fmt"
  "os"

  . "github.com/logrusorgru/aurora"
)

type ResponseType int
type PromptType int

const (
  EMPTY ResponseType = iota
  YES
  NO
  CANCEL
  ALL
)

const (
  STRING PromptType = iota
  YESNO
  YESNOCANCEL
  YESNOCANCELALL
)

type Reader interface {
  ParseResponse(response string) (ResponseType, error)
  ReadInput() string
}

type DefaultReader struct {}

type ReaderOptions struct {
  MessageWidth int
}

func NewDefaultReader(opts ReaderOptions) *DefaultReader {
  return &DefaultReader{}
}

// `ParseResponse` takes the input read from `ReadInput` and converts it into
// a `ResponseType`. This makes checking for a particular response easier.
func (r *DefaultReader) ParseResponse(response string) (ResponseType, error) {
  switch response {
  case "y", "Y", "yes", "Yes":
    return YES, nil
  case "n", "N", "no", "No":
    return NO, nil
  case "c", "C", "cancel", "Cancel":
    return CANCEL, nil
  case "":
    return EMPTY, nil
  default:
    return EMPTY, fmt.Errorf("Cannot parse response: %s", response)
  }
}

// `ReadInput` uses a bufio scanner instead of fmt.Scan functions in order to
// read a single line of text from the user. This may change in the future, but
// for now, it works out well.
func (r *DefaultReader) ReadInput() string {
  scanner := bufio.NewScanner(os.Stdin)
  scanner.Scan()
  return scanner.Text()
}

func (l *Logger) prompt(t PromptType, format *string, args ...interface{}) string {
  // Create a new message.
  msg := NewMessage(PRINT, format, args...)

  // Format the message as a prompt.
  l.Writer.FormatPrompt(t, msg)

  // Write the message.
  l.Writer.WriteMessage(msg)

  // Get the input from the user.
  response := l.Reader.ReadInput()
  term.SaveCursorPosition()
  term.MoveToBeginning()
  term.MoveUp(1)

  term.Clear()

  l.Writer.WriteMessage(msg)

  fmt.Printf("%s", Cyan(response))
  fmt.Print("\n")

  term.RestoreCursorPosition()

  return response
}

// These two functions are part of the `Logger` interface. We include them here
// because they have special functions which relate closely to the `Reader`
// interface.
func (l *Logger) Prompt(t PromptType, args ...interface{}) string {
  return l.prompt(t, nil, args...)
}

func (l *Logger) Promptf(t PromptType, format string, args ...interface{}) string {
  return l.prompt(t, &format, args...)
}

// func (r *DefaultReader) ReadInputMultiline() (lines []string) {
//   scanner := bufio.NewScanner(os.Stdin)
//   for scanner.Scan() {
//     lines = append(lines, scanner.Text())
//   }
//   return
// }

func Prompt(t PromptType, args ...interface{}) string {
  return std.Prompt(t, args...)
}

func Promptf(t PromptType, format string, args ...interface{}) string {
  return std.Promptf(t, format, args...)
}

// These last two declarations are utility functions for working with the
// reader. If a user wants to read input for their own purposes, for example,
// they can interact with these two functions from the std logger.
func ParseResponse(response string) (ResponseType, error) {
  return std.Reader.ParseResponse(response)
}

func ReadInput() string {
  return std.Reader.ReadInput()
}
