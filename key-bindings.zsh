__tbmk_save__() {
  $TBMK save -command="${LBUFFER}"
  LBUFFER=""
  local ret=$?
  zle reset-prompt
  return $ret
}

zle     -N     __tbmk_save__
bindkey '\C-t' __tbmk_save__

__tbmk_search__() {
  output=$($TBMK search -query="${LBUFFER}")
  LBUFFER="${output}"
  local ret=$?
  zle reset-prompt
  return $ret
}

zle     -N     __tbmk_search__
bindkey '\C-@' __tbmk_search__
