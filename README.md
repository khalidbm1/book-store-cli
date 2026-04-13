# book-store-cli

> A bookstore. In your terminal. Because why not.

## What is this?

A CLI tool for managing books — written in Go for absolutely no reason other than **we wanted to play with Go and needed an excuse**.

No database. No frontend. No users. No revenue. Just a `books.json` file and vibes.

## Why does this exist?

```
Me:  I want to PLAY with Go
Also me: Let's build... a bookstore?
Also also me: Yes.
```

That's it. That's the whole origin story.

## Features

- Get books (wow)
- Add books (incredible)
- Update books (life-changing)
- Delete books (controversial)
- Store everything in a flat JSON file like it's 2004 (authentic)

## Usage

```bash
# See all your books (both of them)
./bookstore get --all

# Find that one book you definitely didn't forget the ISBN of
./bookstore get --ISBN 5

# Add a book. All fields required. Yes, all of them. No shortcuts.
./bookstore add --ISBN 6 --title "My Cool Book" --author "Some Guy" --price 200 --image_url http://fake.url/book.png

# Update a book because you typo'd the price as 20000 instead of 20
./bookstore update --ISBN 6 --title "My Cool Book" --price 20

# Delete a book. It's gone forever. No undo. No confirmation. Just gone.
./bookstore delete --ISBN 6
```

## Architecture

```
books.json  <-- The "database"
books.go    <-- The "backend"  
main.go     <-- The "API"
```

Silicon Valley would not be proud. We are.

## Known Issues

- If `books.json` doesn't exist, the whole thing panics. Working as intended.
- The CLI parses ALL flag sets regardless of which command you used. This is fine.
- No pagination. If you have 10,000 books, good luck reading that terminal output.
- Deleting by index in a loop is technically a bug. We shipped it anyway.

## Tech Stack

- Go (the point)
- JSON (the database)
- `flag` package (the framework)
- One `books.json` file (the cloud)

## Contributing

This is a toy project made for fun. PRs welcome but also unnecessary. The code works. Mostly.

## License

Do whatever. It's a bookstore CLI. Nobody is getting rich here.

---

*Built with Go, curiosity, and zero production requirements.*
