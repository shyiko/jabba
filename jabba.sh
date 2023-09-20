# https://github.com/Jabba-Team/jabba
# This file is intended to be "sourced" (i.e. ". ~/.jabba/jabba.sh")

export JABBA_HOME="$JABBA_HOME_TO_EXPORT"

jabba() {
    local fd3
    fd3=$(mktemp /tmp/jabba-fd3.XXXXXX)
    (JABBA_SHELL_INTEGRATION=ON "$JABBA_BIN_TO_EXPORT/jabba" "$@" 3>| "${fd3}")
    local exit_code=$?
    eval "$(cat "${fd3}")"
    command rm -f "${fd3}"
    return ${exit_code}
}

[ -n "$(jabba alias default)" ] && jabba use default
