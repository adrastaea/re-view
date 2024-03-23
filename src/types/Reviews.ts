export interface Reviews {
  Reviews: Review[];
}

export interface Review {
  Id: string;
  Author: string;
  Content: string;
  Score: number;
  Date: Date;
}
