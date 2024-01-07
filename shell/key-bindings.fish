status is-interactive; or exit 0

function __tbmk_save__
    set -l commandlineVal (commandline -t)
    eval tbmk save --command="$commandlineVal"
    commandline -f repaint
end

function __tbmk_search__
    set -l commandlineVal (commandline -t)
    set output (eval tbmk search --query="$commandlineVal")
    commandline -t "$output"
    commandline -f repaint
end

bind \ct __tbmk_save__

bind -k nul __tbmk_search__
