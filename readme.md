

### Checks codebases for duplicate err-ids

```bash
cd "$GOPATH/src/github.com/oresoftware/err-id-checker"
go install
cd '<project>' && err-id-checker
```

### TODO

1. check codebases for err-ids with different regex patterns:

currently it only checks for this pattern:
```
ErrId: "<uuid-v4>"
```

but it would be nice if it checked for uuids with this pattern:

```
"▲<uuid-v4>"
```

aka, use some unicode char to prefix the uuid, and then can grep for that.
