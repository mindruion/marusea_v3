package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
)

type commitOptions struct {
	Name  string
	Value string
}

func main() {
	var err error
	options := []commitOptions{
		{Name: "Feat", Value: ":star: Feat"},
		{Name: "Fix", Value: ":beetle: Fix"},
		{Name: "Docs", Value: ":books: Docs"},
		{Name: "Chore", Value: ":gear: Chore"},
		{Name: "Style", Value: ":rocket: Style"},
		{Name: "Refactor", Value: ":hammer: Refactor"},
		{Name: "Perf", Value: ":zap: Perf"},
		{Name: "Test", Value: ":test_tube: Test"},
		{Name: "Config", Value: ":hammer_and_wrench: Config"},
		{Name: "Renaming", Value: ":label: Renaming"},
		{Name: "Revert", Value: ":hedgehog: Revert"},
	}

	branchCreateItems := []commitOptions{
		{Name: "Yes", Value: "yes"},
		{Name: "No", Value: "no"},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "> {{ .Name | green }}",
		Inactive: "{{ .Name | cyan }}",
		Selected: "\U0001F336 {{ .Name | green | cyan }}",
	}
	promptTemplate := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | green }} ",
	}

	searcher := func(input string, index int) bool {
		co := options[index]
		name := strings.Replace(strings.ToLower(co.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}
	validate := func(input string) error {
		_, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return errors.New("invalid number")
		}
		return nil
	}
	validateLength := func(input string) error {
		if len(input) <= 3 {
			return errors.New("invalid length")
		}
		return nil
	}
	_select := promptui.Select{
		Label:     "Select one valid option",
		Items:     options,
		Templates: templates,
		Size:      10,
		Searcher:  searcher,
	}
	_prompt := promptui.Prompt{
		Label:     "Write bellow description:",
		Templates: promptTemplate,
		Validate:  validateLength,
	}

	_prompt2 := promptui.Prompt{
		Label:     "What is task number:",
		Templates: promptTemplate,
		Validate:  validate,
	}

	_prompt3 := promptui.Select{
		Label:     "Should I create a new branch",
		Items:     branchCreateItems,
		Templates: templates,
		Size:      10,
		Searcher:  searcher,
	}
	i, _, err := _select.Run()
	description, err := _prompt.Run()
	taskNumber, err := _prompt2.Run()
	v, _, err := _prompt3.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	description = strings.TrimSpace(description)
	taskNumber = strings.TrimSpace(taskNumber)
	action := options[i]
	createNewBranch := branchCreateItems[v]
	branchName := ""
	if createNewBranch.Value == "yes" {
		res := strings.ReplaceAll(description, " ", "-")
		branchName = fmt.Sprintf("%s/%s-#%s", strings.ToLower(action.Name), res, taskNumber)
		gitCommand("checkout", "-b", branchName)
	}

	gitCommand("add", ".")
	gitCommand("commit", "-m", fmt.Sprintf("%s: %s #%s", action.Value, description, taskNumber))
	gitCommand("push", "origin", "HEAD")
	cleaCommand()

	if createNewBranch.Value == "yes" {
		fmt.Printf("Changes was pushed successfully to '%s'\n", branchName)
	} else {
		fmt.Printf("Changes was pushed successfully\n")
	}
}

func gitCommand(args ...string) {
	if _, err := exec.Command("git", args...).CombinedOutput(); err != nil {
		cleaCommand()
		panic(fmt.Sprintf("Command 'git %s' failed with: \nError: %s", strings.Join(args, " "), err))
	}
}

func cleaCommand(){
	var clear map[string]func()     //create a map for storing clear funcs
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}

}
