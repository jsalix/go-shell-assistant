package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jsalix/go-shell-assistant/api"
	"github.com/jsalix/go-shell-assistant/llm"
)

var (
	template string
	mode     string
	request  string
	client   *api.KoboldClient
)

const (
	API_URL = "http://192.168.1.55:5001/api"
)

func main() {
	// create koboldcpp client
	client, err := api.NewKoboldClient(API_URL)
	if err != nil {
		fmt.Println("unable to create koboldcpp client", err)
		return
	}

	// allow changing template strategy
	flag.StringVar(&template, "t", "chatml", "sets the prompting template")

	// check which input flags are set to determine mode: explain, generate
	var (
		e string
		g string
	)

	flag.StringVar(&e, "e", "", "explain stdin")
	flag.StringVar(&g, "g", "", "generate command")

	flag.Parse()

	if e != "" {
		mode = "explain"
		request = e
	} else if g != "" {
		mode = "generate"
		request = g
	} else {
		fmt.Println("missing necessary flags")
		os.Exit(1)
	}

	// build prompt, send request to llm, print response
	stopStrings := []string{"\n\n\n", "\n###", "\n---", "\n<|", "</s>", "<|im_end|>", "[UNUSED_TOKEN_145]"}

	var prompt string

	switch mode {
	case "explain":
		// get stdin
		stdin, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println("unable to read stdin", err)
			os.Exit(1)
		}
		input := strings.TrimSuffix(string(stdin), "\n")
		// TODO is it possible to get the last command too? (i.e. with !!) -- might not grab the right command tbh
		prompt = llm.GetExplainPrompt(template, string(input), request)
	case "generate":
		prompt = llm.GetGeneratePrompt(template, request)
		stopStrings = append(stopStrings, "```", "\n") // attempt to limit output to command only
	}

	params := &api.KoboldParams{
		MaxContextLength: 4096,
		MaxLength:        300,
		Temperature:      1,
		DynaTempRange:    0,
		TopP:             1,
		MinP:             0,
		TopK:             0,
		TopA:             0,
		Typical:          1.0,
		Tfs:              1.0,
		RepPen:           1.0,
		RepPenRange:      128,
		RepPenSlope:      0,
		SamplerOrder:     []int{6, 0, 1, 3, 4, 2, 5},
		SamplerSeed:      -1,
		StopSequence:     stopStrings,
		BanTokens:        false,
		TrimStop:         true,
		Prompt:           prompt,
	}

	response, err := client.Generate(params)
	if err != nil {
		fmt.Println("there was an error while waiting on a response from koboldcpp!", err)
	}

	// check if koboldcpp is already processing a request
	if response.Status == "busy" {
		fmt.Println("koboldcpp is busy, try again later")
		os.Exit(1)
		// TODO implement retry mechanism?
	}

	if response.Status == "ok" {
		fmt.Print(trimSuffixes(response.Text, &stopStrings))
	}
}

func trimSuffixes(str string, suffixes *[]string) string {
	for _, suffix := range *suffixes {
		trimmedStr, trimmed := strings.CutSuffix(str, suffix)
		if trimmed {
			return trimmedStr
		}
	}
	return str
}
