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
curl -sL ${BINARY_URL} > ${JABBA_DIR}/bin/jabba && chmod a+x ${JABBA_DIR}/bin/jabba
cat >${JABBA_DIR}/jabba.sh<<-EOF
# https://github.com/shyiko/jabba
# This file is indented to be "sourced" (i.e. `. ~/.jabba/jabba.sh`)

jabba() {
    local fd3=\$(mktemp /tmp/jabba-fd3.XXXXXX)
    (JABBA_SHELL_INTEGRATION=ON ${JABBA_DIR}/bin/jabba "\$@" 3> \${fd3})
    local exit_code=\$?
    eval \$(cat \${fd3})
    rm \${fd3}
    (exit \${exit_code})
}
EOF

SOURCE_JABBA="\n[ -s \"$JABBA_DIR/jabba.sh\" ] && source \"$JABBA_DIR/jabba.sh\""

profile_found=0
files=("$PROFILE" "$HOME/.bashrc" "$HOME/.bash_profile" "$HOME/.zshrc" "$HOME/.profile")
for file in "${files[@]}"
do
    if [ -f "${file}" ]; then
        profile_found=1
        if ! grep -qc '/jabba.sh' "${file}"; then
            echo "Adding source string to ${file}"
            printf "$SOURCE_JABBA\n" >> "${file}"
        else
            echo "Skipped update of ${file} (source string already present)"
        fi
    fi
done

if [ "${profile_found}" == "0" ]; then
  echo "Shell config file(s) not found (tried \$PROFILE, ~/.bashrc, ~/.bash_profile, ~/.zshrc, and ~/.profile)."
  echo "Create one/several of them and run this script again"
  echo "-or-"
  echo "Append the following line(s) to the correct file yourself"
  printf "$SOURCE_JABBA"
  echo
  exit 1
fi

echo ""
echo "Installation completed
(if you have any problems please report them at https://github.com/shyiko/jabba/issue)"

} # this ensures the entire script is downloaded
