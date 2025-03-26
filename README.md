# Real-Time Quiz Game

A real-time quiz game where players compete to answer questions first. The game supports custom questions and real-time matches, making it interactive and engaging.

## Features

- **Real-Time Matches**: Compete with other players in real-time to answer questions first.
- **Custom Questions**: Create and add your own questions to the game.
- **Leaderboard**: Track the highest scores and see who’s leading.
- **Responsive Design**: Play on any device with a clean and responsive interface.

## Technologies Used

- **Frontend**: React, TS, Shadcn
- **Backend**: golang, fasthttp, rabbitmq
- **Database**: Postgre

## Getting Started

Follow these instructions to set up and run the game locally.

### Prerequisites

- **Node.js**
- **golang**
- **Git**

### Installation

1. **Clone the repository**
2. rename these files. remove `.copy`. use your own enviroment configuration:
- "api/application/env.json.copy"
- "api/hub/env.json.copy"
- "site/env.json.copy"
3. go to `api/hub` and type `go run .`
4. go to `api/application` and type `go run .`
5. go to `site` and type `npm run dev`


## Contributing

This is project written purely for learning and i do not see a point in contributing.

## License

This project is doesn't have any license
