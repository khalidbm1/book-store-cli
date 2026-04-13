package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Test behaviors for getBooks, saveBooks, and handlers.

// Helper to create a temporary books.json file and restore after.
func withTempBooksFile(t *testing.T, content []Book, fn func(tmpPath string)) {
	t.Helper()

	// Marshal provided books slice
	data, err := json.Marshal(content)
	if err != nil {
		t.Fatalf("failed to marshal books: %v", err)
	}

	// Create temp dir
	tmpDir := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get wd: %v", err)
	}

	// Change working directory so that ./books.json is resolved inside temp dir
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(oldWd)

	// Write books.json
	if err := os.WriteFile(filepath.Join(tmpDir, "books.json"), data, 0644); err != nil {
		t.Fatalf("failed to write temp books.json: %v", err)
	}

	fn(filepath.Join(tmpDir, "books.json"))
}

func TestGetBooks_FileMissingReturnsError(t *testing.T) {
	// Behavior: getBooks should return an error when books.json does not exist.

	// Use temp dir without creating books.json
	tmpDir := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get wd: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(oldWd)

	_, err = getBooks()
	if err == nil {
		t.Fatalf("expected error when books.json is missing, got nil")
	}
}

func TestGetBooks_InvalidJSONReturnsError(t *testing.T) {
	// Behavior: getBooks should return an error when books.json contains invalid JSON.

	tmpDir := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get wd: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(oldWd)

	// Write invalid JSON
	if err := os.WriteFile("books.json", []byte("not-json"), 0644); err != nil {
		t.Fatalf("failed to write invalid books.json: %v", err)
	}

	_, err = getBooks()
	if err == nil {
		t.Fatalf("expected error for invalid JSON, got nil")
	}
}

func TestSaveBooks_WritesFileWithContent(t *testing.T) {
	// Behavior: saveBooks should write the provided books slice to books.json.

	tmpDir := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get wd: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(oldWd)

	books := []Book{{ISBN: "123", Title: "Test"}}
	if err := saveBooks(books); err != nil {
		t.Fatalf("saveBooks returned error: %v", err)
	}

	data, err := os.ReadFile("books.json")
	if err != nil {
		t.Fatalf("failed to read books.json: %v", err)
	}

	var got []Book
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("failed to unmarshal saved data: %v", err)
	}

	if len(got) != 1 || got[0].ISBN != "123" || got[0].Title != "Test" {
		t.Fatalf("unexpected saved books: %+v", got)
	}
}

func TestHandleAddBook_UpdatesExistingBook(t *testing.T) {
	// Behavior: handleAddBook should update an existing book when ISBN already exists.

	initial := []Book{{
		ISBN:    "111",
		Title:   "Old Title",
		Author:  "Old Author",
		Price:   10.0,
		Stock:   1,
		Reviews: []string{"ok"},
		Rating:  3.0,
	}}

	withTempBooksFile(t, initial, func(tmpPath string) {
		// Prepare flags
		addCli := flag.NewFlagSet("add", flag.ContinueOnError)
		isbn := addCli.String("isbn", "111", "")
		title := addCli.String("title", "New Title", "")
		author := addCli.String("author", "New Author", "")
		price := addCli.Float64("price", 20.0, "")
		imageURL := addCli.String("image", "http://image", "")
		description := addCli.String("description", "desc", "")
		category := addCli.String("category", "cat", "")
		publishedDate := addCli.String("published", "2020-01-01", "")
		stock := addCli.Int("stock", 5, "")
		reviews := addCli.String("reviews", "good,better", "")
		rating := addCli.Float64("rating", 4.5, "")

		// Simulate CLI args for Parse
		oldArgs := os.Args
		os.Args = []string{"cmd", "add"}
		defer func() { os.Args = oldArgs }()

		// Call handler
		handleAddBook(addCli, isbn, title, author, price, imageURL, description, category, publishedDate, stock, reviews, rating)

		// Verify file contents
		data, err := os.ReadFile(tmpPath)
		if err != nil {
			t.Fatalf("failed to read books.json: %v", err)
		}
		var got []Book
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("expected 1 book, got %d", len(got))
		}
		if got[0].Title != "New Title" || got[0].Author != "New Author" || got[0].Price != 20.0 {
			t.Fatalf("book not updated as expected: %+v", got[0])
		}
	})
}

func TestHandleAddBook_AppendsNewBook(t *testing.T) {
	// Behavior: handleAddBook should append a new book when ISBN does not exist.

	initial := []Book{}

	withTempBooksFile(t, initial, func(tmpPath string) {
		addCli := flag.NewFlagSet("add", flag.ContinueOnError)
		isbn := addCli.String("isbn", "222", "")
		title := addCli.String("title", "Brand New", "")
		author := addCli.String("author", "Author", "")
		price := addCli.Float64("price", 30.0, "")
		imageURL := addCli.String("image", "http://image2", "")
		description := addCli.String("description", "desc2", "")
		category := addCli.String("category", "cat2", "")
		publishedDate := addCli.String("published", "2021-01-01", "")
		stock := addCli.Int("stock", 10, "")
		reviews := addCli.String("reviews", "nice", "")
		rating := addCli.Float64("rating", 5.0, "")

		oldArgs := os.Args
		os.Args = []string{"cmd", "add"}
		defer func() { os.Args = oldArgs }()

		handleAddBook(addCli, isbn, title, author, price, imageURL, description, category, publishedDate, stock, reviews, rating)

		data, err := os.ReadFile(tmpPath)
		if err != nil {
			t.Fatalf("failed to read books.json: %v", err)
		}
		var got []Book
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("expected 1 book, got %d", len(got))
		}
		if got[0].ISBN != "222" || got[0].Title != "Brand New" {
			t.Fatalf("unexpected book: %+v", got[0])
		}
	})
}

// Additional behaviors to test:
// 1) handleGetBooks with --all prints header and all books without exiting.
// 2) handleGetBooks with --ISBN for existing book prints details and does not exit.
// 3) handleGetBooks with --ISBN for missing book prints "Book not found" and exits with code 1.
// 4) handleDeleteBook with existing ISBN removes the book and saves file.
// 5) handleDeleteBook with missing ISBN prints "Book not found" and exits with code 1.

// To test exit paths, we temporarily replace osExit.
var (
	origStdout *os.File
	origStderr *os.File
)

// captureOutput captures stdout during fn execution.
func captureOutput(t *testing.T, fn func()) string {
	t.Helper()

	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("failed to read from pipe: %v", err)
	}
	return buf.String()
}

func TestHandleGetBooks_AllPrintsHeaderAndBooks(t *testing.T) {
	initial := []Book{{
		ISBN:   "1",
		Title:  "Title1",
		Author: "Author1",
		Price:  10,
		Rating: 4.0,
	}}

	withTempBooksFile(t, initial, func(tmpPath string) {
		_ = tmpPath // file path not needed directly; getBooks uses CWD

		getCli := flag.NewFlagSet("get", flag.ContinueOnError)
		all := getCli.Bool("all", true, "")
		isbn := getCli.String("ISBN", "", "")

		oldArgs := os.Args
		os.Args = []string{"cmd", "get", "--all"}
		defer func() { os.Args = oldArgs }()

		out := captureOutput(t, func() {
			handleGetBooks(getCli, all, isbn)
		})

		if !strings.Contains(out, "ISBN") || !strings.Contains(out, "Title1") {
			t.Fatalf("expected output to contain header and book, got: %s", out)
		}
	})
}

func TestHandleGetBooks_ISBNNotFoundExits(t *testing.T) {
	initial := []Book{{
		ISBN:   "1",
		Title:  "Title1",
		Author: "Author1",
		Price:  10,
		Rating: 4.0,
	}}

	withTempBooksFile(t, initial, func(tmpPath string) {
		_ = tmpPath

		origExit := osExit
		osExit = func(code int) { panic(code) }
		defer func() { osExit = origExit }()

		getCli := flag.NewFlagSet("get", flag.ContinueOnError)
		all := getCli.Bool("all", false, "")
		isbn := getCli.String("ISBN", "999", "")

		oldArgs := os.Args
		os.Args = []string{"cmd", "get", "--ISBN", "999"}
		defer func() { os.Args = oldArgs }()

		var exited bool
		out := captureOutput(t, func() {
			defer func() {
				if r := recover(); r != nil {
					exited = true
				}
			}()
			handleGetBooks(getCli, all, isbn)
		})

		if !exited {
			t.Fatal("expected osExit to be called")
		}
		if !strings.Contains(out, "Book not found") {
			t.Fatalf("expected 'Book not found' in output, got: %s", out)
		}
	})
}

func TestHandleDeleteBook_RemovesExistingBook(t *testing.T) {
	initial := []Book{{
		ISBN:   "1",
		Title:  "Title1",
		Author: "Author1",
	}}

	withTempBooksFile(t, initial, func(tmpPath string) {
		deleteCli := flag.NewFlagSet("delete", flag.ContinueOnError)
		isbn := deleteCli.String("ISBN", "1", "")

		oldArgs := os.Args
		os.Args = []string{"cmd", "delete", "--ISBN", "1"}
		defer func() { os.Args = oldArgs }()

		captureOutput(t, func() {
			handleDeleteBook(deleteCli, isbn)
		})

		data, err := os.ReadFile(tmpPath)
		if err != nil {
			t.Fatalf("failed to read books.json: %v", err)
		}
		var got []Book
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}
		if len(got) != 0 {
			t.Fatalf("expected 0 books after delete, got %d", len(got))
		}
	})
}

func TestHandleDeleteBook_ISBNNotFoundPrintsMessage(t *testing.T) {
	initial := []Book{{
		ISBN:   "1",
		Title:  "Title1",
		Author: "Author1",
	}}

	withTempBooksFile(t, initial, func(tmpPath string) {
		_ = tmpPath

		origExit := osExit
		osExit = func(code int) { panic(code) }
		defer func() { osExit = origExit }()

		deleteCli := flag.NewFlagSet("delete", flag.ContinueOnError)
		isbn := deleteCli.String("ISBN", "999", "")

		oldArgs := os.Args
		os.Args = []string{"cmd", "delete", "--ISBN", "999"}
		defer func() { os.Args = oldArgs }()

		var exited bool
		out := captureOutput(t, func() {
			defer func() {
				if r := recover(); r != nil {
					exited = true
				}
			}()
			handleDeleteBook(deleteCli, isbn)
		})

		if !exited {
			t.Fatal("expected osExit to be called")
		}
		if !strings.Contains(out, "Book not found") {
			t.Fatalf("expected 'Book not found' in output, got: %s", out)
		}
	})
}

// Dummy reference to avoid unused import error for errors package in case of future extension.
var _ = errors.New
