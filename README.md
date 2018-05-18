# p2p2p
peer2peer2peer - Naive peer2peer networking with basic API

# Operation

The way p2p2p's networking works is that a node without the `--bootnodes` flag is a network virgin. It's pool is open and so it listens indefinitely. It is unable to message anyone.

If a new node contacts it, as a member of its `--bootnodes` flag, the two swap their phonebooks, expanding both phonebooks.

The rest of the nodes on the network are not immediately informed of this new node. But any node, if given a node name not already in its pool, will ask its network if anyone has it and if it is returned, will incorporate and continue. This is modeled similarly to ARP requests, in a fashion.

#### In Dev
    - `go run main.go run --help`
    - `go run main.go run -n node00 -p 3032`
    - `go run main.go run -n node01 -p 3033`
    - `go run main.go run -p 3030 -n node02 -b localhost:3032,localhost:3033`
#### Binary
    - `./p2p2p run ...`

# API

- `/health`
  - Check to see if the API is serving. Also provides known node count.
- `/whisper/{name}?message={message}`
  - Sends {message} to {name}