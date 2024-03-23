export interface ReviewsList {
  Reviews: ReviewItem[];
}

export interface ReviewItem {
  Id: string;
  Author: string;
  Content: string;
  Score: number;
  Date: Date;
}
