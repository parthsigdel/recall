## recall
Store the commands you keep forgetting but know you'll need again.

Remove the inconveniency of searching the same commands online (or your shell history), copying them, pasting them , and refining them again according to your needs. 
Have them ready to be executed instantly — right on your terminal prompt, no copy-pasting.

## Demonstration/Usage
**Search and run:**
Hit `Ctrl+F` to open the search. Find your command and press `Enter` — it'll be placed on your terminal prompt, ready to run.
<img width="1920" height="1080" alt="search_in_recall" src="https://github.com/user-attachments/assets/774eb1a8-e95c-4d14-b234-9f71d7d507b8" />

#
**Add a command:**
```bash
recall -s
```
<img width="1117" height="617" alt="add_in_recall" src="https://github.com/user-attachments/assets/4fef23f2-1dcc-41b6-8f65-550872d40deb" />

#
**Delete a command:**
Hit `Ctrl+F` to open the search, find the command you want to remove, and press `Ctrl+D`.
<img width="1920" height="1080" alt="delete_in_recall" src="https://github.com/user-attachments/assets/1e16d462-7eb4-4d50-a44d-9281741f2dc2" />


## Installation guide

### macOS / Linux

```
curl -fsSL https://raw.githubusercontent.com/parthsigdel/recall/main/install.sh | sh
```

### Windows (Powershell)

```
irm https://raw.githubusercontent.com/parthsigdel/recall/main/install.ps1 | iex
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

`powershell` / `pwsh` (run `$PROFILE` to find your profile path):
```
Invoke-Expression (recall -init powershell)
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

`powershell`:
```
. $PROFILE
```

## NOTE
If you have already set up a keybinding involving CTRL+F, please skip the shell integration process mentioned above. Instead, just use "recall" to start it. 
