# 💬 Chatty

Welcome to **chatty** — a terminal-based, real-time chat app built for learning and playing around with Golang’s concurrency and WebSocket magic.

## What is chatty?

Chatty is my personal sandbox for exploring real-time data streaming using Go’s powerful **goroutines** and **channels**, with WebSocket keeping the conversations flowing smoothly. To make the terminal UI a bit more fun and interactive, I’m using the **bubble tea** library to keep things fresh.

Under the hood, chatty leverages the trusted **Gorilla WebSocket** library. When you join (currently) a single chat room, the server creates a client layer as a hands-on middleman between your WebSocket connection and the server itself. Each client spins up two goroutines — one to listen for incoming messages and another to send outgoing ones — making sure conversations flow without blocking or lag.

## Why chatty?

Because sometimes breaking things in code beats binge-watching shows. This project lets me dive deep into concurrency, networking, and terminal UI — and it’s purely a **personal learning project**. So, if it misbehaves, it’s just me fumbling through new territory.

## How does it work?

- The **server** manages a chat room, handling WebSocket connections and coordinating message broadcasts.
- Upon joining, each user is represented as a **client** that bridges the WebSocket connection and the server.
- Two goroutines run per client: one constantly reads messages from your connection, and the other writes outgoing messages, keeping things snappy and responsive.
- The **client app**, running in your terminal, uses the bubble tea library for a smooth UI, sending and receiving messages in real-time.

## Features

- Real-time messaging through WebSockets.
- Terminal UI designed with bubble tea for a smooth interaction.
- Concurrency handled by goroutines and channels — Go’s secret sauce.
- No login fuss; just jump in and say hi.

## What’s still cooking?

- No support for multiple chat rooms (but that’s on the roadmap).
- Client UI it's still a work in progress.
- Error handling is still a work in progress.

## Get started

Clone the repo, run the server, then join the chat from your terminal:

```sh
git clone https://github.com/gargalloeric/chatty.git

cd chatty

make run/api # start the server

make run/client # start the client
```

## Want to play around?

Feel free to fork, tweak, or completely rebuild. This is all about learning, experimenting, and having a bit of fun with Go’s real-time capabilities.