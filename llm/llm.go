package llm

import (
	"fmt"
	"strings"
)

const PERF_STIM = "You are an extremely efficient shell terminal assistant, please ensure all your responses are of highest quality to guarantee the continued success of my career. "

func GetExplainPrompt(template string, input string, request string) string {
	var p strings.Builder
	p.WriteString(getSystemSeq(template))
	p.WriteString(PERF_STIM +
		"You will be provided the output from a command and a user request, respond directly and succinctly " +
		"in two sentences or less.")
	p.WriteString(getInputSeq(template))
	p.WriteString(fmt.Sprintf("```%v```\n%v", input, request))
	p.WriteString(getResponseSeq(template))
	return p.String()
}

func GetGeneratePrompt(template string, request string) string {
	var p strings.Builder
	p.WriteString(getSystemSeq(template))
	p.WriteString(PERF_STIM + "Generate a command or sequence of commands, and ONLY the necessary command(s), " +
		"to satisfy the following request.")
	p.WriteString(getInputSeq(template))
	p.WriteString(request)
	p.WriteString(getResponseSeq(template))
	p.WriteString("```terminal\n>")
	return p.String()
}

func getSystemSeq(template string) string {
	switch template {
	case "chatml":
		return "<|im_start|>system\n"
	case "mistral":
		return "[INST] "
	case "internlm":
		return "[UNUSED_TOKEN_146]system\n"
	}

	return "### Instruction:\n"
}

func getInputSeq(template string) string {
	switch template {
	case "chatml":
		return "<|im_end|>\n<|im_start|>user\n"
	case "mistral":
		return " "
	case "internlm":
		return "[UNUSED_TOKEN_145]\n[UNUSED_TOKEN_146]user\n"
	}

	return "\n\n### Input:\n"
}

func getResponseSeq(template string) string {
	switch template {
	case "chatml":
		return "<|im_end|>\n<|im_start|>assistant\n"
	case "mistral":
		return " [/INST] "
	case "internlm":
		return "[UNUSED_TOKEN_145]\n[UNUSED_TOKEN_146]assistant\n"
	}

	return "\n\n### Response:\n"
}
