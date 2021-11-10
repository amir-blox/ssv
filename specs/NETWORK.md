# SSV Specification - Network

**WIP**

    
## Fundamentals

### Stack

**TODO**

libp2p + discv5 ... 

### Transport

**TODO**

Supported protocols - TCP/UDP
Encryption
Multiplexing

### Identity

There are two keys for each peer in the network:

##### Network Key
`Network Key` is used to create peer ID. \
All messages from a peer are signed using this key and verified by other peers with the corresponding public key. \
Unless provided in configuration (`NetworkPrivateKey` / `NETWORK_PRIVATE_KEY`), the key will be generated and saved locally for future use. 

##### Operator Key
`Operator Key` is used for decryption of shares keys that are used for signing consensus messages and duties. \
Note that an operator won't be functional in case the key was lost.


### Messaging

Messages in the network are being transported p2p with one of the following methods:

#### Streams

Libp2p allows to create a bidirectional stream between two peers and implement a wire messaging protocol. \
See more information in [IPFS specs > communication-model - streams](https://ipfs.io/ipfs/QmVqNrDfr2dxzQUo4VN3zhG4NV78uYFmRpgSktWDc2eeh2/specs/7-properties/#71-communication-model---streams).

Streams are used in cases where message audience is a single peer.

#### PubSub

PubSub is used as an infrastructure for broadcasting messages among a group (AKA subnet) of operator nodes.

GossipSub ([v1.1](https://github.com/libp2p/specs/blob/master/pubsub/gossipsub/gossipsub-v1.1.md)) is the pubsub protocol used in SSV. \
In short, each node save metadata regards topic subscription of other peers in the network. \
With that information, messages are propagated to the most relevant peers (subscribed or neighbors of a subscribed peer) and therefore reduce the overall traffic.

## Protocols

Network interaction is achieved by using the following protocols:

### 1. Consensus

**TODO**
- IBFT/QBFT consensus
- state/decided propagation

#### Message Structure

The basic message structure includes the following fields:

```protobuf
syntax = "proto3";

// Message represents the object that is being passed around the network
message Message {
  // type is the IBFT state / stage
  RoundState type   = 1;
  // round is the current round where the message was sent
  uint64 round      = 2;
  // lambda is the message identifier
  bytes lambda      = 3;
  // sequence number is an incremental number for each instance, much like a block number would be in a blockchain
  uint64 seq_number = 4;
  // value holds the message data in bytes
  bytes value       = 5;
}

// RoundState is the available types of IBFT state / stage
enum RoundState {
  // NotStarted is when no instance has started yet
  NotStarted  = 0;
  // PrePrepare is the first stage in IBFT
  PrePrepare  = 1;
  // Prepare is the second stage in IBFT
  Prepare     = 2;
  // Commit is when an instance receives a qualified quorum of prepare msgs, then sends a commit msg
  Commit      = 3;
  // ChangeRound is sent upon round change
  ChangeRound = 4;
  // Decided is when an instance receives a qualified quorum of commit msgs
  Decided     = 5;
  // Stopped is the state of an instance that stopped running
  Stopped     = 6;
}
```

`SignedMessage` is the wrapping object that adds a signature and the corresponding singers:

```protobuf
syntax = "proto3";
import "gogo.proto";

// SignedMessage is a wrapper on top of Message for supporting signatures
message SignedMessage{
  // message is the raw message to sign
  Message message = 1 [(gogoproto.nullable) = false];
  // signature is a signature of the message
  bytes signature = 2 [(gogoproto.nullable) = false];
  // signer_ids are the IDs of the signing operators
  repeated uint64 signer_ids = 3;
}
```

**NOTE** all pubsub messages are being wrapped by libp2p's message structure

#### Topics/Subnets

**TODO**

- topics
  - per validator
  - subnets fixed vs. dynamic number of subnets


### 2. History Sync

**TODO**

- why streams
- rate limiting?

#### Stream Protocols

**TODO**

##### Heights Decided

`/sync/highest_decided/0.0.1`

##### Decided By Range 

`/sync/decided_by_range/0.0.1`

##### Last Change Round

`/sync/last_change_round/0.0.1`

## Networking

### Discovery

**TODO**

#### Alternatives

**TODO**

- Kademlia DHT

### Forks

**TODO**

## Major Decisions

**TODO**

## Open points

...


- stack (libp2p, discv5)
- configuration
- topics
  - per validator
  - subnets fixed vs. dynamic number of subnets
- discovery
  - discv5 vs. Kademlia DHT
  - ENR structure
- messaging
  - transport, encryption, multiplexing
  - authentication (peerID correlation to operator key)
  - heartbeat?
  - ibft/consensus
    - pubsub vs. streams
    - protocol
  - history sync
    - why streams
    - rate limiting?
    - protocol
  - decided
    - on main topic - why?