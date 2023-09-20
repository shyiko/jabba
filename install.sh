#!/bin/bash -e
{ # this ensures the entire script is downloaded

JABBA_HOME=${JABBA_HOME:-$JABBA_DIR} # JABBA_DIR is here for backward-compatibility
JABBA_VERSION=${JABBA_VERSION:-latest}

if [ "$JABBA_HOME" == "" ] || [ "$JABBA_HOME" == "$HOME/.jabba" ]; then
    JABBA_HOME=$HOME/.jabba
    JABBA_HOME_TO_EXPORT=\$HOME/.jabba
else
    JABBA_HOME_TO_EXPORT=$JABBA_HOME
fi

# Install location customizations
JABBA_BIN=${JABBA_BIN:-"$JABBA_HOME/bin"}
JABBA_BIN_TO_EXPORT=${JABBA_BIN_TO_EXPORT:-"$JABBA_HOME_TO_EXPORT/bin"}
JABBA_SHARE=${JABBA_SHARE:-$JABBA_HOME}

has_command() {
    if ! command -v "$1" > /dev/null 2>&1
    then echo 1;
    else echo 0;
    fi
}

SKIP_RC=
while :; do
    case "$1" in
    --skip-rc) # skip rc file scripts
        SKIP_RC="1"
        ;;
    *)
        break
        ;;
    esac
    shift
done

# curl looks for HTTPS_PROXY while wget for https_proxy
https_proxy=${https_proxy:-$HTTPS_PROXY}
HTTPS_PROXY=${HTTPS_PROXY:-$https_proxy}

if [ "$JABBA_GET" == "" ]; then
    if [ 0 -eq "$(has_command curl)" ]; then
        JABBA_GET="curl -sL"
    elif [ 0 -eq "$(has_command wget)" ]; then
        JABBA_GET="wget -qO-"
    else
        echo "[ERROR] This script needs wget or curl to be installed."
        exit 1
    fi
fi

if [ "$JABBA_VERSION" == "latest" ]; then
    # resolving "latest" to an actual tag
    JABBA_VERSION=$($JABBA_GET https://Jabba-Team.github.io/jabba/latest)
fi

# http://semver.org/spec/v2.0.0.html
if [[ ! "$JABBA_VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.+-]+)?$ ]]; then
    echo "'$JABBA_VERSION' is not a valid version."
    exit 1
fi

case "$OSTYPE" in
    darwin*)
    case "$(uname -m)" in
      x86_64*)
        OSARCH=amd64
      ;;
      arm64*)
        OSARCH=arm64
      ;;
    esac
    printf "Downloading for mac for arch %s" "${OSARCH}"
    BINARY_URL=https://github.com/Jabba-Team/jabba/releases/download/${JABBA_VERSION}/jabba-${JABBA_VERSION}-darwin-${OSARCH}
    ;;
    linux*)
    case "$(uname -m)" in
        arm*|aarch*)
        if [ "$(getconf LONG_BIT)" = "64" ]; then OSARCH=arm64; else OSARCH=arm; fi
        ;;
        s390x)
        OS_ARCH=s390x
        echo "OS_ARCH='$OS_ARCH' is not a valid architecture at this point."
        exit 1
        ;;
        s390*)
        if [ "$(getconf LONG_BIT)" = "64" ]; then OSARCH=s390x; else OSARCH=s390; fi
        echo "OS_ARCH='$OS_ARCH' is not a valid architecture at this point."
        exit 1
        ;;
        powerpc*)
        OS_ARCH=powerpc
        echo "OS_ARCH='$OS_ARCH' is not a valid architecture at this point."
        exit 1
        ;;
        ppc64*)
        OS_ARCH=ppc64le
        echo "OS_ARCH='$OS_ARCH' is not a valid architecture at this point."
        exit 1
        ;;
        *)
        if [ "$(getconf LONG_BIT)" = "64" ]; then OSARCH=amd64; else OSARCH=386; fi
        ;;
    esac
    BINARY_URL="https://github.com/Jabba-Team/jabba/releases/download/${JABBA_VERSION}/jabba-${JABBA_VERSION}-linux-${OSARCH}"
    ;;
    cygwin*|msys*)
    OS_ARCH=$(echo 'echo %PROCESSOR_ARCHITECTURE% & exit' | cmd | tail -n 1 | xargs) # xargs used to trim whitespace
    if [ "$OS_ARCH" == "AMD64" ]; then
        BINARY_URL="https://github.com/Jabba-Team/jabba/releases/download/${JABBA_VERSION}/jabba-${JABBA_VERSION}-windows-amd64.exe"
    elif [ "$OS_ARCH" == "x86" ]; then
        BINARY_URL="https://github.com/Jabba-Team/jabba/releases/download/${JABBA_VERSION}/jabba-${JABBA_VERSION}-windows-386.exe"
    else
        echo "OS_ARCH='$OS_ARCH' is not a valid architecture at this point."
        exit 1
    fi
    ;;
    *)
    echo "Unsupported OS $OSTYPE. If you believe this is an error -
please create a ticket at https://github.com/Jabba-Team/jabba/issues."
    exit 1
    ;;
esac

echo "Installing v$JABBA_VERSION..."
echo

if [ ! -f "${JABBA_BIN}/jabba" ]; then
    JABBA_SELF_DESTRUCT_AFTER_COMMAND="true"
fi

mkdir -p "${JABBA_SHARE}"
mkdir -p "${JABBA_BIN}"

if [ "$JABBA_MAKE_INSTALL" == "true" ]; then
    cp jabba "${JABBA_BIN}"
else
    $JABBA_GET "${BINARY_URL}" > "${JABBA_BIN}/jabba" && chmod a+x "${JABBA_BIN}/jabba"
fi

if ! "${JABBA_BIN}/jabba" --version &>/dev/null; then
    echo "${JABBA_BIN}/jabba does not appear to be a valid binary.

Check your Internet connection / proxy settings and try again.
If the problem persists - please create a ticket at https://github.com/Jabba-Team/jabba/issues."
    exit 1
fi

if [ "$JABBA_COMMAND" != "" ]; then
    "${JABBA_BIN}/jabba" "$JABBA_COMMAND"
    if [ "$JABBA_SELF_DESTRUCT_AFTER_COMMAND" == "true" ]; then
        rm -f "${JABBA_BIN}/jabba"
        rmdir "${JABBA_BIN}"
        exit 0
    fi
fi

sed -e "s=\$JABBA_HOME_TO_EXPORT=$JABBA_HOME_TO_EXPORT=g" -e "s=\$JABBA_BIN_TO_EXPORT=$JABBA_BIN_TO_EXPORT=g" jabba.sh > "${JABBA_SHARE}/jabba.sh"

SOURCE_JABBA="\n[ -s \"$JABBA_SHARE/jabba.sh\" ] && source \"$JABBA_SHARE/jabba.sh\""

if [ ! "$SKIP_RC" ]; then
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
        touch "${file}"
        if ! grep -qc '/jabba.sh' "${file}"; then
            echo "Adding source string to ${file}"
            printf "%s\n" "$SOURCE_JABBA" >> "${file}"
        else
            echo "Skipped update of ${file} (source string already present)"
        fi
    done

    if [ -f "$(which zsh 2>/dev/null)" ]; then
        file="$HOME/.zshrc"
        touch "${file}"
        if ! grep -qc '/jabba.sh' "${file}"; then
            echo "Adding source string to ${file}"
            printf "%s\n" "$SOURCE_JABBA" >> "${file}"
        else
            echo "Skipped update of ${file} (source string already present)"
        fi
    fi
fi

sed -e "s=\$JABBA_HOME_TO_EXPORT=$JABBA_HOME_TO_EXPORT=g" -e "s=\$JABBA_BIN_TO_EXPORT=$JABBA_BIN_TO_EXPORT=g" jabba.fish > "${JABBA_SHARE}/jabba.fish"

FISH_SOURCE_JABBA="\n[ -s \"$JABBA_SHARE/jabba.fish\" ]; and source \"$JABBA_SHARE/jabba.fish\""

if [ -f "$(which fish 2>/dev/null)" ]; then
    file="$HOME/.config/fish/config.fish"
    mkdir -p "$(dirname "${file}")"
    touch "${file}"
    if ! grep -qc '/jabba.fish' "${file}"; then
        echo "Adding source string to ${file}"
        printf "%s\n" "$FISH_SOURCE_JABBA" >> "${file}"
    else
        echo "Skipped update of ${file} (source string already present)"
    fi
fi

echo ""
echo "Installation completed
(if you have any problems please report them at https://github.com/Jabba-Team/jabba/issues)"

} # this ensures the entire script is downloaded
