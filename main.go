package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	todos     []todo
	cursor    int
	editField textinput.Model
	addField  textinput.Model
}

var fileName = "tpomo_todos.json"
var pomTime time.Duration = time.Second * 5

func initialModel() model {
	todos := getTodos(fileName)
	editField := textinput.New()
	editField.Prompt = "Editing: "
	addField := textinput.New()
	addField.Prompt = "Create New: "

	return model{
		todos:     todos,
		cursor:    0,
		editField: editField,
		addField:  addField,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if m.editField.Focused() {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "esc":
				m.editField.Blur()
				return m, nil

			case "enter":
				m.editField.Blur()
				m.todos[m.cursor].Name = m.editField.Value()
				saveTodos(fileName, m.todos)
				return m, nil
			}
		}

		m.editField, cmd = m.editField.Update(msg)
		return m, cmd
	} else if m.addField.Focused() {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "esc":
				m.addField.Blur()
				m.addField.SetValue("")
				return m, nil

			case "enter":
				m.addField.Blur()
				var newTodo = todo{
					Name: m.addField.Value(),
					Done: false,
				}
				m.todos = append(m.todos, newTodo)
				saveTodos(fileName, m.todos)
				m.addField.SetValue("")
				return m, nil
			}
		}

		m.addField, cmd = m.addField.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if len(m.todos) == 0 {
				break
			}
			if m.cursor == 0 {
				m.cursor = len(m.todos) - 1
				break
			}
			m.cursor--
		case "down", "j":
			if len(m.todos) == 0 {
				break
			}
			if len(m.todos)-1 == m.cursor {
				m.cursor = 0
				break
			}
			m.cursor++
		case "a":
			if m.addField.Focused() {
				m.addField.Blur()
				m.addField.Placeholder = ""
			} else {
				m.addField.Focus()
				m.addField.Placeholder = "Added new todo"
			}
		case "d":
			m.todos = append(m.todos[:m.cursor], m.todos[m.cursor+1:]...)
			saveTodos(fileName, m.todos)
		case " ":
			m.todos[m.cursor].Done = !m.todos[m.cursor].Done
			saveTodos(fileName, m.todos)
		case "e":
			if m.editField.Focused() {
				m.editField.Blur()
			} else {
				m.editField.Focus()
				m.editField.SetValue(m.todos[m.cursor].Name)
			}
		}

	}
	return m, nil
}

func (m model) View() string {
	s := "Your Todos:\n\n"
	help := "exit: ctrl-c, esc | confirm: submit | delete: d | edit: e | check: space"

	for i, todo := range m.todos {
		cursor := ""

		if m.cursor == i {
			cursor = "👈"
		}

		done := "🔲"
		if todo.Done {
			done = "✅"
		}

		s += fmt.Sprintf("%s %s %s\n", done, todo.Name, cursor)
	}

	edit := ""
	if m.editField.Focused() {
		edit = m.editField.View()
		help = "exit: ctrl-c, esc | confirm: submit"
	}

	add := ""
	if m.addField.Focused() {
		add = m.addField.View()
		help = "exit: ctrl-c, esc | confirm: submit"
	}

	timer := m.timer.View()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		s,
		edit,
		add,
		help,
		timer,
	)
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("There was an error.")
		os.Exit(1)
	}
}
