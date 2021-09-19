# NATS Streaming/Jetstream Replicator [SJR]

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/snapp-incubator/stan-js-replicator/ci?label=ci&logo=github&style=flat-square)

## Introduction

This project replicates messages from streaming channels to jetstream. but why?
At Snapp when we started using the nats-streaming there was no jetstream for nats,
so streaming was our only solution to have reliable delivery for important events,
but now streaming is going to be deprecated as mentioned [here](https://github.com/nats-io/nats-streaming-server#warning--deprecation-notice-warning) and jetstream
has new features which we will use.

Right now all of our projects use streaming and
we want to have an smooth migration so with this project we can replicate events
on the jetstream side then remove them from streaming side.

## How it works?

It uses the configuration as below to first creates the stream on jetstream and then binds the given topics to it as subjects.
then it starts subscribing on given topics from nats streaming and publishes them into jetstream.
if you want to bind topics into different streams you need to run different instances of sjr.

```yaml
---
monitoring:
  enabled: true
  address: ":8080"
nats:
  url: "nats://127.0.0.1:4222"
logger:
telemetry:
streaming:
  url: "nats://127.0.0.1:4222"
  group: "sjr"
  clientid: "snapp"
channel: "koochooloo"
topics:
  - k.1
  - k.2
stream:
  maxage: "1h"
  storagetype: 1
  replicas: 1
```
