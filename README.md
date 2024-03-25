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

- Vercel Account

### Deployment

1. Log in to Vercel
2. Select "Add New...", choose "Project" from the drop down, and choose "Import Third-Party Git Repository"
3. Enter the URL of this repository
4. Accept default settings for now and click "Deploy". This deployment will fail until you complete the settings changes in the next steps.
5. Navigate to your project page, go to the Storage tab, and select Create Database.
6. Create a new PostgreSQL database.
7. Connect the database to the project. The environment variables will be automatically added to the project.
8. Navigate to the Settings tab of your project, and in "General" settings, change the "Node.js version" to "18.x". (This solves a compatibility issue with the serverless functions)
9. You are now ready to deploy your instance. Go back to "Deployments".
10. Click the "..." button on the deployment that failed, and select "Redeploy".
11. Your instance should now be up and running. Click the domain link to view your instance.

### Development

1. Clone the repository
2. Install dependencies

   ```bash
   npm install
   ```

3. Link to vercel serverless functions

   ```bash
   vercel link
   ```

4. Start the development server

   ```bash
   vercek dev
   ```

5. Open the local server that vercel link returned in your terminal in your browser.
