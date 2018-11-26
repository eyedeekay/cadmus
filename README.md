# Cadmus - An IRC Bot, Logger of channels

From the Greek god [Cadmus](https://en.wikipedia.org/wiki/Cadmus)

> In Greek mythology, Cadmus (/ˈkædməs/; Greek: Κάδμος Kadmos), was the
> founder and first king of Thebes.
> ...
> Cadmus was credited by the ancient Greeks (Herodotus[4] is an example)
> with introducing the original alphabet to the Greeks, who adapted it to
> form their Greek alphabet. 

And so `Cadmus` is an IRC Bot that logs IRC Channels.

## Requirements

Cadmus has no special requirements. Simply invite it to a channel you want
logged and it will keep logs of the channel.

## Installation

```#!bash
$ go get github.com/prologic/cadmus
```

## Getting Started

Simply run `cadmus`:

```#!bash
$ ./cadmus
```

## How it works

- Cadmus will connect to a configured server.
- When Cadmus is invited to a channel; it will immediately join.
- Cadmus will then log all activity on the channel.

## Related Projects

* [eris](https://github.com/prologic/eris) -- a modern IRC Server / Daemon written in Go that has a heavy focus on security and privacy
* [soter](https://github.com/prologic/soter) -- an IRC Bot written in Go that protects IRC Channels by persisting channel modes and topics

## License

Cadmus is licensed under the MIT License.
