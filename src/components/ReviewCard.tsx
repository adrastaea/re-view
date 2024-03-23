import React from "react";
import { Review } from "../types/Reviews";

const ReviewCard: React.FC<{ review: Review }> = ({ review }) => {
  return (
    <div className="rounded-lg border-2 p-4 shadow-sm hover:shadow-md">
      <div className="justify-left flex flex-row items-center gap-2">
        <div className="flex-shrink-0 text-xl font-bold text-gray-800">
          {review.Author}
        </div>
        <div className="flex-grow">{"⭐️".repeat(review.Score)}</div>
        <div className="flex-grow text-right text-sm font-light">
          {new Date(review.Date).toLocaleString("en-US", {
            month: "short",
            day: "numeric",
            year: "numeric",
            hour: "numeric",
            minute: "numeric",
            hour12: true,
          })}
        </div>
      </div>
      <div className="pt-2">{review.Content}</div>
    </div>
  );
};

export default ReviewCard;
