import ReviewsColumn from "./components/ReviewsColumn";
import "./App.css";

export interface Review {
  id: number;
  author: string;
  content: string;
  score: number;
  date: string;
}

function App() {
  return (
    <>
      <div>
        <h1>Reviews</h1>
        <ReviewsColumn />
      </div>
    </>
  );
}

export default App;
