package basic

import (
	"context"
	"testing"
)

func TestContextKey(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, Key("openid"), "不是群主")
	ctx = context.WithValue(ctx, Key2("openid"), "群主")

	if ctx.Value("openid") != nil {
		t.Fatal("fail openid must be nil")
	}
	if ctx.Value(Key("openid")) == ctx.Value((Key2("openid"))) {
		t.Fatal("fail Key(openid) must not equal Key2(openid)")
	}
}
