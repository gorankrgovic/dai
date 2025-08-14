package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gorankrgovic/dai/internal/buildinfo"
	"github.com/spf13/cobra"
)

var (
	parrotMode string // "", "party", "insult", "wise"
)

func versionString() string {
	v := buildinfo.Version
	if v == "" {
		v = "dev"
	}

	if c := buildinfo.Commit; c != "" {
		if len(c) > 8 {
			c = c[:8]
		}
		v = fmt.Sprintf("%s (%s)", v, c)
	}
	if d := buildinfo.Date; d != "" {
		v = fmt.Sprintf("%s %s", v, d)
	}
	return v
}

var rootCmd = &cobra.Command{
	Use:     "dai",
	Short:   "DAI — Debug & Develop AI CLI",
	Long:    "DAI — Debug & Develop AI. Triage and autofix directly from terminal.",
	Version: versionString(),
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&parrotMode, "parrot", "p", "", "Summon the DAI parrot (modes: party, insult, wise)")
}

// Execute is the main entry point
func Execute() {
	mode := parseParrotMode(os.Args)
	if mode != "" {
		showParrot(mode)
		os.Exit(0)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func parseParrotMode(args []string) string {
	for _, a := range args {
		if a == "--parrot" || a == "-p" {
			return "basic"
		}
		if strings.HasPrefix(a, "--parrot=") {
			return strings.ToLower(strings.TrimPrefix(a, "--parrot="))
		}
		if strings.HasPrefix(a, "-p=") {
			return strings.ToLower(strings.TrimPrefix(a, "-p="))
		}
	}
	return ""
}

func showParrot(mode string) {
	rand.Seed(time.Now().UnixNano())

	// Poruke
	basicMsgs := []string{
		"Polly wants a pull request!",
		"Deploy or not deploy… squawk!",
		"Refactor me, human!",
		"Tests or it didn’t happen.",
		"One more log, one less bug.",
	}
	insultMsgs := []string{
		"Squawk! Your commit smells like Friday 5:59 PM!",
		"Your CI is green because it runs nothing — impressive.",
		"Bro, the linter cried reading this.",
		"You packed a monolith into a microservice… with no network. Genius.",
		"Tip: ‘TODO’ is not a feature spec.",
		"Wow, an AI-powered form validator… truly changing the world.",
		"Your AI model just predicted the sun will rise tomorrow. Stunning innovation.",
		"You rebranded regex as ‘machine learning’. VC bait 101.",
		"Enterprise-grade? You mean it takes 3 months to deploy a button.",
		"Your code review process has more bureaucracy than a Balkan DMV.",
		"Squawk! I’ve seen faster releases from a government website.",
		"You wrapped GPT output in JSON and called it a proprietary LLM. Sure, boss.",
		"Your AI ‘assistant’ just sent 5 meetings to plan a meeting.",
		"Congratulations, you replaced junior devs with a chatbot that hallucinates.",
		"Your JIRA board has more tickets than your app has users.",
		"AI ethics committee? More like PR firewall for bad press.",
		"You trained your model on 90% Stack Overflow, 10% wishful thinking.",
		"Enterprise cloud migration: now 10x slower, but at least it’s expensive.",
		"Your AI roadmap is just buzzwords in alphabetical order.",
		"Squawk! Even my parrot could automate that workflow — for peanuts.",
	}
	wiseMsgs := []string{
		"Measure twice, deploy once.",
		"Small diffs ship faster than big dreams.",
		"Failures are logs in disguise.",
		"Readability scales, hacks don’t.",
		"Optimize last; test first.",
	}

	msg := ""
	switch mode {
	case "insult":
		msg = pick(insultMsgs)
	case "wise":
		msg = pick(wiseMsgs)
	default: // basic & party
		msg = pick(basicMsgs)
	}

	// Bubble
	bubble := renderSpeechBubble(msg, 48)

	// Deluxe ASCII
	parrotArt := `
                      @@@@@@@@@@@@@@@@@@
                  @@@@@@@@@@@@@@@@@@@@@@@@@@
              @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
            @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
          @@@@@@@@@@@@@@@@::::::::::::::@@@@@@@@@@
        @@@@@@@@@@@@@@::::::    @@@@@@  ::::@@@@@@
        @@@@@@@@@@@@@@::  ::  @@@@@@@@@@::==::==@@
      @@@@@@@@@@@@@@::  ::    @@@@@@@@@@::::==@@@@@@
      @@@@@@@@@@@@@@@@@@@@::    @@@@@@  ::==@@@@==@@@@
    @@@@@@@@@@@@@@@@        ::          ::@@@@==::==@@@@
    @@@@@@@@@@@@      @@@@@@@@@@::  ::::@@@@======  ==@@@@
    @@@@@@@@@@    @@@@::  ::  ::  ::::@@@@==========  ==@@
  @@@@@@@@      @@  ::  ::  ::  ::@@@@@@@@@@@@========::@@@@
@@@@@@==@@@@  @@  ::  ::  ::  ::@@@@@@@@@@@@@@@@======::==@@
@@@@@@@@      @@::  ::  ::  ::::@@==@@@@  @@@@@@==========@@
@@@@====@@    @@  ::  ::  ::::@@@@==@@@@  @@@@@@@@====::==@@
@@==@@        @@::  ::  ::::==@@====@@@@@@  @@@@@@@@======@@
@@==    @@    @@@@::::::::==::@@======@@@@  @@@@@@@@====@@@@
==@@@@@@    ==@@@@@@@@====::==@@======@@@@      @@@@==@@@@@@
@@@@@@    ==  @@@@@@@@@@@@@@@@@@@@@@====@@@@@@  @@@@@@==@@
@@@@  @@==  ==@@@@@@@@@@@@==@@  @@@@@@@@@@@@@@  @@@@==@@@@
@@==@@==  ==  ==@@@@@@@@==@@@@                  @@@@@@@@
==@@==  @@  ==  ========@@@@@@                @@@@@@@@
@@@@==@@  ==  ====@@@@==@@@@@@   dD       @@@@@@@@
@@====@@====@@==@@==@@==@@@@@@@@ -The DAI Dude
==@@@@@@==@@==@@@@@@==@@@@@@@@@@
`

	// PARTY
	if mode == "party" {
		fetti(5, 64, 90*time.Millisecond) // 5 waves of confetti
	}

	// Write
	fmt.Println(bubble)
	fmt.Println(parrotArt)

	// PARTY outro
	if mode == "party" {
		fetti(3, 64, 80*time.Millisecond)
	}

	// Hint on error
	if mode != "basic" && mode != "party" && mode != "insult" && mode != "wise" {
		fmt.Println("ℹ️  Unknown parrot mode. Try: basic (default), party, insult, wise")
	}
}

func pick[T any](arr []T) T { return arr[rand.Intn(len(arr))] }

// Minimal terminal
func fetti(waves, width int, delay time.Duration) {
	colors := []string{
		"\x1b[31m", "\x1b[32m", "\x1b[33m",
		"\x1b[34m", "\x1b[35m", "\x1b[36m",
	}
	syms := []rune{'*', '.', '^', 'o', '~', '+', '×'}
	reset := "\x1b[0m"
	for i := 0; i < waves; i++ {
		var b strings.Builder
		indent := rand.Intn(10)
		b.WriteString(strings.Repeat(" ", indent))
		lineLen := width - indent
		if lineLen < 20 {
			lineLen = 20
		}
		for j := 0; j < lineLen; j++ {
			c := colors[rand.Intn(len(colors))]
			s := syms[rand.Intn(len(syms))]
			b.WriteString(c)
			b.WriteRune(s)
		}
		b.WriteString(reset)
		fmt.Println(b.String())
		time.Sleep(delay)
	}
}

// Render talk-bubble
func renderSpeechBubble(text string, width int) string {
	lines := wrap(text, width)
	var sb strings.Builder
	sb.WriteString("  " + strings.Repeat("_", width+2) + "\n")
	for _, line := range lines {
		sb.WriteString(fmt.Sprintf(" / %-*s \\\n", width, line))
	}
	sb.WriteString("  " + strings.Repeat("-", width+2) + "\n")

	sb.WriteString("          \\\n")
	sb.WriteString("           \\\n")
	return sb.String()
}

func wrap(s string, width int) []string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return []string{""}
	}
	var lines []string
	var line string
	for _, w := range words {
		if len(line)+len(w)+1 > width {
			lines = append(lines, line)
			line = w
		} else {
			if line == "" {
				line = w
			} else {
				line += " " + w
			}
		}
	}
	if line != "" {
		lines = append(lines, line)
	}
	return lines
}
