

list:
    just --list

run:
    go run .

# https://github.com/charmbracelet/bubbletea#debugging
debugger:
    #!/usr/bin/env fish
    dlv debug --listen=:2345 --headless --api-version=2 .

debugger-kill:
    #!/usr/bin/env fish
    kill (ps aux | grep -v grep | grep 'dlv debug' | awk '{print $2}')

git-push m="updates":
    git status
    git add . || true
    git commit -m "{{ m }}" || true
    git push origin master
    git push origin master --tags
