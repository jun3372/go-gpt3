package gpt3

import (
	"context"
	"errors"
)

type GPT3client struct {
	client    Client
	maxtokens int
}

func MakeGPT3Client(apikey string, maxtokens int, options ...ClientOption) *GPT3client {
	return &GPT3client{
		client: NewClient(
			apikey,
			options...),
		maxtokens: maxtokens,
	}
}

func (c *GPT3client) DoStream(ctx context.Context, say []ChatCompletionMessage, fn func(cr CompletionResponseInterface)) error {
	if len(say) == 0 {
		return errors.New("您得说些什么。")
	}
	if c.client.DefaultEngine() == Gpt35TurboEngine ||
		c.client.DefaultEngine() == Gpt35Turbo0301Engine {
		request := ChatCompletionRequest{
			Model:     c.client.DefaultEngine(),
			Messages:  say,
			MaxTokens: IntPtr(c.maxtokens),
		}
		return c.client.ChatCompletionStream(ctx, request, fn)
	}
	// 组装 内容
	text := make([]string, len(say))
	for idx, v := range say {
		text[idx] = v.Content
	}
	// text := strings.Builder{}
	// for _, v := range say {
	// 	text.WriteString(v.Content)
	// }
	// tstr := text.String()
	// // max 4096 限制
	// if l := len(tstr); l > 4096 {
	// 	tstr = string([]rune(tstr[l-4096:])[2:])
	// }
	request := CompletionRequest{
		Prompt:    text,
		MaxTokens: IntPtr(c.maxtokens),
	}
	return c.client.CompletionStream(ctx, request, fn)
}

func (c *GPT3client) DoOnce(ctx context.Context, say []ChatCompletionMessage) (CompletionResponseInterface, error) {
	if len(say) == 0 {
		return nil, errors.New("您得说些什么。")
	}
	if c.client.DefaultEngine() == Gpt35TurboEngine ||
		c.client.DefaultEngine() == Gpt35Turbo0301Engine {
		request := ChatCompletionRequest{
			Model:     Gpt35TurboEngine,
			Messages:  say,
			MaxTokens: IntPtr(c.maxtokens),
		}
		return c.client.ChatCompletion(ctx, request)
	}
	request := CompletionRequest{
		Prompt:    []string{say[0].Content},
		MaxTokens: IntPtr(c.maxtokens),
	}
	return c.client.Completion(ctx, request)
}
