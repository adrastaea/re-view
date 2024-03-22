import React, { useEffect, useState } from "react";
import { getReviews } from "../api/getReviews";
import { Review } from "../types/Reviews";
import ReviewCard from "./ReviewCard";

const ReviewsColumn: React.FC = () => {
  const [reviews, setReviews] = useState<Review[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    const fetchReviews = async () => {
      try {
        setLoading(true);
        const response = await getReviews();
        setReviews(response.Reviews);
        setLoading(false);
      } catch (error) {
        setError("Failed to load reviews.");
        setLoading(false);
      }
    };

    fetchReviews();
  }, []);

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;

  console.log(reviews);
  return (
    <div>
      {reviews.map((review) => (
        <ReviewCard key={review.Id} review={review} />
      ))}
    </div>
  );
};

export default ReviewsColumn;
