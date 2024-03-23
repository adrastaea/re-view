import { AppsList } from "../types/Apps";

// const REVIEWS_API_URL = "http://localhost:8080/api/topApps";

/**
 * Fetches reviews from the server.
 * @returns A promise that resolves to an array of App objects.
 * @throws {Error} When the fetch operation fails or returns a non-ok status.
 */
export async function getTopApps(): Promise<AppsList> {
  try {
    const response = await fetch('api/topApps');
    if (!response.ok) {
      throw new Error(`Failed to fetch data: ${response.status} (${response.statusText})`);
    }
    const data: AppsList = await response.json();
    return data;
  } catch (error) {
    throw new Error(`An error occurred while fetching reviews: ${error instanceof Error ? error.message : String(error)}`);
  }
}
