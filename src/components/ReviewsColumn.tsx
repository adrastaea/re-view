import React, { useEffect, useState } from "react";
import { getReviews } from "../api/getReviews";
import { Review } from "../types/Reviews";
import ReviewCard from "./ReviewCard";
import reactLogo from "../assets/react.svg";

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

  if (loading)
    return (
      <div className="flex items-center p-8 text-3xl">
        <img
          src={reactLogo}
          alt="React Logo"
          className="mr-4 h-12 w-12 animate-spin"
        />
        Loading...
      </div>
    );
  if (error) return <div className="flex p-8 text-3xl">Error: {error}</div>;

  console.log(reviews);
  return (
    <div className="container flex flex-col gap-4 p-4">
      {reviews.map((review) => (
        <ReviewCard key={review.Id} review={review} />
      ))}
    </div>
  );
};

export default ReviewsColumn;
