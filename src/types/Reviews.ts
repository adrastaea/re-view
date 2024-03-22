// app/types/Reviews.ts
export interface Reviews {
  Reviews: Review[];
}

export interface Review {
  Id: string;
  Author: string;
  Content: string;
  Score: string;
  Date: string;
}
