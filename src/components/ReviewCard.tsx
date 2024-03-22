// app/reviews/ReviewCard.tsx
import React from "react";
import { Review } from "../types/Reviews";

const ReviewCard: React.FC<{ review: Review }> = ({ review }) => {
  return (
    <div
      style={{
        margin: "10px",
        border: "1px solid #ddd",
        padding: "10px",
        borderRadius: "5px",
      }}
    >
      <h3>{review.Author}</h3>
      <p>{review.Content}</p>
      <p>Score: {review.Score}</p>
      <p>Date: {review.Date}</p>
    </div>
  );
};

export default ReviewCard;
