__tbmk_save__() {
  offset=${READLINE_POINT}
  READLINE_POINT=0

  command=$(echo "${READLINE_LINE:0}" | grep -oE '(.)+$')
  ./tbmk save -command="$command"
}

__tbmk_search__() {
  offset=${READLINE_POINT}
  READLINE_POINT=0

  local output
  query=$(echo "${READLINE_LINE:0}" | grep -oE '(.)+$')
  output=$(./tbmk search -query="$query")
  READLINE_LINE=${output#*$'\t'}
  if [[ -z "$READLINE_POINT" ]]; then
    echo "$READLINE_LINE"
  else
    READLINE_POINT=0x7fffffff
  fi
}

bind -x '"\C-t":__tbmk_save__'
bind -x '"\C-@":__tbmk_search__'
