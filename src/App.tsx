import ReviewsColumn from "./components/ReviewsColumn";

export interface Review {
  id: number;
  author: string;
  content: string;
  score: number;
  date: string;
}

function App() {
  return (
    <div className="justify-top flex min-h-screen flex-col items-center">
      <h1 className="pt-12 text-5xl font-bold text-[#6a7edc]">Re:View</h1>
      <ReviewsColumn />
    </div>
  );
}

export default App;
