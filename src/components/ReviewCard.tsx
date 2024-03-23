// app/reviews/ReviewCard.tsx
import React from "react";
import { Review } from "../types/Reviews";

const ReviewCard: React.FC<{ review: Review }> = ({ review }) => {
  return (
    <div className="rounded-lg border-2 p-4 shadow-sm hover:shadow-md">
      <div className="justify-left flex flex-row gap-2">
        <h3 className="text-xl font-bold text-gray-900">{review.Author}</h3>
        <p className="flex-shrink">{"⭐️".repeat(review.Score)}</p>
        <p className="flex-grow text-right">
          {new Date(review.Date).toLocaleString("en-US", {
            month: "long",
            day: "numeric",
            year: "numeric",
            hour: "numeric",
            minute: "numeric",
            second: "numeric",
            hour12: true,
          })}
        </p>
      </div>
      <p>{review.Content}</p>
    </div>
  );
};

export default ReviewCard;
