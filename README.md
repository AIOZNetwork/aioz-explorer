## AIOZ Blockchain Explorer

#### Introduction
This project is written to track down most actions done on AIOZ blockchain: blocks, transactions, messages, events,... , enhanced by Cockroachdb.

#### How to start
There are 3 services need to be run:
1. Cockroachdb cluster (for dev, 1 node is enough)
2. GoB node, which is AIOZ node with GoB implemented inside, used to index data
3. Light Client Daemon (LCD) server: serve blockchain services like broadcasting transactions

## Notes
#### Database
1. Mind the user, password (in case of security enabled), host, port, database name, schema,...
2. Check each table properties (index, primary key,...)

#### Replay blocks
There are 2 binary that runs block replay:
- aiozreplay: for aiozd in general
- replay: for explorer's needs

#### Other nodes

1. There are changes since v0.39.0 when you want to replay block:
- Previous version sets `_state.LastHeightValidatorsChanged = 1`
- If `s.LastBlockHeight` causes a valset change,
- We set `s.LastHeightValidatorsChanged = s.LastBlockHeight + 1 + 1`
- Extra +1 due to _**nextValSet**_ delay.
- In `replay/utils.go`, set `_state.LastHeightValidatorsChanged = rollbackBlock1.Header.Height + 1 + 1` for GoB replay
- In `/cmd/aiozreplay/utils.go` set `_state.LastHeightValidatorsChanged = rollbackBlock1.Header.Height + 1 + 1` for `aiozd` replay
- _Not yet testing on low block count with nothing changed in validators set_
