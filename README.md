## ion

Ion is a minimal terminal wrapper for Windows that:
- Provides custom command aliases (with dynamic arguments) to compensate for Windows’ limited alias support.
- Stores personal secrets locally using SQLite, with encryption derived from a salt.
- Runs inside a TUI and executes your usual PowerShell commands in place.

### Why
Windows lacks a straightforward, portable alias system. Ion gives you a simple, local, self-contained way to:
- Define and use aliases that can call regular commands or Ion subcommands
- Save and retrieve secrets locally (no cloud)
- Keep everything in one terminal experience

## Features
- Alias management: add, update, rename, remove, list (supports JSON), search, and dynamic placeholders using `${ion}`
- Secret management: add (with optional salt and tags), update, rename, tag, list (JSON and decrypted views), search, remove, and copy to clipboard via `use`
- TUI wrapper around PowerShell with support for `cd`, `cls`, `clear`, and a zen mode

## Install, Build, and Run

Prerequisites:
- Go installed (module-aware). On Windows, use PowerShell.

Build:
```powershell
cd C:\dev\ion
go build -o .\bin\ion.exe .\cmd\main
```

Run:
```powershell
.\u005cbiniction.exe
```

Alternatively (no binary):
```powershell
go run .\cmd\main
```

Optional: add `C:\dev\ion\bin` to your PATH to run `ion` from anywhere.

## Data Storage
- Ion creates a data folder under your OS user config directory, then `ion\data`.
- SQLite database lives at: `<UserConfigDir>\ion\data\data.db`.

## Security Notes
- Secret values are encrypted locally using AES-GCM; the key is derived from the provided salt. Keep your salt(s) safe.
- `ion secret use <name>` copies the decrypted value to the clipboard.

## Usage Overview

### General help
```powershell
ion help
```
Outputs:
```
ion help

          ion zen - toggle zen mode
          ion help - show this help
          ion alias help - show the help for the alias command
          ion secret help - show the help for the secret command
```

### Zen mode
```powershell
ion zen
```

### Aliases
Help:
```powershell
ion alias help
```
Key actions and examples:
```powershell
# Add (supports "name=command" or spaced `name = command`)
ion alias add my-list=ion secret list

# Update
ion alias update my-list=ion secret list -j

# Rename
ion alias rename my-list my-list-json

# Remove
ion alias remove my-list-json

# List (human table or JSON)
ion alias list
ion alias list -j

# Search (fuzzy)
ion alias search my

# Dynamic placeholders with ${ion}
ion alias add greet=echo ${ion} said ${ion}
greet Bob Hello
# output: Bob said Hello
```

### Secrets
Help:
```powershell
ion secret help
```
Key actions and examples:
```powershell
# Add (salt optional; tags optional)
ion secret add email myPass123
ion secret add -s mySalt email myPass123
ion secret add -t work personal -s mySalt email myPass123

# Update value
ion secret update email newPass456

# Rename
ion secret rename email personal-email

# Tag (last arg is the name)
ion secret tag work laptop personal-email

# List (table; JSON; decrypted)
ion secret list
ion secret list -j
ion secret list -d
ion secret list -j -d

# Search (fuzzy)
ion secret search personal

# Remove
ion secret remove personal-email

# Copy decrypted value to clipboard
ion secret use personal-email
```

### Running system commands
Type your usual PowerShell commands; Ion executes them under the hood:
```powershell
cls
clear
cd ..
dir
```

## Pipeline of multiple commands
Use `&&` between commands. Ion will resolve aliases and Ion-subcommands per segment before executing.

## Roadmap
- Autocomplete
- UI finish/polish
- GitHub Actions (CI)

## Notes
- Ion is designed for Windows and uses PowerShell for system execution.
- Aliases are resolved inside Ion; they do not modify your shell’s global alias state.

