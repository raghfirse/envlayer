# envlayer

Hierarchical environment variable manager that merges `.env` files by environment context.

---

## Installation

```bash
go install github.com/yourusername/envlayer/cmd/envlayer@latest
```

Or add it as a library:

```bash
go get github.com/yourusername/envlayer
```

---

## Usage

`envlayer` merges `.env` files in order of precedence — base values are overridden by environment-specific ones.

**File structure:**
```
.env
.env.production
.env.local
```

**Load and merge in your Go application:**

```go
package main

import (
    "fmt"
    "github.com/yourusername/envlayer"
)

func main() {
    env, err := envlayer.Load("production")
    if err != nil {
        panic(err)
    }

    fmt.Println(env.Get("DATABASE_URL"))
}
```

`envlayer` will automatically merge `.env` → `.env.production` → `.env.local`, with later files taking priority.

**CLI usage:**

```bash
envlayer run --env production -- ./myapp
```

---

## How It Works

| File | Priority |
|------|----------|
| `.env` | Base (lowest) |
| `.env.<environment>` | Environment-specific |
| `.env.local` | Local overrides (highest) |

---

## License

MIT © 2024 yourusername