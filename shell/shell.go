package shell

import (
	"fmt"
)

const initBash = `
_recall() {
    local cmd
    cmd=$(recall) || return
    READLINE_LINE="$cmd"
    READLINE_POINT=${#READLINE_LINE}
}
bind -x '"\C-f":"_recall"'
`

const initZsh = `
_recall() {
    local cmd
    cmd=$(recall) || return
    BUFFER="$cmd"
    CURSOR=${#BUFFER}
    zle reset-prompt
}
zle -N _recall
bindkey '^F' _recall
`

const initFish = `
function _recall
    set cmd (recall)
    or return
    commandline --replace -- $cmd
    commandline --function repaint
end
bind \cf _recall
`

const initPowershell = `
Set-PSReadLineKeyHandler -Chord 'Ctrl+f' -ScriptBlock {
    param($key, $arg)

    $cmd = recall
    if ([string]::IsNullOrEmpty($cmd)) { return }

    [Microsoft.PowerShell.PSConsoleReadLine]::RevertLine()
    [Microsoft.PowerShell.PSConsoleReadLine]::Insert($cmd)
}
`

func Run(shell string) (string, error) {
	switch shell {
	case "bash":
		return initBash, nil
	case "zsh":
		return initZsh, nil
	case "fish":
		return initFish, nil
	case "powershell":
		return initPowershell, nil
	default:
		return "", fmt.Errorf("Unsupported shell: %s \n", shell)
	}
}
