import React from "react";
import { ReviewItem } from "../types/Reviews";

const ReviewCard: React.FC<{ review: ReviewItem }> = ({ review }) => {
  return (
    <div className="rounded-lg border-2 p-4 shadow-sm hover:shadow-md">
      <div className="justify-left flex flex-row flex-wrap items-center gap-2">
        <div className="flex-shrink-0 text-lg font-bold text-gray-800 md:text-2xl">
          {review.Author}
        </div>
        <div className="flex-grow text-sm md:text-lg">
          {"⭐️".repeat(review.Score)}
        </div>
      </div>
      <div className="text-left text-sm font-light">
        {new Date(review.Date).toLocaleString("en-US", {
          month: "short",
          day: "numeric",
          year: "numeric",
          hour: "numeric",
          minute: "numeric",
          hour12: true,
        })}
      </div>
      <div className="pt-2">{review.Content}</div>
    </div>
  );
};

export default ReviewCard;
