## recall
Store the commands you keep forgetting but know you'll need again.

Remove the inconveniency of searching the same commands online (or your shell history), copying them, pasting them , and refining them again according to your needs. 
Have them ready to be executed instantly — right on your terminal prompt, no copy-pasting.

## Demonstration/Usage
**Search and run:**
Hit `Ctrl+F` to open the search. Find your command and press `Enter` — it'll be placed on your terminal prompt, ready to run.

**Add a command:**
```bash
recall -s
```
You'll be prompted for a title and the command:
```
title: <your title>
command: <your command>
```

Hit enter, and it's stored instantly.

**Delete a command:**

Hit `Ctrl+F` to open the search, find the command you want to remove, and press `Ctrl+D`.


## Installation guide
```
go install github.com/CoderParth/recall@latest
```

Then make sure your Go bin is on your PATH. If not already added, please add this to your shell config:

`bash` / `zsh`:
```
export PATH=$PATH:$(go env GOPATH)/bin
```

`fish`:
```
fish_add_path (go env GOPATH)/bin
```

## Shell integration — IMPORTANT
Please add the following to your shell config file, according to your shell.

`bash` (`~/.bashrc`):
```
eval "$(recall -init bash)"
```

`zsh` (`~/.zshrc`):
```
eval "$(recall -init zsh)"
```

`fish` (`~/.config/fish/config.fish`):
```
recall -init fish | source
```

Then reload your config file:

`bash`:
```
source ~/.bashrc 
```

`zsh`:
```
source ~/.zshrc  
```

`fish`:
```
source ~/.config/fish/config.fish
```
