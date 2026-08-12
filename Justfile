

list:
    just --list

run:
    go run .

debugger:
    dlv debug --listen=:2345 --headless --api-version=2 .

git-push m="updates":
    git status
    git add . || true
    git commit -m "{{ m }}" || true
    git push origin master
    git push origin master --tags
