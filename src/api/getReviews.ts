import { ReviewsList } from "../types/Reviews";

// const REVIEWS_API_URL = "http://localhost:8080/api/reviews";

/**
 * Fetches reviews from the server.
 * @returns A promise that resolves to an array of Review objects.
 * @throws {Error} When the fetch operation fails or returns a non-ok status.
 */
export async function getReviews(id: string): Promise<ReviewsList> {
  try {
    const url = `api/reviews?id=${id}`;
    const response = await fetch(url);
    if (!response.ok) {
      throw new Error(`Failed to fetch data: ${response.status} (${response.statusText}) (${response.text()}`);
    }
    const data: ReviewsList = await response.json();
    return data;
  } catch (error) {
    throw new Error(`An error occurred while fetching reviews: ${error instanceof Error ? error.message : String(error)}`);
  }
}
