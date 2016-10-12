#!/bin/bash -e
{ # this ensures the entire script is downloaded

JABBA_DIR=${JABBA_DIR:-$HOME/.jabba}
JABBA_VERSION=${JABBA_VERSION:-latest}

# curl looks for HTTPS_PROXY while wget for https_proxy
https_proxy=${https_proxy:-$HTTPS_PROXY}
HTTPS_PROXY=${HTTPS_PROXY:-$https_proxy}

if [ "$JABBA_GET" == "" ]; then
    if [ -f "$(which curl 2>/dev/null)" ]; then
        JABBA_GET="curl -sL"
    else
        JABBA_GET="wget -qO-"
    fi
fi

if [ "$JABBA_VERSION" == "latest" ]; then
    # resolving "latest" to an actual tag
    JABBA_VERSION=$($JABBA_GET https://shyiko.github.com/jabba/latest)
fi

# http://semver.org/spec/v2.0.0.html
if [[ ! "$JABBA_VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.+-]+)?$ ]]; then
    echo "'$JABBA_VERSION' is not a valid version."
    exit 1
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
please create a ticket at https://github.com/shyiko/jabba/issue."
    exit 1
    ;;
esac

echo "Installing v$JABBA_VERSION..."
echo

mkdir -p ${JABBA_DIR}/bin

if [ "$JABBA_MAKE_INSTALL" == "true" ]; then
    cp jabba ${JABBA_DIR}/bin
else
    $JABBA_GET ${BINARY_URL} > ${JABBA_DIR}/bin/jabba && chmod a+x ${JABBA_DIR}/bin/jabba
fi

if ! ${JABBA_DIR}/bin/jabba --version &>/dev/null; then
    echo "${JABBA_DIR}/bin/jabba does not appear to be a valid binary.

Check your Internet connection / proxy settings and try again.
If the problem persists - please create a ticket at https://github.com/shyiko/jabba/issue."
    exit 1
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
echo "    rm -f \${fd3}"
echo "    (exit \${exit_code})"
echo "}"
echo ""
echo "if [ ! -z \"\$(jabba alias default)\" ]; then"
echo "    jabba use default"
echo "fi"
} > ${JABBA_DIR}/jabba.sh

SOURCE_JABBA="\n[ -s \"$JABBA_DIR/jabba.sh\" ] && source \"$JABBA_DIR/jabba.sh\""

files=("$HOME/.bashrc")

if [ -f "$HOME/.bash_profile" ]; then
    files+=("$HOME/.bash_profile")
elif [ -f "$HOME/.bash_login" ]; then
    files+=("$HOME/.bash_login")
elif [ -f "$HOME/.profile" ]; then
    files+=("$HOME/.profile")
else
    files+=("$HOME/.bash_profile")
fi

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
echo "    eval (cat \$fd3 | sed \"s/^export/set -x/g\" | sed \"s/^unset/set -e/g\" | tr '=' ' ' | sed \"s/:/\\\" \\\"/g\" | tr '\\\\n' ';')"
echo "    rm -f \$fd3"
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
