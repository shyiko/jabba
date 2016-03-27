#!/bin/bash -e
{ # this ensures the entire script is downloaded

JABBA_DIR=${JABBA_DIR:-$HOME/.jabba}
JABBA_VERSION=${JABBA_VERSION:-latest}

if [ "$JABBA_VERSION" == "latest" ]; then
    # resolving "latest" to an actual tag
    JABBA_VERSION=$(curl -o- https://api.github.com/repos/shyiko/jabba/releases/latest | grep 'tag_name' | cut -d\" -f4)
fi

case "$OSTYPE" in
    darwin*)
    BINARY_URL=https://github.com/shyiko/jabba/releases/download/${JABBA_VERSION}/jabba-${JABBA_VERSION}-darwin-amd64
    ;;
    linux*)
    if [ `getconf LONG_BIT` = "64" ]; then OSARCH=amd64; else OSARCH=386; fi
    BINARY_URL=https://github.com/shyiko/jabba/releases/download/${JABBA_VERSION}/jabba-${JABBA_VERSION}-linux-${OSARCH}
    ;;
    *)
    echo "Unsupported OS $OSTYPE. If you believe this is an error -
please create a ticket at https://github.com/shyiko/jabba/issue. Thank you"
    exit 1
    ;;
esac

echo "Installing v$JABBA_VERSION..."
echo

mkdir -p ${JABBA_DIR}/bin

if [ "$JABBA_MAKE_INSTALL" == "true" ]; then
    cp jabba ${JABBA_DIR}/bin
else
    curl -sL ${BINARY_URL} > ${JABBA_DIR}/bin/jabba && chmod a+x ${JABBA_DIR}/bin/jabba
fi

{
echo "# https://github.com/shyiko/jabba"
echo "# This file is indented to be \"sourced\" (i.e. \". ~/.jabba/jabba.sh\")"
echo ""
echo "jabba() {"
echo "    local fd3=\$(mktemp /tmp/jabba-fd3.XXXXXX)"
echo "    (JABBA_SHELL_INTEGRATION=ON ${JABBA_DIR}/bin/jabba \"\$@\" 3> \${fd3})"
echo "    local exit_code=\$?"
echo "    eval \$(cat \${fd3})"
echo "    rm \${fd3}"
echo "    (exit \${exit_code})"
echo "}"
echo ""
echo "[ ! -z \"\$(jabba alias default)\" ] && jabba use default"
} > ${JABBA_DIR}/jabba.sh

SOURCE_JABBA="\n[ -s \"$JABBA_DIR/jabba.sh\" ] && source \"$JABBA_DIR/jabba.sh\""

files=("$HOME/.bashrc" "$HOME/.bash_profile" "$HOME/.profile")
for file in "${files[@]}"
do
    touch ${file}
    if ! grep -qc '/jabba.sh' "${file}"; then
        echo "Adding source string to ${file}"
        printf "$SOURCE_JABBA\n" >> "${file}"
    else
        echo "Skipped update of ${file} (source string already present)"
    fi
done

if [ -f "$(which zsh 2>/dev/null)" ]; then
    file="$HOME/.zshrc"
    touch ${file}
    if ! grep -qc '/jabba.sh' "${file}"; then
        echo "Adding source string to ${file}"
        printf "$SOURCE_JABBA\n" >> "${file}"
    else
        echo "Skipped update of ${file} (source string already present)"
    fi
fi

{
echo "# https://github.com/shyiko/jabba"
echo "# This file is indented to be \"sourced\" (i.e. \". ~/.jabba/jabba.fish\")"
echo ""
echo "function jabba"
echo "    set fd3 (mktemp /tmp/jabba-fd3.XXXXXX)"
echo "    env JABBA_SHELL_INTEGRATION=ON ${JABBA_DIR}/bin/jabba \$argv 3> \$fd3"
echo "    set exit_code \$status"
echo "    eval (cat \$fd3 | sed \"s/^export/set -x/g\" | sed \"s/^unset/set -e/g\" | tr '=:' ' ' | tr '\\\\n' ';')"
echo "    rm \$fd3"
echo "    return \$exit_code"
echo "end"
echo ""
echo "[ ! -z (echo (jabba alias default)) ]; and jabba use default"
} > ${JABBA_DIR}/jabba.fish

FISH_SOURCE_JABBA="\n[ -s \"$JABBA_DIR/jabba.fish\" ]; and source \"$JABBA_DIR/jabba.fish\""

if [ -f "$(which fish 2>/dev/null)" ]; then
    file="$HOME/.config/fish/config.fish"
    mkdir -p $(dirname ${file})
    touch ${file}
    if ! grep -qc '/jabba.fish' "${file}"; then
        echo "Adding source string to ${file}"
        printf "$FISH_SOURCE_JABBA\n" >> "${file}"
    else
        echo "Skipped update of ${file} (source string already present)"
    fi
fi

echo ""
echo "Installation completed
(if you have any problems please report them at https://github.com/shyiko/jabba/issue)"

} # this ensures the entire script is downloaded
