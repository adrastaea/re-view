import { ReviewsList } from "../types/Reviews";

// const REVIEWS_API_URL = "http://localhost:8080/api/reviews";

/**
 * Fetches reviews from the server.
 * @returns A promise that resolves to an array of Review objects.
 * @throws {Error} When the fetch operation fails or returns a non-ok status.
 */
export async function getReviews(id: string): Promise<ReviewsList> {
  try {
    console.log("getReviews: ", id);
    const url = `api/reviews`;
    console.log(url);
    const response = await fetch(url);
    if (!response.ok) {
      throw new Error(`Failed to fetch data: ${response.status} (${response.statusText})`);
    }
    const data: ReviewsList = await response.json();
    return data;
  } catch (error) {
    throw new Error(`An error occurred while fetching reviews: ${error instanceof Error ? error.message : String(error)}`);
  }
}
