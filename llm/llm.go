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
		"You will be provided the output from a shell command and a request regarding the output, respond directly and succinctly " +
		"in under two sentences.")
	p.WriteString(getInputSeq(template))
	p.WriteString(fmt.Sprintf("```%v```\n%v", input, request))
	p.WriteString(getResponseSeq(template))
	return p.String()
}

func GetGeneratePrompt(template string, request string) string {
	var p strings.Builder
	p.WriteString(getSystemSeq(template))
	p.WriteString(PERF_STIM + "You will be given a request. Generate a command or sequence of commands, and ONLY the necessary command(s), " +
		"to satisfy the request.")
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
	case "capy":
		return "USER: "
	case "deepseek":
		return "User: "
	}

	return "### Instruction:\n"
}

func getInputSeq(template string) string {
	switch template {
	case "chatml":
		return "<|im_end|>\n<|im_start|>user\n"
	case "mistral":
		return " [/INST] Got it. Please provide the details.</s> [INST] "
	case "internlm":
		return "[UNUSED_TOKEN_145]\n[UNUSED_TOKEN_146]user\n"
	case "capy":
		return " ASSISTANT: Got it. Please provide the details.</s> USER: "
	case "deepseek":
		return "\n\nAssistant: Got it. Please provide the details.<｜end▁of▁sentence｜>User: "
	}

	return "\n\n### Input:\n"
}

func getResponseSeq(template string) string {
	switch template {
	case "chatml":
		return "<|im_end|>\n<|im_start|>assistant\n"
	case "mistral":
		return " [/INST]"
	case "internlm":
		return "[UNUSED_TOKEN_145]\n[UNUSED_TOKEN_146]assistant\n"
	case "capy":
		return " ASSISTANT:"
	case "deepseek":
		return "\n\nAssistant:"
	}

	return "\n\n### Response:\n"
}
