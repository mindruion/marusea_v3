//package main
//
//import (
//	"errors"
//	"fmt"
//	"github.com/manifoldco/promptui"
//	"regexp"
//	"strconv"
//)
//
//func normalize(phoneNumber string) string {
//	re := regexp.MustCompile("\\D")
//	return re.ReplaceAllString(phoneNumber, "")
//}
//
//func main() {
//	items := []map[string]string{

//	branchCreateItems := []string{"Yes", "No"}
//
//	var result1,result2,result3,result4 string
//	var err error
//
//	validate := func(input string) error {
//		_, err := strconv.ParseFloat(input, 64)
//		if err != nil {
//			return errors.New("invalid number")
//		}
//		return nil
//	}
//	templates := &promptui.SelectTemplates{
//		Active:   "{{ . | blue }} ",
//		Inactive: "{{ . | black }} ",
//		Selected: "{{ . | green }} ",
//	}
//	promptTemplate := &promptui.PromptTemplates{
//			Prompt:  "{{ . }} ",
//			Valid:   "{{ . | green }} ",
//			Invalid: "{{ . | red }} ",
//			Success: "{{ . | bold }} ",
//		}
//
//	_select := promptui.Select{
//		Label:    "Write bellow description:",
//		Items:    items,
//		Templates: templates,
//	}
//
//	_prompt := promptui.Prompt{
//		Label:    "Write bellow description:",
//		Templates: promptTemplate,
//	}
//
//	_prompt2 := promptui.Prompt{
//		Label:    "What is task number:",
//		//Templates: prompt_templates,
//		Validate: validate,
//	}
//
//	_prompt3 := promptui.Select{
//		Label:    "Should I create a new branch?",
//		Items: branchCreateItems,
//		Templates: templates,
//	}
//
//	_, result1, err = _select.Run()
//	result2, err = _prompt.Run()
//	result3, err = _prompt2.Run()
//	_, result4, err = _prompt3.Run()
//
//	if err != nil {
//		fmt.Printf("Prompt failed %v\n", err)
//		return
//	}
//
//	fmt.Printf("You choose %s\n", result1)
//	fmt.Printf("You choose %s\n", result2)
//	fmt.Printf("You choose %s\n", result3)
//	fmt.Printf("You choose %s\n", result4)
//}

package main

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
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
		{Name: "Feat", Value: "feat"},
		{Name: "Fix", Value: "fix"},
		{Name: "Docs", Value: "docs"},
		{Name: "Chore", Value: "chore"},
		{Name: "Style", Value: "style"},
		{Name: "Refactor", Value: "refactor"},
		{Name: "Perf", Value: "perf"},
		{Name: "Test", Value: "test"},
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
		Validate: validateLength,
	}

	_prompt2 := promptui.Prompt{
		Label: "What is task number:",
		Templates: promptTemplate,
		Validate: validate,
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
	_, createNewBranch, err := _prompt3.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	description = strings.TrimSpace(description)
	taskNumber = strings.TrimSpace(taskNumber)
	action := options[i]

	if createNewBranch == "Yes" {
		res := strings.ReplaceAll(description, " ", "-")
		branchName := fmt.Sprintf("%s/%s-#%s", action, res, taskNumber)
		gitCommand("checkout", "-b", branchName)
	}

	gitCommand("add", ".")
	gitCommand("commit", "-m", fmt.Sprintf("%s: %s # %s", action, description, taskNumber))
	gitCommand("push", "origin", "HEAD")
}

func gitCommand(args...string){
	if c, err :=exec.Command("git", args...).CombinedOutput(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s\n", c)
	}
}
