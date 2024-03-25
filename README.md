# Re:View

[![Link](https://img.shields.io/badge/Link-View%20Live-blue)](https://re-view.vercel.app/)

A simple dashboard of recent reviews from the current top 10 most popular apps on the Apple App Store.

## Stack

- React
- Vite
- Tailwind CSS
- Go
- Vercel
  - Frontend Hosting
  - Go Serverless Functions
  - PostgresSQL DB

## Running your own instance

### Prerequisites

- Go
- Node.js

### Deployment

1. Clone the repository
2. Run `npm install` in the root directory
3. Run `npm run dev` to start the frontend
4. Open a new terminal and navigate to the `service` directory
5. Run `go run .` to start the server
6. Open <http://localhost:3000> in your browser
