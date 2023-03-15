package gpt3

import (
	"context"
)

type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionRequest is a request for the chat/completions API
type ChatCompletionRequest struct {
	Model    EngineType              `json:"model"`
	Messages []ChatCompletionMessage `json:"messages"`
	// The maximum number of tokens allowed for the generated answer. By default, the number of tokens the model can return will be (4096 - prompt tokens).
	MaxTokens *int `json:"max_tokens,omitempty"`
	// Sampling temperature to use
	Temperature *float32 `json:"temperature,omitempty"`
	// Alternative to temperature for nucleus sampling
	TopP *float32 `json:"top_p,omitempty"`
	// How many chat completion choices to generate for each input message.
	N *int `json:"n"`
	// Up to 4 sequences where the API will stop generating tokens. Response will not contain the stop sequence.
	Stop []string `json:"stop,omitempty"`
	// Number between -2.0 and 2.0. Positive values penalize new tokens based on whether they appear in the text so far, increasing the model's likelihood to talk about new topics.
	PresencePenalty float32 `json:"presence_penalty"`
	// Number between -2.0 and 2.0. Positive values penalize new tokens based on their existing frequency in the text so far, decreasing the model's likelihood to repeat the same line verbatim.
	FrequencyPenalty float32 `json:"frequency_penalty"`
	// Modify the likelihood of specified tokens appearing in the completion.
	// Accepts a json object that maps tokens (specified by their token ID in the tokenizer) to an associated bias value from -100 to 100. Mathematically, the bias is added to the logits generated by the model prior to sampling. The exact effect will vary per model, but values between -1 and 1 should decrease or increase likelihood of selection; values like -100 or 100 should result in a ban or exclusive selection of the relevant token.
	LogitBias map[string]string `json:"logit_bias,omitempty"`

	// Whether to stream back results or not. Don't set this value in the request yourself
	// as it will be overriden depending on if you use CompletionStream or Completion methods.
	Stream bool `json:"stream,omitempty"`
}

type ChatCompletionResponseChoiceMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// CompletionResponseChoice is one of the choices returned in the response to the Completions API
type ChatCompletionResponseChoice struct {
	// Index        int                                 `json:"index"`
	Message      ChatCompletionResponseChoiceMessage `json:"delta"`
	FinishReason string                              `json:"finish_reason"`
}

/*
{
	"id": "chatcmpl-6pTTLGwAbhbdi2SfN8yUdAnQLgSyT",
	"object": "chat.completion.chunk",
	"created": 1677726039,
	"model": "gpt-3.5-turbo-0301",
	"choices": [{
		"delta": {
			"content": "块"
		},
		"index": 0,
		"finish_reason": null
	}]
}
*/

// CompletionResponse is the full response from a request to the completions API
type ChatCompletionResponse struct {
	// ID      string                         `json:"id"`
	// Object  string                         `json:"object"`
	// Created int                            `json:"created"`
	// Model   string                         `json:"model"`
	Choices []ChatCompletionResponseChoice `json:"choices"`
	Usage   ChatCompletionResponseUsage    `json:"usage"`
}

func (cr *ChatCompletionResponse) CanContinue() bool {
	if cr != nil && len(cr.Choices) > 0 {
		return cr.Choices[0].FinishReason == "length"
	}
	return false
}

func (cr *ChatCompletionResponse) Text() string {
	if cr != nil && len(cr.Choices) > 0 {
		return cr.Choices[0].Message.Content
	}
	return ""
}

func (cr *ChatCompletionResponse) Role() string {
	if cr != nil && len(cr.Choices) > 0 {
		return cr.Choices[0].Message.Role
	}
	return "assistant"
}

func (cr *ChatCompletionResponse) TotalTokens() int {
	if cr != nil {
		return cr.Usage.TotalTokens
	}
	return 0
}

func (cr *ChatCompletionResponse) Reset() {
	if cr != nil {
		*cr = ChatCompletionResponse{}
	}
	return
}

// CompletionResponseUsage is the object that returns how many tokens the completion's request used
type ChatCompletionResponseUsage struct {
	// PromptTokens     int `json:"prompt_tokens"`
	// CompletionTokens int `json:"completion_tokens"`
	TotalTokens int `json:"total_tokens"`
}

func (c *client) ChatCompletion(ctx context.Context, request ChatCompletionRequest) (*ChatCompletionResponse, error) {
	request.Stream = false
	req, err := c.newRequest(ctx, "POST", "/chat/completions", request)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := new(ChatCompletionResponse)
	if err := getResponseObject(resp, output); err != nil {
		return nil, err
	}
	return output, nil
}

func (c *client) ChatCompletionStream(ctx context.Context, request ChatCompletionRequest, onData func(CompletionResponseInterface)) error {
	return c.ChatCompletionStreamWithEngine(ctx, request, onData)
}

func (c *client) ChatCompletionStreamWithEngine(
	ctx context.Context,
	request ChatCompletionRequest,
	onData func(CompletionResponseInterface),
) error {
	request.Stream = true
	req, err := c.newRequest(ctx, "POST", "/chat/completions", request)
	if err != nil {
		return err
	}
	return c.sendAndOnData(req, new(ChatCompletionResponse), onData)
}
