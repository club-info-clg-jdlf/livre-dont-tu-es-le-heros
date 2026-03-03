package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"time"
)

const (
	historyName   = "history.json"
	firstStepName = "intro"
)

// Thanks to https://mholt.github.io/json-to-go/
type history struct {
	MsPerChar time.Duration `json:"msPerChar"`
	Steps     steps         `json:"steps"`
}
type option struct {
	Goto   string `json:"goto"`
	Text   string `json:"text"`
	Hidden bool   `json:"hidden,omitempty"`
	Number int    `json:"number,omitempty"`
}
type options []option

type step struct {
	Title   string  `json:"title"`
	Text    string  `json:"text,omitempty"`
	Options options `json:"options,omitempty"`
}

type steps []step

func (steps *steps) getStep(name string) (step, error) {
	for _, step := range *steps {
		if step.Title == name {
			return step, nil
		}
	}

	return step{}, fmt.Errorf("(*main.steps).getStep(): step %s not found", strconv.Quote(name))
}

func (options *options) addNumbers()

func printWithTime(t time.Duration, a ...any) {
	s := fmt.Sprint(a...)

	for _, r := range s {
		fmt.Print(r)
		time.Sleep(t)
	}
}

func makeLinkWithText(url string, text string) string {
	return fmt.Sprintf(
		"\033]8;;%s\033\\%s\033]8;;\033\\",
		url,
		text,
	)
}

func makeLink(url string) string {
	return makeLinkWithText(url, url)
}

func doesntContain[T comparable](s []T, e T) bool {
	return !slices.Contains(s, e)
}

func main() {

	execPath, err := os.Executable()
	if err != nil {
		log.Fatalf("main.main(): Error finding the executable path: %s", err.Error())
	}

	execDirNotParsed := filepath.Dir(execPath) // Maybe it contains unresolved symbolic links

	execDir, err := filepath.EvalSymlinks(execDirNotParsed)
	if err != nil {
		log.Fatalf(
			"main.main(): Error evaluating symbolic links of the executable path (%s): %s",
			execDirNotParsed,
			err.Error(),
		)
	}

	relativeHistoryPath := filepath.Join(execDir, historyName)

	historyPath, err := filepath.Abs(relativeHistoryPath)
	if err != nil {
		log.Fatalf(
			"main.main(): Error finding the absolute path of the executable (%s): %s",
			relativeHistoryPath,
			err.Error(),
		)
	}

	var (
		historyJson []byte
		history     history = history{}
	)

	historyJson, err = os.ReadFile(historyPath)
	if err != nil {
		log.Fatalf("main.main(): Error reading '%s': %s", makeLink("file://"+historyPath), err.Error())
	}

	// fmt.Printf("%s\n\n", historyJson)

	err = json.Unmarshal(historyJson, &history)
	if err != nil {
		log.Fatalf("main.main(): Error parsing '%s': %s", makeLink("file://"+historyPath), err.Error())
	}

	t := history.MsPerChar

	step, err := history.Steps.getStep(firstStepName)
	if err != nil {
		log.Fatalf("main.main(): Failed to find step %s: %s", strconv.Quote(firstStepName), err.Error())
	}

	for {
		printWithTime(t, step.Text, "\n\n")

		takenIndexes := []int{}
		i := 0
		for _, o := range step.Options {
			if o.Number != 0 && doesntContain(takenIndexes, o.Number) {
				takenIndexes = append(takenIndexes, o.Number)
			}
		}

		for _, o := range step.Options {
			// if
		}
	}

	// fmt.Printf("%#v\n", history)
}
