#!/bin/sh -eu

. /etc/os-release

SELF=$(readlink -f "$0")
PROJECT_ROOT=$(dirname "$SELF")

install_package() {
  /usr/bin/apt-get update -qq
  echo ">> Installing ${1}"
  /usr/bin/apt-get install -y --no-install-recommends "$1" || true
}

cd "$PROJECT_ROOT"

if test $# -gt 0; then
  PUPPET_VERSION="$1"
fi

if test -n "${PUPPET_VERSION:-}"; then
  if ! test -s "puppet${PUPPET_VERSION}-release-${VERSION_CODENAME}.deb"; then
    if ! test -s /usr/bin/curl; then
      install_package curl
    fi

    echo ">> Downloading https://apt.puppetlabs.com/puppet${PUPPET_VERSION}-release-${VERSION_CODENAME}.deb"
    /usr/bin/curl -sLO "https://apt.puppetlabs.com/puppet${PUPPET_VERSION}-release-${VERSION_CODENAME}.deb"
  fi

  echo ">> Installing Puppetlabs APT configuration"
  /usr/bin/apt install "./puppet${PUPPET_VERSION}-release-${VERSION_CODENAME}.deb"
  install_package puppet-agent
else
  install_package puppet
fi

export PATH="${PATH}:/opt/puppetlabs/puppet/bin"

echo ">> Running puppet with input.pp"
puppet apply --test "input.pp" || true

PUPPET_VERSION=$(puppet --version)
PUPPET_REPORT_FILE=$(puppet config print lastrunreport)

cp -v "$PUPPET_REPORT_FILE" "last_run_report-${PUPPET_VERSION}.yaml"

