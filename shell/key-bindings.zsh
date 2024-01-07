__tbmk_save__() {
  tbmk save -command="${LBUFFER}"
  LBUFFER=""
  local ret=$?
  zle reset-prompt
  return $ret
}

zle     -N     __tbmk_save__
bindkey '\C-t' __tbmk_save__

__tbmk_search__() {
  output=$(tbmk search -query="${LBUFFER}")
  LBUFFER="${output}"
  local ret=$?
  zle reset-prompt
  return $ret
}

zle     -N     __tbmk_search__
bindkey '\C-@' __tbmk_search__
