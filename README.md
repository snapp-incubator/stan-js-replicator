# NATS Streaming/Jetstream Replicator

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
