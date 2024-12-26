package views

import (
	"context"
	"io"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestLoginView(t *testing.T) {
	r, w := io.Pipe()
	go func() {
		_ = Login().Render(context.Background(), w)
		_ = w.Close()
	}()

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		t.Fatalf("failed to read template: %v", err)
	}

	if doc.Find(`form`).Length() == 0 {
		t.Error("expected form content to be rendered but was not")
	}
}
