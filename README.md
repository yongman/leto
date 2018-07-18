# Leto

## 0. What is Leto mean?

In Greek mythology, [Leto](https://en.wikipedia.org/wiki/Leto) (/ˈliːtoʊ/) is a daughter of the Titans Coeus and Phoebe, the sister of Asteria.

## 1. What is Leto?

Leto is another reference example use of [Hashicorp Raft](https://github.com/hashicorp/raft). The API is [redis protocol](https://redis.io/topics/protocol) compatiable.

[Raft](https://raft.github.io/)  is a consensus algorithm that is designed to be easy to understand. It's equivalent to Paxos in fault-tolerance and performance. The difference is that it's decomposed into relatively independent subproblems, and it cleanly addresses all major pieces needed for practical systems. We hope Raft will make consensus available to a wider audience, and that this wider audience will be able to develop a variety of higher quality consensus-based systems than are available today.

## 2. Why do this?

You can have better comprehension about how `raft protocal` works if you use it. This helps me a lot.


## 3. Run sample

**3.1 show helps**
```
bin/leto -h
Usage of bin/leto:
  -id string
        node id
  -join string
        join to already exist cluster
  -listen string
        server listen address (default ":5379")
  -raftbind string
        raft bus transport bind address (default ":15379")
  -raftdir string
        raft data directory (default "./")
```

**3.2 Start first node**

```
bin/leto -id id1 -raftdir ./id1
```
the first node will be listen user request and node join request in port `5379`, and use port `15379` for raft transport.

**3.3 Start second node**

```
bin/leto -id id2 -raftdir ./id2 -listen ":6379" -raftbind ":16379" -join "127.0.0.1:5379"
```

**3.4 Start third node**

```
bin/leto -id id3 -raftdir ./id3 -listen ":7379" -raftbind ":17379" -join "127.0.0.1:5379"
```

**3.5 Test**

Requst first node
```
redis-cli -p 5379
127.0.0.1:5379> set a b
OK
127.0.0.1:5379> get a
b
127.0.0.1:5379>
```

Write to second node, data has been replicated to this node. And it will return `not leader error` if write to it.

```
redis-cli -p 6379
127.0.0.1:6379> get a
b
127.0.0.1:6379> set a b
(error) not leader
127.0.0.1:6379>
```

Now, we  `shutdown` the first node, the second node voted to be leader.
```
redis-cli -p 6379
127.0.0.1:6379> get a
b
127.0.0.1:6379> set a b
OK
127.0.0.1:6379>
```

## 4. Support commands

- GET
- SET
- JOIN (communicate with peer when start node)
- LEAVE (remove dead node from raft group)
- PING
