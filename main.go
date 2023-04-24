package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

var url string

type model struct {
	status int
	err    error
}

func checkServer(url string) tea.Cmd {
	return func() tea.Msg {
		// create an HTTP client and make a GET request
		c := &http.Client{Timeout: 10 * time.Second}
		res, err := c.Get(url)

		if err != nil {
			// there was an error making our request
			// wrap the error we received in a message and return it
			return errMsg{err}
		}

		// we received a response from the server
		// return the HTTP status code as a message
		return statusMsg(res.StatusCode)
	}
}

type statusMsg int

type errMsg struct{ err error }

// for messages that contain errors it's often handy to also implement the error interface on the message
func (e errMsg) Error() string { return e.err.Error() }

func (m model) Init() tea.Cmd {
	if len(os.Args) <= 1 {
		fmt.Println("You must enter a URL to check")
		os.Exit(1)
	}
	url = os.Args[1]
	return checkServer(url)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case statusMsg:
		// the server returned a status message
		// save it to our model
		// exit the program because there is nothing left to do
		// the final view will still be rendered before exiting
		m.status = int(msg)
		return m, tea.Quit

	case errMsg:
		// there was an error, note it in the model and quit
		m.err = msg
		return m, tea.Quit

	case tea.KeyMsg:
		// allow user to quit using 'ctrl+c' or 'q' even though it should run quickly
		if msg.Type == tea.KeyCtrlC || msg.Type.String() == "q" {
			return m, tea.Quit
		}
	}

	// if any other messages somehow come through, do nothing
	return m, nil
}

func (m model) View() string {
	// if there's an error, print it and do nothing else
	if m.err != nil {
		return fmt.Sprintf("\n We had some trouble: %v\n\nn", m.err)
	}

	// tell the user we're doing something
	s := fmt.Sprintf("Checking %s... ", url)

	// when the server reponds with a status, add it to the output
	if m.status > 0 {
		s += fmt.Sprintf("%d %s!", m.status, http.StatusText(m.status))
	}

	// return the string for renndering
	return "\n" + s + "\n\n"
}

func main() {
	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Printf("Oh no, there was an error: %v\n", err)
		os.Exit(1)
	}
}
