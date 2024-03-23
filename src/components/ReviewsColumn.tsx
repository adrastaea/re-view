import React, { useEffect, useState } from "react";
import { getReviews } from "../api/getReviews";
import { Review } from "../types/Reviews";
import ReviewCard from "./ReviewCard";
import reactLogo from "../assets/react.svg";
import { App } from "../types/Apps";

const ReviewsColumn: React.FC<{ selectedApp: App | null }> = ({
  selectedApp,
}) => {
  const [reviews, setReviews] = useState<Review[] | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    const fetchReviews = async () => {
      console.log("fetchReviews: ", selectedApp);
      if (selectedApp !== null) {
        try {
          setLoading(true);
          const response = await getReviews(selectedApp.Id);
          setReviews(response.Reviews);
          setLoading(false);
        } catch (error: any) {
          setError("Failed to load reviews: " + error.message);
          setLoading(false);
        }
      } else {
        setLoading(false);
        setReviews([]);
      }
    };

    fetchReviews();
  }, [selectedApp]);

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

  if (!reviews)
    return (
      <div className="flex p-8 text-3xl">No reviews in previous 48 hours</div>
    );
  return (
    <div className="container flex flex-col gap-4 p-4">
      {reviews.map((review) => (
        <ReviewCard key={review.Id} review={review} />
      ))}
    </div>
  );
};

export default ReviewsColumn;
