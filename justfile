set shell := ["bash", "-uc"]

envrc := '''
  watch_file **/*.nix
  use_flake .#
'''

default:
  @just --choose \
  --chooser="\
  just --list --list-heading '' | \
  grep -v -e '^\s*default$' -e '^\s*exit$' | \
  sort | \
  { echo 'exit'; cat; } | \
  fzf --no-sort \
  || :"

exit:
  @exit 0

bootstrap:
  echo '{{envrc}}' > .envrc

build:
  go build

test:
  go test -v

update:
  nix flake update
  direnv reload
  gomod2nix generate
